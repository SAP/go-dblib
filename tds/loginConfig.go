// SPDX-FileCopyrightText: 2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/SAP/go-dblib/dsn"
)

// LoginConfigRemoteServer contains the name and the password to the
// server.
type LoginConfigRemoteServer struct {
	Name, Password string
}

// LoginConfig contains configuration to login to the server.
type LoginConfig struct {
	DSN      *dsn.Info
	Hostname string

	// TODO name
	HostProc string
	AppName  string
	ServName string

	Language string
	CharSet  string

	RemoteServers []LoginConfigRemoteServer

	// Encrypt allows any TDSMsgId but only negotiation-relevant security
	// bits such as TDS_MSG_SEC_ENCRYPT will be recognized.
	Encrypt TDSMsgId
}

// NewLoginConfig creates a new login-configuration by using dsn
// information and setting default configuration-values in regard to the
// ASE database server (Should be be adjusted by clients).
func NewLoginConfig(dsn *dsn.Info) (*LoginConfig, error) {
	conf := &LoginConfig{}

	conf.DSN = dsn

	if dsn.ClientHostname != "" {
		conf.Hostname = dsn.ClientHostname
	} else {
		hostname, err := os.Hostname()
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve hostname: %w", err)
		}
		conf.Hostname = hostname
	}
	if len(conf.Hostname) > 30 {
		conf.Hostname = conf.Hostname[:30]
	}

	conf.HostProc = strconv.Itoa(os.Getpid())

	conf.ServName = conf.DSN.Host
	if len(conf.ServName) > 30 {
		conf.ServName = conf.ServName[:30]
	}

	// Should be overwritten by clients
	conf.AppName = "github.com/SAP/go-ase/libase/tds"

	conf.CharSet = "utf8"
	conf.Language = "us_english"

	conf.Encrypt = TDS_MSG_SEC_ENCRYPT4

	return conf, nil
}

// TDS default login-configuration values.
const (
	TDS_MAXNAME   = 30
	TDS_NETBUF    = 4
	TDS_RPLEN     = 255
	TDS_VERSIZE   = 4
	TDS_PROGNLEN  = 10
	TDS_OLDSECURE = 2
	TDS_HA        = 6
	TDS_SECURE    = 2
	TDS_PKTLEN    = 6
	TDS_DUMMY     = 4
)

