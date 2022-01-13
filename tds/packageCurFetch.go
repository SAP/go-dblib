// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import "fmt"

//go:generate stringer -type=CursorFetchType

// CursorFetchType is the type to indicate in which direction to fetch
// cursor results.
type CursorFetchType uint

const (
	TDS_CUR_NEXT CursorFetchType = iota + 1
	TDS_CUR_PREV
	TDS_CUR_FIRST
	TDS_CUR_LAST
	TDS_CUR_ABS
	TDS_CUR_REL
)

var _ Package = (*CurFetchPackage)(nil)

// CurFetchPackage is used to fetch more rows from a cursor.
type CurFetchPackage struct {
	CursorID int32
	Name     string

	Type      CursorFetchType
	RowNumber int32
}

// ReadFrom implements the tds.Package interface.
func (pkg *CurFetchPackage) ReadFrom(ch BytesChannel) error {
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

	fetchType, err := ch.Uint8()
	if err != nil {
		return err
	}
	pkg.Type = CursorFetchType(fetchType)
	n++

	if pkg.Type == TDS_CUR_ABS || pkg.Type == TDS_CUR_REL {
		pkg.RowNumber, err = ch.Int32()
		if err != nil {
			return ErrNotEnoughBytes
		}
		n += 4
	}

	if n != int(totalLength) {
		return fmt.Errorf("expected to read %d bytes, read %d bytes instead", totalLength, n)
	}

	return nil
}

// WriteTo implements the tds.Package interface.
func (pkg CurFetchPackage) WriteTo(ch BytesChannel) error {
	if err := ch.WriteByte(byte(TDS_CURFETCH)); err != nil {
		return err
	}

	// 4 cursorID
	// if cursorID == 0:
	//   1 name len
	//   x name
	// 1 type
	// if type ABS/REL:
	//   4 rownumber
	totalLength := 4 + 1
	if pkg.CursorID == 0 {
		totalLength += 1 + len(pkg.Name)
	}
	if pkg.Type == TDS_CUR_ABS || pkg.Type == TDS_CUR_REL {
		totalLength += 4
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

	if err := ch.WriteUint8(uint8(pkg.Type)); err != nil {
		return err
	}

	if pkg.Type == TDS_CUR_ABS || pkg.Type == TDS_CUR_REL {
		if err := ch.WriteInt32(pkg.RowNumber); err != nil {
			return err
		}
	}

	return nil
}

func (pkg CurFetchPackage) String() string {
	return fmt.Sprintf("%T(%d, %s, %s, %d)", pkg, pkg.CursorID, pkg.Name, pkg.Type, pkg.RowNumber)
}
