// SPDX-FileCopyrightText: 2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import "fmt"

//go:generate stringer -type=CursorOStatus

// CursorOStatus is the type for a value indicating the status of
// a cursor when opening it.
type CursorOStatus uint

const (
	TDS_CUR_OSTAT_UNUSED CursorOStatus = iota
	TDS_CUR_OSTAT_HASARGS
	TDS_CUR_CONSEC_UPDS
)

var _ Package = (*CurOpenPackage)(nil)

// CurOpenPackage is used to open a declared cursor.
type CurOpenPackage struct {
	CursorID int32
	Name     string
	Status   CursorOStatus
}

// ReadFrom implements the tds.Package interface.
func (pkg *CurOpenPackage) ReadFrom(ch BytesChannel) error {
	totalLength, err := ch.Uint16()
	if err != nil {
		return ErrNotEnoughBytes
	}

	pkg.CursorID, err = ch.Int32()
	if err != nil {
		return ErrNotEnoughBytes
	}
	n := 4

	if pkg.CursorID == 0 {
		nameLength, err := ch.Uint8()
		if err != nil {
			return ErrNotEnoughBytes
		}
		n++

		pkg.Name, err = ch.String(int(nameLength))
		if err != nil {
			return ErrNotEnoughBytes
		}
		n += int(nameLength)
	}

	status, err := ch.Uint8()
	if err != nil {
		return ErrNotEnoughBytes
	}
	pkg.Status = CursorOStatus(status)
	n++

	if n != int(totalLength) {
		return fmt.Errorf("expected to read %d bytes, read %d bytes instead", totalLength, n)
	}

	return nil
}

// WriteTo implements the tds.Package interface.
func (pkg CurOpenPackage) WriteTo(ch BytesChannel) error {
	if err := ch.WriteByte(byte(TDS_CUROPEN)); err != nil {
		return err
	}

	// 4 cursorid
	// if cursorId == 0:
	//   1 name len
	//   x name
	// 1 status
	totalLength := 5
	if pkg.CursorID == 0 {
		totalLength += 1 + len(pkg.Name)
	}

	if err := ch.WriteUint16(uint16(totalLength)); err != nil {
		return err
	}

	if err := ch.WriteInt32(pkg.CursorID); err != nil {
		return err
	}

	if pkg.CursorID == 0 {
		if err := ch.WriteUint8(uint8(len(pkg.Name))); err != nil {
			return err
		}

		if err := ch.WriteString(pkg.Name); err != nil {
			return err
		}
	}

	if err := ch.WriteUint8(uint8(pkg.Status)); err != nil {
		return err
	}

	return nil
}

// String implements the fmt.Stringer interface.
func (pkg CurOpenPackage) String() string {
	return fmt.Sprintf("%T(%d, %s, %s)", pkg, pkg.CursorID, pkg.Name, pkg.Status)
}