func (config *LoginConfig) pack() (Package, error) {
	buf := &bytes.Buffer{}

	// lhostname, lhostlen
	if err := writeString(buf, config.Hostname, TDS_MAXNAME); err != nil {
		return nil, fmt.Errorf("error writing hostname: %w", err)
	}

	// lusername, lusernlen
	if err := writeString(buf, config.DSN.Username, TDS_MAXNAME); err != nil {
		return nil, fmt.Errorf("error writing username: %w", err)
	}

	// lpw, lpwnlen
	var err error
	switch config.Encrypt {
	case TDS_MSG_SEC_ENCRYPT, TDS_MSG_SEC_ENCRYPT2, TDS_MSG_SEC_ENCRYPT3, TDS_MSG_SEC_ENCRYPT4:
		err = writeString(buf, "", TDS_MAXNAME)
	default:
		err = writeString(buf, config.DSN.Password, TDS_MAXNAME)
	}
	if err != nil {
		return nil, fmt.Errorf("error writing password: %w", err)
	}

	// lhostproc, lhplen
	if err := writeString(buf, config.HostProc, TDS_MAXNAME); err != nil {
		return nil, fmt.Errorf("error writing hostproc: %w", err)
	}

	// lint2
	if _, err := writeBasedOnEndian(buf, 3, 2); err != nil {
		return nil, fmt.Errorf("error writing int2: %w", err)
	}

	// lint4
	if _, err := writeBasedOnEndian(buf, 1, 0); err != nil {
		return nil, fmt.Errorf("error writing int4: %w", err)
	}

	// lchar -> ASCII
	if err := buf.WriteByte(6); err != nil {
		return nil, fmt.Errorf("error writing char: %w", err)
	}

	// lflt
	if _, err := writeBasedOnEndian(buf, 10, 4); err != nil {
		return nil, fmt.Errorf("error writing flt: %w", err)
	}

	// ldate
	if _, err := writeBasedOnEndian(buf, 9, 8); err != nil {
		return nil, fmt.Errorf("error writing date: %w", err)
	}

	// lusedb
	if err := buf.WriteByte(1); err != nil {
		return nil, fmt.Errorf("error writing usedb: %w", err)
	}

	// ldmpld
	if err := buf.WriteByte(1); err != nil {
		return nil, fmt.Errorf("error writing dmpld: %w", err)
	}

	// only relevant for server-server comm
	// linterfacespare
	if err := buf.WriteByte(0); err != nil {
		return nil, fmt.Errorf("error writing interfacespare: %w", err)
	}

	// ltype
	if err := buf.WriteByte(0); err != nil {
		return nil, fmt.Errorf("error writing type: %w", err)
	}

	// deprecated
	// lbufsize
	if _, err := buf.Write(make([]byte, TDS_NETBUF)); err != nil {
		return nil, fmt.Errorf("error writing bufsize: %w", err)
	}

	// lspare
	if _, err := buf.Write(make([]byte, 3)); err != nil {
		return nil, fmt.Errorf("error writing spare: %w", err)
	}

	// lappname, lappnlen
	if err := writeString(buf, config.AppName, TDS_MAXNAME); err != nil {
		return nil, fmt.Errorf("error writing appname: %w", err)
	}

	// lservname, lservnlen
	if err := writeString(buf, config.ServName, TDS_MAXNAME); err != nil {
		return nil, fmt.Errorf("error writing servname: %w", err)
	}

	// TODO only relevant for server-server comm, replace?
	// lrempw, lrempwlen
	switch config.Encrypt {
	case TDS_MSG_SEC_ENCRYPT, TDS_MSG_SEC_ENCRYPT2, TDS_MSG_SEC_ENCRYPT3, TDS_MSG_SEC_ENCRYPT4:
		err = writeString(buf, "", TDS_RPLEN)
	default:
		err = writeString(buf, "", TDS_RPLEN)
	}
	if err != nil {
		return nil, fmt.Errorf("error writing rempw: %w", err)
	}

	// ltds
	if _, err := buf.Write([]byte{0x5, 0x0, 0x0, 0x0}); err != nil {
		return nil, fmt.Errorf("error writing tds version: %w", err)
	}

	// lprogname, lprognlen
	if err := writeString(buf, libraryName, TDS_PROGNLEN); err != nil {
		return nil, fmt.Errorf("error writing progname: %w", err)
	}

	// lprogvers
	if _, err := buf.Write(libraryVersion.Bytes()); err != nil {
		return nil, fmt.Errorf("error writing progversion: %w", err)
	}

	// lnoshort - do not convert short data types
	if err := buf.WriteByte(0); err != nil {
		return nil, fmt.Errorf("error writing noshort: %w", err)
	}

	// lflt4
	if _, err := writeBasedOnEndian(buf, 13, 12); err != nil {
		return nil, fmt.Errorf("error writing flt4: %w", err)
	}

	// ldate4
	if _, err := writeBasedOnEndian(buf, 17, 16); err != nil {
		return nil, fmt.Errorf("error writing date4: %w", err)
	}

	// llanguage, llanglen
	if err := writeString(buf, config.Language, TDS_MAXNAME); err != nil {
		return nil, fmt.Errorf("error writing language: %w", err)
	}

	// lsetlang - notify of language changes
	if err := buf.WriteByte(1); err != nil {
		return nil, fmt.Errorf("error writing setlang: %w", err)
	}

	// loldsecure - deprecated
	if _, err := buf.Write(make([]byte, TDS_OLDSECURE)); err != nil {
		return nil, fmt.Errorf("error writing oldsecure: %w", err)
	}

	// lseclogin
	switch config.Encrypt {
	case TDS_MSG_SEC_ENCRYPT:
		err = buf.WriteByte(0x01)
	case TDS_MSG_SEC_ENCRYPT2:
		err = buf.WriteByte(0x1 | 0x20)
	case TDS_MSG_SEC_ENCRYPT3, TDS_MSG_SEC_ENCRYPT4:
		err = buf.WriteByte(0x1 | 0x20 | 0x80)
	default:
		err = buf.WriteByte(0x0)
	}
	if err != nil {
		return nil, fmt.Errorf("error writing seclogin: %w", err)
	}

	// lsecbulk - deprecated
	if err := buf.WriteByte(1); err != nil {
		return nil, fmt.Errorf("error writing secbulk: %w", err)
	}

	// lhalogin
	// TODO - values need to be determined by config to allow for
	// failover reconnects in clusters
	if err := buf.WriteByte(1); err != nil {
		return nil, fmt.Errorf("error writing ailover: %w", err)
	}

	// lhasessionid
	// TODO session id for HA failover, find out if this needs to be
	// user set or retrieved from the server
	if _, err := buf.Write(make([]byte, TDS_HA)); err != nil {
		return nil, fmt.Errorf("error writing hasessionid: %w", err)
	}

	// lsecspare - unused
	// TODO TDS_SECURE unknown
	if _, err := buf.Write(make([]byte, TDS_SECURE)); err != nil {
		return nil, fmt.Errorf("error writing secspare: %w", err)
	}

	// lcharset, lcharsetlen
	if err := writeString(buf, config.CharSet, TDS_MAXNAME); err != nil {
		return nil, fmt.Errorf("error writing charset: %w", err)
	}

	// lsetcharset - notify of charset changes
	if err := buf.WriteByte(1); err != nil {
		return nil, fmt.Errorf("error writing setcharset: %w", err)
	}

	// lpacketsize - 256 to 65535 bytes
	// Write default packet size - will be renegotiated anyhow.
	if err := writeString(buf, "512", TDS_PKTLEN); err != nil {
		return nil, fmt.Errorf("error writing packetsize: %w", err)
	}

	// ldummy - apparently unused
	if _, err := buf.Write(make([]byte, TDS_DUMMY)); err != nil {
		return nil, fmt.Errorf("error writing dummy: %w", err)
	}

	return &TokenlessPackage{Data: buf}, nil
}
