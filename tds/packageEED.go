// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import (
	"fmt"
	"strings"
)

//go:generate stringer -type=EEDStatus

// EEDStatus is used to indicate the status of an EED package.
type EEDStatus uint8

const (
	TDS_NO_EED      EEDStatus = 0x00
	TDS_EED_FOLLOWS EEDStatus = 0x1
	TDS_EED_INFO    EEDStatus = 0x2
)

// EEDPackage is used to communicate information and error messages.
type EEDPackage struct {
	MsgNumber  uint32
	State      uint8
	Class      uint8
	SQLState   []byte
	Status     EEDStatus
	TranState  uint16
	Msg        string
	ServerName string
	ProcName   string
	LineNr     uint16
}

// ReadFrom implements the tds.Package interface.
func (pkg *EEDPackage) ReadFrom(ch BytesChannel) error {
	length, err := ch.Uint16()
	if err != nil {
		return ErrNotEnoughBytes
	}

	pkg.MsgNumber, err = ch.Uint32()
	if err != nil {
		return ErrNotEnoughBytes
	}
	n := 4

	pkg.State, err = ch.Uint8()
	if err != nil {
		return ErrNotEnoughBytes
	}
	n++

	pkg.Class, err = ch.Uint8()
	if err != nil {
		return ErrNotEnoughBytes
	}
	n++

	sqlStateLen, err := ch.Uint8()
	if err != nil {
		return ErrNotEnoughBytes
	}
	n++

	pkg.SQLState, err = ch.Bytes(int(sqlStateLen))
	if err != nil {
		return ErrNotEnoughBytes
	}
	n += int(sqlStateLen)

	status, err := ch.Uint8()
	if err != nil {
		return ErrNotEnoughBytes
	}
	pkg.Status = EEDStatus(status)
	n++

	pkg.TranState, err = ch.Uint16()
	if err != nil {
		return ErrNotEnoughBytes
	}
	n += 2

	msgLength, err := ch.Uint16()
	if err != nil {
		return ErrNotEnoughBytes
	}
	n += 2

	msg, err := ch.String(int(msgLength))
	if err != nil {
		return ErrNotEnoughBytes
	}
	// Some messages contain a trailing newline, but not all.
	pkg.Msg = strings.TrimSuffix(msg, "\n")
	n += int(msgLength)

	serverLength, err := ch.Uint8()
	if err != nil {
		return ErrNotEnoughBytes
	}
	n++

	pkg.ServerName, err = ch.String(int(serverLength))
	if err != nil {
		return ErrNotEnoughBytes
	}
	n += int(serverLength)

	procLength, err := ch.Uint8()
	if err != nil {
		return ErrNotEnoughBytes
	}
	n++

	pkg.ProcName, err = ch.String(int(procLength))
	if err != nil {
		return ErrNotEnoughBytes
	}
	n += int(procLength)

	pkg.LineNr, err = ch.Uint16()
	if err != nil {
		return ErrNotEnoughBytes
	}
	n += 2

	if n != int(length) {
		return fmt.Errorf("expected to read %d bytes, read %d bytes instead", length, n)
	}

	return nil
}

// WriteTo implements the tds.Package interface.
func (pkg EEDPackage) WriteTo(ch BytesChannel) error {
	if err := ch.WriteByte(byte(TDS_EED)); err != nil {
		return fmt.Errorf("failed to write TDS Token %s: %w", TDS_EED, err)
	}

	// 4 msgnumber
	// 1 state
	// 1 class
	// x sqlstate
	// 1 status
	// 2 transtate
	// x msg
	// x servername
	// x procname
	// 2 linenr
	length := 11 + len(pkg.SQLState) + len(pkg.Msg) + len(pkg.ServerName) + len(pkg.ProcName)

	if err := ch.WriteUint16(uint16(length)); err != nil {
		return fmt.Errorf("failed to write length: %w", err)
	}

	if err := ch.WriteUint32(pkg.MsgNumber); err != nil {
		return fmt.Errorf("failed to write message number: %w", err)
	}

	if err := ch.WriteUint8(pkg.State); err != nil {
		return fmt.Errorf("failed to write state: %w", err)
	}

	if err := ch.WriteUint8(pkg.Class); err != nil {
		return fmt.Errorf("failed to write class: %w", err)
	}

	if err := ch.WriteUint8(uint8(len(pkg.SQLState))); err != nil {
		return fmt.Errorf("failed to write SQL state len: %w", err)
	}

	if err := ch.WriteBytes(pkg.SQLState); err != nil {
		return fmt.Errorf("failed to write SQL state: %w", err)
	}

	if err := ch.WriteByte(byte(pkg.Status)); err != nil {
		return fmt.Errorf("failed to write status: %w", err)
	}

	if err := ch.WriteUint16(pkg.TranState); err != nil {
		return fmt.Errorf("failed to write tran state: %w", err)
	}

	if err := ch.WriteUint16(uint16(len(pkg.Msg))); err != nil {
		return fmt.Errorf("failed to write message length: %w", err)
	}

	if err := ch.WriteString(pkg.Msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err := ch.WriteUint8(uint8(len(pkg.ServerName))); err != nil {
		return fmt.Errorf("failed to write server name length: %w", err)
	}

	if err := ch.WriteString(pkg.ServerName); err != nil {
		return fmt.Errorf("failed to write server name: %w", err)
	}

	if err := ch.WriteUint8(uint8(len(pkg.ProcName))); err != nil {
		return fmt.Errorf("failed to write proc name length: %w", err)
	}

	if err := ch.WriteString(pkg.ProcName); err != nil {
		return fmt.Errorf("failed to write proc name: %w", err)
	}

	if err := ch.WriteUint16(pkg.LineNr); err != nil {
		return fmt.Errorf("failed to write line nr: %w", err)
	}

	return nil
}

func (pkg EEDPackage) String() string {
	return fmt.Sprintf("%T(%s - %d: %s)", pkg, pkg.Status, pkg.MsgNumber, pkg.Msg)
}
