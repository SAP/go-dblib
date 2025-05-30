// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-multierror"
)

// Conn handles a TDS-based connection.
//
// Note: This is not the underlying structure for driver.Conn - that is
// Channel.
type Conn struct {
	conn io.ReadWriteCloser
	Caps *CapabilityPackage
	info *Info

	odce odceCipher

	ctx                 context.Context
	ctxCancel           context.CancelFunc
	tdsChannelCurFreeId uint32
	tdsChannels         map[int]*Channel
	tdsChannelsLock     *sync.RWMutex
	errCh               chan error

	// packetSize is the negotiated packet size
	packetSize int
}

// Dial returns a prepared and dialed Conn.
//
// A new child context will be created from the passed context and used
// to abort any interaction with the server - hence closing the parent
// context will abort all interaction with the server.
func NewConn(ctx context.Context, info *Info) (*Conn, error) {
	// Dial returns a prepared and dialed Conn.
	c, err := net.Dial(info.Network, fmt.Sprintf("%s:%s", info.Host, info.Port))
	if err != nil {
		return nil, fmt.Errorf("error opening connection: %w", err)
	}

	if info.TLSEnable || strings.TrimSpace(strings.Replace(info.Port, "0", "", -1)) == "443" {
		tlsConfig := &tls.Config{}
		tlsConfig.ServerName = info.Host
		tlsConfig.InsecureSkipVerify = info.TLSSkipValidation

		if info.TLSHostname != "" {

			hostname := strings.TrimPrefix(info.TLSHostname, "CN=")

			tlsConfig.ServerName = hostname
		}

		if info.TLSCAFile != "" {
			bs, err := os.ReadFile(info.TLSCAFile)
			if err != nil {
				return nil, fmt.Errorf("error reading file at ssl-ca path '%s': %w",
					info.TLSCAFile, err)
			}

			tlsConfig.RootCAs = x509.NewCertPool()

			withCaCert := false

			for {
				var block *pem.Block
				block, bs = pem.Decode(bs)
				if block == nil {
					break
				}

				caCert, err := x509.ParseCertificate(block.Bytes)
				if err != nil {
					return nil, fmt.Errorf("error parsing CA PEM at ssl-ca path '%s': %w",
						info.TLSCAFile, err)
				}

				tlsConfig.RootCAs.AddCert(caCert)
				withCaCert = true
				if len(bs) == 0 {
					break
				}
			}

			if !withCaCert {
				return nil, fmt.Errorf("could not parse any valid CA certificate from file '%s'", info.TLSCAFile)
			}
		}

		tlsClient := tls.Client(c, tlsConfig)
		if err := tlsClient.Handshake(); err != nil {
			return nil, fmt.Errorf("error during TLS handshake with server: %w", err)
		}
		c = tlsClient
	}

	tds := &Conn{
		info:       info,
		conn:       c,
		packetSize: 512,
	}

	if err := tds.setCapabilities(); err != nil {
		return nil, fmt.Errorf("error setting capabilities on connection: %w", err)
	}

	tds.odce = aes_256_cbc

	tds.ctx, tds.ctxCancel = context.WithCancel(ctx)
	// Channels cannot have ID 0 - but channel with the id 0 is used to
	// communicate general packets such as login/logout.
	tds.tdsChannelCurFreeId = uint32(0)
	tds.tdsChannels = make(map[int]*Channel)
	tds.tdsChannelsLock = &sync.RWMutex{}
	tds.errCh = make(chan error, 10)

	// A goroutine automatically reads payloads from the server and
	// passes them to the corresponding channel.
	// Payloads sent to the server are sent in the thread the client
	// uses.
	go tds.ReadFrom()

	return tds, nil
}

// Close closes a Conn and its unclosed Channels.
//
// Teardown and closing on the client side is guaranteed, even if Close
// returns an error. An error is only returned if the communication with
// the server fails or if channels report errors during closing.
//
// If an error is returned it is a *multierror.Error with all errors.
func (tds *Conn) Close() error {
	var me error

	var tdsChannels []*Channel
	// tdsChannels := make([]*Channel, len(tds.tdsChannels))
	for _, channel := range tds.tdsChannels {
		tdsChannels = append(tdsChannels, channel)
	}

	for _, channel := range tdsChannels {
		if err := channel.Close(); err != nil {
			me = multierror.Append(me, fmt.Errorf("error closing channel: %w", err))
		}
	}

	tds.ctxCancel()

	if err := tds.conn.Close(); err != nil {
		me = multierror.Append(me, fmt.Errorf("error closing connection: %w", err))
	}

	return me
}

// PacketSize returns the negotiated packet size.
func (tds *Conn) PacketSize() int {
	// Must be pointer-receive as it is passed to Channels to acquire
	// the negotiated packet size.
	return tds.packetSize
}

// PacketBodySize returns the negotiated packet size minus the packet
// header size.
func (tds *Conn) PacketBodySize() int {
	// Must be pointer-receive as it is passed to Channels to acquire
	// the negotiated packet size.
	return tds.packetSize - PacketHeaderSize
}

func (tds *Conn) getValidChannelId() (int, error) {
	curId := int(tds.tdsChannelCurFreeId)

	if curId > math.MaxUint16 {
		// TODO create error
		return 0, fmt.Errorf("exhausted all channel IDs")
	}

	// increment ID before recursing or returning
	atomic.AddUint32(&tds.tdsChannelCurFreeId, 1)

	if _, ok := tds.tdsChannels[curId]; ok {
		// ChannelId is already used, recurse
		return tds.getValidChannelId()
	}

	return curId, nil
}

