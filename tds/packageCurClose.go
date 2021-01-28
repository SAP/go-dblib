// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import "fmt"

//go:generate stringer -type=CursorCloseOption

// CursorCloseOption is the type of values for CurClose packages.
type CursorCloseOption uint

const (
	TDS_CUR_COPT_UNUSED  CursorCloseOption = 0x0
	TDS_CUR_COPT_DEALLOC CursorCloseOption = 0x1
)

// TODO
type CurClosePackage struct {
	CursorID int32
	Name     string

	Options CursorCloseOption
}

// ReadFrom implements the tds.Package interface.
func (pkg *CurClosePackage) ReadFrom(ch BytesChannel) error {
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

	options, err := ch.Uint8()
	if err != nil {
		return ErrNotEnoughBytes
	}
	pkg.Options = CursorCloseOption(options)
	n++

	if n != int(totalLength) {
		return fmt.Errorf("expected to read %d bytes, read %d bytes instead", totalLength, n)
	}

	return nil
}

// WriteTo implements the tds.Package interface.
func (pkg CurClosePackage) WriteTo(ch BytesChannel) error {
	if err := ch.WriteByte(byte(TDS_CURCLOSE)); err != nil {
		return err
	}

	// 4 cursorID
	// if cursorID == 0:
	//   1 name len
	//   x name
	// 1 options
	totalLength := 4 + 1
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

	if err := ch.WriteUint8(uint8(pkg.Options)); err != nil {
		return err
	}

	return nil
}

// String implements the fmt.Stringer interface.
func (pkg CurClosePackage) String() string {
	strOpts := deBitmaskString(int(pkg.Options), int(TDS_CUR_COPT_DEALLOC),
		func(i int) string { return CursorCloseOption(i).String() },
		TDS_CUR_COPT_UNUSED.String(),
	)

	return fmt.Sprintf("%T(%d, %s, %s)", pkg, pkg.CursorID, pkg.Name, strOpts)
}
