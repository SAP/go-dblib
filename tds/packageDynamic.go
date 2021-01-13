// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import (
	"errors"
	"fmt"
	"math"
)

//go:generate stringer -type=DynamicOperationType

// DynamicOperationType is a bitmask to indicate the options of
// a dynamic SQL statement.
type DynamicOperationType byte

// Types of a dynamic-operation-package.
const (
	TDS_DYN_INVALID    DynamicOperationType = 0x00
	TDS_DYN_PREPARE    DynamicOperationType = 0x01
	TDS_DYN_EXEC       DynamicOperationType = 0x02
	TDS_DYN_DEALLOC    DynamicOperationType = 0x04
	TDS_DYN_EXEC_IMMED DynamicOperationType = 0x08
	TDS_DYN_PROCNAME   DynamicOperationType = 0x10
	TDS_DYN_ACK        DynamicOperationType = 0x20
	TDS_DYN_DESCIN     DynamicOperationType = 0x40
	TDS_DYN_DESCOUT    DynamicOperationType = 0x80
)

//go:generate stringer -type=DynamicStatusType

// DynamicStatusType is a bitmask indicating the status of a dynamic SQL
// statement.
type DynamicStatusType byte

// Types of a dynamic-status-package.
const (
	TDS_DYNAMIC_UNUSED            DynamicStatusType = 0x0
	TDS_DYNAMIC_HASARGS           DynamicStatusType = 0x1
	TDS_DYNAMIC_SUPPRESS_FMT      DynamicStatusType = 0x2
	TDS_DYNAMIC_BATCH_PARAMS      DynamicStatusType = 0x4
	TDS_DYNAMIC_SUPPRESS_PARAMFMT DynamicStatusType = 0x8
)

// DynamicPackage is used to allocate and deallocate as well as request
// information on a dynamic SQL statement.
type DynamicPackage struct {
	Type   DynamicOperationType
	Status DynamicStatusType
	ID     string
	Stmt   string

	wide bool
}

// NewDynamicPackage returns a dynamic package.
func NewDynamicPackage(wide bool) *DynamicPackage {
	dyn := new(DynamicPackage)
	dyn.wide = wide
	return dyn
}

// ReadFrom implements the tds.Package interface.
func (pkg *DynamicPackage) ReadFrom(ch BytesChannel) error {
	var totalLength int
	var err error
	if pkg.wide {
		var length uint32
		length, err = ch.Uint32()
		totalLength = int(length)
	} else {
		var length uint16
		length, err = ch.Uint16()
		totalLength = int(length)
	}
	if err != nil {
		return ErrNotEnoughBytes
	}

	dynamicType, err := ch.Byte()
	if err != nil {
		return ErrNotEnoughBytes
	}
	pkg.Type = DynamicOperationType(dynamicType)
	n := 1

	dynamicStatus, err := ch.Byte()
	if err != nil {
		return ErrNotEnoughBytes
	}
	pkg.Status = DynamicStatusType(dynamicStatus)
	n++

	idLen, err := ch.Uint8()
	if err != nil {
		return ErrNotEnoughBytes
	}
	n++

	pkg.ID, err = ch.String(int(idLen))
	if err != nil {
		return ErrNotEnoughBytes
	}
	n += int(idLen)

	if pkg.Type&TDS_DYN_PREPARE == TDS_DYN_PREPARE || pkg.Type&TDS_DYN_EXEC_IMMED == TDS_DYN_EXEC_IMMED {
		var stmtLen int
		if pkg.wide {
			var length uint32
			length, err = ch.Uint32()
			stmtLen = int(length)
			n += 4
		} else {
			var length uint16
			length, err = ch.Uint16()
			stmtLen = int(length)
			n += 2
		}
		if err != nil {
			return ErrNotEnoughBytes
		}

		pkg.Stmt, err = ch.String(int(stmtLen))
		if err != nil {
			return ErrNotEnoughBytes
		}
		n += int(stmtLen)
	}

	if n != totalLength {
		return fmt.Errorf("expected to read %d bytes, read %d bytes instead", totalLength, n)
	}

	return nil
}