// ReadFrom creates packets from payloads from the server and writes
// them to the corresponding Channel.
func (tds *Conn) ReadFrom() {
	for {
		if err := tds.ctx.Err(); err != nil {
			return
		}

		packet := &Packet{}
		_, err := packet.ReadFrom(tds.ctx, tds.conn, time.Duration(tds.info.PacketReadTimeout)*time.Second)
		if err != nil && !errors.Is(err, io.EOF) {
			tds.errCh <- fmt.Errorf("error reading packet: %w", err)
			continue
		}

		tds.tdsChannelsLock.RLock()
		tdsChan, ok := tds.tdsChannels[int(packet.Header.Channel)]
		tds.tdsChannelsLock.RUnlock()
		if !ok {
			tds.errCh <- fmt.Errorf("received packet for invalid channel %d", packet.Header.Channel)
			continue
		}

		// Errors are recorded in the channels' error channel.
		tdsChan.WritePacket(packet)

		// err from packet.ReadFrom
		if errors.Is(err, io.EOF) {
			return
		}
	}
}

func (tds *Conn) setCapabilities() error {
	caps, err := NewCapabilityPackage(
		[]RequestCapability{
			// Support language requests
			TDS_REQ_LANG,
			// Support RPC requests
			// TODO: TDS_REQ_RPC,
			// Support procedure event notifications
			// TODO: TDS_REQ_EVT,
			// Support multiple commands per request
			TDS_REQ_MSTMT,
			// Support bulk copy
			// TODO: TDS_REQ_BCP,
			// Support cursors requests
			// TODO: TDS_REQ_CURSOR,
			// Support dynamic SQL
			TDS_REQ_DYNF,
			// Support MSG requests
			TDS_REQ_MSG,
			// RPC will use TDS_DBRPC and TDS_PARAMFMT / TDS_PARAM
			TDS_REQ_PARAM,

			// Enable all optional data types
			TDS_DATA_INT1,
			TDS_DATA_INT2,
			TDS_DATA_INT4,
			TDS_DATA_BIT,
			TDS_DATA_CHAR,
			TDS_DATA_VCHAR,
			TDS_DATA_BIN,
			TDS_DATA_VBIN,
			TDS_DATA_MNY8,
			TDS_DATA_MNY4,
			TDS_DATA_DATE8,
			TDS_DATA_DATE4,
			TDS_DATA_FLT4,
			TDS_DATA_FLT8,
			TDS_DATA_NUM,
			TDS_DATA_TEXT,
			TDS_DATA_IMAGE,
			TDS_DATA_DEC,
			TDS_DATA_LCHAR,
			TDS_DATA_LBIN,
			TDS_DATA_INTN,
			TDS_DATA_DATETIMEN,
			TDS_DATA_MONEYN,
			TDS_DATA_SENSITIVITY,
			TDS_DATA_BOUNDARY,
			TDS_DATA_FLTN,
			TDS_DATA_BITN,
			TDS_DATA_INT8,
			TDS_DATA_UINT2,
			TDS_DATA_UINT4,
			TDS_DATA_UINT8,
			TDS_DATA_UINTN,
			TDS_DATA_NLBIN,
			TDS_IMAGE_NCHAR,
			TDS_BLOB_NCHAR_16,
			TDS_BLOB_NCHAR_8,
			TDS_BLOB_NCHAR_SCSU,
			TDS_DATA_DATE,
			TDS_DATA_TIME,
			TDS_DATA_INTERVAL,
			TDS_DATA_UNITEXT,
			TDS_DATA_SINT1,
			TDS_REQ_LARGEIDENT,
			TDS_REQ_BLOB_NCHAR_16,
			TDS_DATA_XML,
			TDS_DATA_BIGDATETIME,
			TDS_DATA_USECS,
			//TODO: TDS_DATA_LOBLOCATOR,

			// Support streaming
			//TODO: TDS_OBJECT_CHAR,
			//TODO: TDS_OBJECT_BINARY,

			// Support expedited and non-expedited attentions
			TDS_CON_OOB,
			TDS_CON_INBAND,
			// Use urgent notifications
			TDS_REQ_URGEVT,

			// Create procs from dynamic statements
			TDS_PROTO_DYNPROC,

			// Request status byte in TDS_PARAMS responses
			// Allows to handel nullbytes
			TDS_DATA_COLUMNSTATUS,
			// Support newer versions of tokens
			TDS_REQ_CURINFO3,
			TDS_REQ_DBRPC2,
			// TDS_PARAMFMT2
			TDS_WIDETABLES,

			// Support scrollable cursors
			TDS_CSR_SCROLL,
			TDS_CSR_SENSITIVE,
			TDS_CSR_INSENSITIVE,
			TDS_CSR_SEMISENSITIVE,
			TDS_CSR_KEYSETDRIVEN,

			// Renegotiate packet size after login negotiation
			TDS_REQ_SRVPKTSIZE,

			// Support cluster failover and migration
			//TODO: TDS_CAP_CLUSTERFAILOVER,
			//TODO: TDS_REQ_MIGRATE,

			// Support batched parameters
			TDS_REQ_DYN_BATCH,
			TDS_REQ_LANG_BATCH,
			TDS_REQ_RPC_BATCH,

			// Support on demand encryption
			TDS_REQ_COMMAND_ENCRYPTION,

			// Client will only perform readonly operations
			//TODO: TDS_REQ_READONLY,
		},
		[]ResponseCapability{
			// Ignore format control
			TDS_RES_NO_TDSCONTROL,
		},
		[]SecurityCapability{},
	)

	if err != nil {
		return fmt.Errorf("error creating capability package: %w", err)
	}

	tds.Caps = caps
	return nil
}
