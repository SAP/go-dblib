// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import "fmt"

var _ Package = (*CurUpdatePackage)(nil)

// CurUpdatePackage is used to update a cursor.
type CurUpdatePackage struct {
	CursorID int32
	Name     string

	Status    CursorOStatus
	TableName string
	Stmt      string
}

// ReadFrom implements the tds.Package interface.
func (pkg *CurUpdatePackage) ReadFrom(ch BytesChannel) error {
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
	pkg.Status = CursorOStatus(status)
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

	// TODO Only in language option case?
	stmtLength, err := ch.Uint16()
	if err != nil {
		return ErrNotEnoughBytes
	}
	n += 2

	pkg.Stmt, err = ch.String(int(stmtLength))
	if err != nil {
		return ErrNotEnoughBytes
	}
	n += int(stmtLength)
	// end TODO

	if n != int(totalLength) {
		return fmt.Errorf("expected to read %d bytes, read %d bytes instead", totalLength, n)
	}

	return nil
}

// WriteTo implements the tds.Package interface.
func (pkg CurUpdatePackage) WriteTo(ch BytesChannel) error {
	if err := ch.WriteByte(byte(TDS_CURUPDATE)); err != nil {
		return err
	}

	// 4 cursorID
	// if cursorID == 0:
	//   1 name len
	//   x name
	// 1 status
	// 1 table name len
	// x table name
	// if pkg.Stmt:
	//   2 stmt len
	//   x stmt
	totalLength := 4 + 1 + 1 + len(pkg.TableName)
	if pkg.CursorID == 0 {
		totalLength += 1 + len(pkg.Name)
	}
	if len(pkg.Stmt) > 0 {
		totalLength += 2 + len(pkg.Stmt)
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

	if len(pkg.Stmt) > 0 {
		if err := ch.WriteUint16(uint16(len(pkg.Stmt))); err != nil {
			return err
		}

		if err := ch.WriteString(pkg.Stmt); err != nil {
			return err
		}
	}

	return nil
}

// String implements the fmt.Stringer interface.
func (pkg CurUpdatePackage) String() string {
	return fmt.Sprintf("%T(%d, %s, %s, %s, %q)", pkg, pkg.CursorID, pkg.Name,
		pkg.Status, pkg.TableName, pkg.Stmt)
}