// WriteTo implements the tds.Package interface.
func (pkg *DynamicPackage) WriteTo(ch BytesChannel) error {
	if pkg.Type == TDS_DYN_INVALID {
		return errors.New("dynamic type is invalid")
	}

	token := TDS_DYNAMIC
	if pkg.wide {
		token = TDS_DYNAMIC2
	}

	if err := ch.WriteByte(byte(token)); err != nil {
		return err
	}

	// maxLength is the maximum length of bytes for the package
	maxLength := math.MaxInt16
	if pkg.wide {
		maxLength = math.MaxInt32
	}
	// userMaxLength is the maximum length of bytes for values the user
	// can supply
	userMaxLength := maxLength

	// 1  dynamicType
	// 1  dynamicStatus
	// 1  id length
	// x  id
	totalLength := 3 + len(pkg.ID)
	userMaxLength -= 3
	if pkg.Type&TDS_DYN_PREPARE == TDS_DYN_PREPARE || pkg.Type&TDS_DYN_EXEC_IMMED == TDS_DYN_EXEC_IMMED {
		// 2  stmt length if !pkg.wide
		// 4 stmt length if pkg.wide
		// x  stmt
		totalLength += 2 + len(pkg.Stmt)
		userMaxLength -= 2
		if pkg.wide {
			// add two more bytes for TDS_DYNAMIC2
			totalLength += 2
			userMaxLength -= 2
		}
	}

	if totalLength >= maxLength {
		return fmt.Errorf("tds: query too long, statement ID and query can at most be %d bytes long",
			userMaxLength)
	}

	if pkg.wide {
		if err := ch.WriteUint32(uint32(totalLength)); err != nil {
			return err
		}
	} else {
		if err := ch.WriteUint16(uint16(totalLength)); err != nil {
			return err
		}
	}

	if err := ch.WriteByte(byte(pkg.Type)); err != nil {
		return err
	}
	n := 1

	if err := ch.WriteByte(byte(pkg.Status)); err != nil {
		return err
	}
	n++

	if err := ch.WriteUint8(uint8(len(pkg.ID))); err != nil {
		return err
	}
	n++

	if err := ch.WriteString(pkg.ID); err != nil {
		return err
	}
	n += len(pkg.ID)

	if pkg.Type&TDS_DYN_PREPARE == TDS_DYN_PREPARE || pkg.Type&TDS_DYN_EXEC_IMMED == TDS_DYN_EXEC_IMMED {
		if pkg.wide {
			if err := ch.WriteUint32(uint32(len(pkg.Stmt))); err != nil {
				return err
			}
			n += 4
		} else {
			if err := ch.WriteUint16(uint16(len(pkg.Stmt))); err != nil {
				return err
			}
			n += 2
		}

		if err := ch.WriteString(pkg.Stmt); err != nil {
			return err
		}
		n += len(pkg.Stmt)
	}

	if n != totalLength {
		return fmt.Errorf("expected to write %d bytes, wrote %d bytes instead", totalLength, n)
	}

	return nil
}

func (pkg DynamicPackage) String() string {
	strTypes := deBitmaskString(int(pkg.Type), int(TDS_DYN_DESCOUT),
		func(i int) string { return DynamicOperationType(i).String() },
		TDS_DYN_PREPARE.String(),
	)

	strStati := deBitmaskString(int(pkg.Status), int(TDS_DYNAMIC_SUPPRESS_FMT),
		func(i int) string { return DynamicStatusType(i).String() },
		TDS_DYNAMIC_UNUSED.String(),
	)

	return fmt.Sprintf("%T(%s, %s - %s: %s)", pkg, strTypes, strStati, pkg.ID, pkg.Stmt)
}
