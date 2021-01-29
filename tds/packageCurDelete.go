// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import "fmt"

//go:generate stringer -type=CursorDeleteStatus

// CursorDeleteStatus is the type for a currently unused bitmask.
type CursorDeleteStatus uint

const (
	TDS_CUR_DELSTAT_UNUSED CursorDeleteStatus = iota
)

var _ Package = (*CurDeletePackage)(nil)

// CurDeletePackage is used to deallocate a cursor.
type CurDeletePackage struct {
	CursorID int32
	Name     string

	Status    CursorDeleteStatus
	TableName string
}

// ReadFrom implements the tds.Package interface.
func (pkg *CurDeletePackage) ReadFrom(ch BytesChannel) error {
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
		nameLen, err := ch.Uint8()
		if err != nil {
			return ErrNotEnoughBytes
		}
		n++

		pkg.Name, err = ch.String(int(nameLen))
		if err != nil {
			return ErrNotEnoughBytes
		}
		n += int(nameLen)
	}

	status, err := ch.Uint8()
	if err != nil {
		return ErrNotEnoughBytes
	}
	pkg.Status = CursorDeleteStatus(status)
	n++

	tableNameLength, err := ch.Uint8()
	if err != nil {
		return ErrNotEnoughBytes
	}
	n++

	pkg.TableName, err = ch.String(int(tableNameLength))
	if err != nil {
		return ErrNotEnoughBytes
	}
	n += int(tableNameLength)

	if n != int(totalLength) {
		return fmt.Errorf("expected to read %d bytes, read %d bytes instead", totalLength, n)
	}

	return nil
}

// WriteTo implements the tds.Package interface.
func (pkg CurDeletePackage) WriteTo(ch BytesChannel) error {
	if err := ch.WriteByte(byte(TDS_CURDELETE)); err != nil {
		return err
	}

	// 4 cursorID
	// if cursorID == 0:
	//   1 name len
	//   x name
	// 1 status
	// 1 tablename len
	// x tablename
	totalLength := 4 + 1 + 1 + len(pkg.TableName)
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

	if err := ch.WriteUint8(uint8(len(pkg.TableName))); err != nil {
		return err
	}

	if err := ch.WriteString(pkg.TableName); err != nil {
		return err
	}

	return nil
}

func (pkg CurDeletePackage) String() string {
	return fmt.Sprintf("%T(%d, %s, %s, %s)", pkg, pkg.CursorID, pkg.Name,
		pkg.TableName, pkg.Status)
}
