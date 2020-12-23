// SPDX-FileCopyrightText: 2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import (
	"errors"
	"fmt"
	"io"
)

// ErrNotEnoughBytes is returned by packages' ReadFrom if the
// BytesChannel does not have enough bytes to parse the package fully.
var ErrNotEnoughBytes = errors.New("not enough bytes in channel to parse package")

// Package is the interface providing the ReadFrom and WriteTo methods.
type Package interface {
	// ReadFrom reads bytes from the passed channel until either the
	// channel is closed or the package has all required information.
	// The read bytes are parsed into the package struct.
	ReadFrom(BytesChannel) error

	// WriteTo writes bytes to the passed channel until either the
	// channel is closed or the package has written all required
	// information.
	WriteTo(BytesChannel) error

	fmt.Stringer
}

// LookupPackage returns the Package struct for a token.
func LookupPackage(token Token) (Package, error) {
	switch token {
	case TDS_EED:
		return &EEDPackage{}, nil
	case TDS_ERROR:
		return &ErrorPackage{}, nil
	case TDS_LOGINACK:
		return &LoginAckPackage{}, nil
	case TDS_DONE:
		return &DonePackage{}, nil
	case TDS_DONEPROC:
		return &DoneProcPackage{}, nil
	case TDS_DONEINPROC:
		return &DoneInProcPackage{}, nil
	case TDS_MSG:
		return &MsgPackage{}, nil
	case TDS_PARAMFMT:
		return &ParamFmtPackage{}, nil
	case TDS_PARAMFMT2:
		return &ParamFmtPackage{wide: true}, nil
	case TDS_ROWFMT:
		return &RowFmtPackage{}, nil
	case TDS_ROWFMT2:
		return &RowFmtPackage{wide: true}, nil
	case TDS_PARAMS:
		return &ParamsPackage{}, nil
	case TDS_ROW:
		return &RowPackage{}, nil
	case TDS_CAPABILITY:
		return NewCapabilityPackage(nil, nil, nil)
	case TDS_ENVCHANGE:
		return &EnvChangePackage{}, nil
	case TDS_LANGUAGE:
		return &LanguagePackage{}, nil
	case TDS_ORDERBY:
		return &OrderByPackage{}, nil
	case TDS_ORDERBY2:
		return &OrderBy2Package{}, nil
	case TDS_RETURNSTATUS:
		return &ReturnStatusPackage{}, nil
	case TDS_LOGOUT:
		return &LogoutPackage{}, nil
	case TDS_DYNAMIC:
		return &DynamicPackage{}, nil
	case TDS_DYNAMIC2:
		return &DynamicPackage{wide: true}, nil
	default:
		return NewTokenlessPackage(), nil
	}
}

// IsError returns true if the package signals an error - either in
// the communication or in the protocol.
func IsError(pkg Package) bool {
	switch pkg.(type) {
	case *EEDPackage, *ErrorPackage:
		return true
	}

	return false
}

// IsDone returns true if the package terminates the stream.
func IsDone(pkg Package) bool {
	switch pkg.(type) {
	case *DonePackage:
		return true
	}

	return false
}

// BytesChannel defines the required methods for Packages to read and
// write information to a stream.
type BytesChannel interface {
	// Position marks the index of the packet and index of the byte the
	// channel currently is at.
	// The position is considered volatile and only valid until the next
	// call to DiscardUntilPosition.
	Position() (int, int)
	SetPosition(int, int)
	DiscardUntilCurrentPosition()

	io.Reader
	io.Writer

	Bytes(n int) ([]byte, error)
	WriteBytes([]byte) error

	Byte() (byte, error)
	WriteByte(byte) error

	Uint8() (uint8, error)
	WriteUint8(uint8) error

	Int8() (int8, error)
	WriteInt8(int8) error

	Uint16() (uint16, error)
	WriteUint16(uint16) error

	Int16() (int16, error)
	WriteInt16(int16) error

	Uint32() (uint32, error)
	WriteUint32(uint32) error

	Int32() (int32, error)
	WriteInt32(int32) error

	Uint64() (uint64, error)
	WriteUint64(uint64) error

	Int64() (int64, error)
	WriteInt64(int64) error

	String(int) (string, error)
	WriteString(string) error
}
