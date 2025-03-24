// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import "fmt"

//go:generate stringer -type=CursorOption

// CursorOption is a bitmask to indicate options for a cursor.
type CursorOption uint

const (
	TDS_CUR_DOPT_UNUSED          CursorOption = 0x0
	TDS_CUR_DOPT_RDONLY          CursorOption = 0x1
	TDS_CUR_DOPT_UPDATABLE       CursorOption = 0x2
	TDS_CUR_DOPT_SENSITIVE       CursorOption = 0x4
	TDS_CUR_DOPT_DYNAMIC         CursorOption = 0x8
	TDS_CUR_DOPT_IMPLICIT        CursorOption = 0x10
	TDS_CUR_DOPT_INSENSITIVE     CursorOption = 0x20
	TDS_CUR_DOPT_SEMISENSITIVE   CursorOption = 0x40
	TDS_CUR_DOPT_KEYSETDRIVEN    CursorOption = 0x80
	TDS_CUR_DOPT_SCROLLABLE      CursorOption = 0x100
	TDS_CUR_DOPT_RELLOCKSONCLOSE CursorOption = 0x200
)

//go:generate stringer -type=CursorDStatus

// CursorDStatus is the type to indicate the status of a cursor when
// declaring it.
type CursorDStatus uint

const (
	TDS_CUR_DSTAT_UNUSED  CursorDStatus = 0x0
	TDS_CUR_DSTAT_HASARGS CursorDStatus = 0x1
)

var _ Package = (*CurDeclarePackage)(nil)

// CurDeclarePackage is used to declare cursors.
type CurDeclarePackage struct {
	Name    string
	Options CursorOption
	Status  CursorDStatus
	Stmt    string

	columns []string

	wide bool
}

// NewCurDeclarePackage takes all required values and returns
// a CurDeclarePackage with the required values and sensible defaults.
//
// TODO: An error is only returned if ...
func NewCurDeclarePackage(name, stmt string, status CursorDStatus, options CursorOption) (*CurDeclarePackage, error) {
	pkg := &CurDeclarePackage{
		Name:    name,
		Options: options,
		Status:  status,
		Stmt:    stmt,
		wide:    true,
	}

	// TODO check if wide is supported

	return pkg, nil
}

// ReadFrom implements the tds.Package interface.
func (pkg *CurDeclarePackage) ReadFrom(ch BytesChannel) error {
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

	nameLength, err := ch.Uint8()
	if err != nil {
		return ErrNotEnoughBytes
	}
	n := 1

	pkg.Name, err = ch.String(int(nameLength))
	if err != nil {
		return ErrNotEnoughBytes
	}
	n += len(pkg.Name)

	if pkg.wide {
		var opts uint32
		opts, err = ch.Uint32()
		pkg.Options = CursorOption(opts)
		n += 4
	} else {
		var opts uint8
		opts, err = ch.Uint8()
		pkg.Options = CursorOption(opts)
		n++
	}
	if err != nil {
		return ErrNotEnoughBytes
	}

	status, err := ch.Uint8()
	if err != nil {
		return ErrNotEnoughBytes
	}
	pkg.Status = CursorDStatus(status)
	n++

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
	n += len(pkg.Stmt)

	columnCount, err := ch.Uint16()
	if err != nil {
		return ErrNotEnoughBytes
	}
	n += 2

	for i := 0; i < int(columnCount); i++ {
		nameLength, err := ch.Uint8()
		if err != nil {
			return ErrNotEnoughBytes
		}
		n++

		name, err := ch.String(int(nameLength))
		if err != nil {
			return ErrNotEnoughBytes
		}
		n += int(nameLength)

		pkg.columns = append(pkg.columns, name)
	}

	if n != totalLength {
		return fmt.Errorf("expected to read %d bytes, read %d bytes instead", totalLength, n)
	}

	return nil
}

// WriteTo implements the tds.Package interface.
func (pkg CurDeclarePackage) WriteTo(ch BytesChannel) error {
	token := TDS_CURDECLARE
	if pkg.wide {
		token = TDS_CURDECLARE3
	}
	if err := ch.WriteByte(byte(token)); err != nil {
		return err
	}

	// 1 name length
	// x name
	// 1 or 4 options
	// 1 status
	// 2 or 4 stmt len
	// x stmt
	// 2 columns len
	// x columns (1 name len, x name per col)
	totalLength := 1 + len(pkg.Name) + 1 + 1 + 2 + len(pkg.Stmt) + 2
	if pkg.wide {
		// +3 options
		// +2 stmt len
		totalLength += 3 + 2
	}
	for _, col := range pkg.columns {
		totalLength += 1 + len(col)
	}

	var err error
	if pkg.wide {
		err = ch.WriteUint32(uint32(totalLength))
	} else {
		err = ch.WriteUint16(uint16(totalLength))
	}
	if err != nil {
		return err
	}

	if err := ch.WriteUint8(uint8(len(pkg.Name))); err != nil {
		return err
	}
	n := 1

	if err := ch.WriteString(pkg.Name); err != nil {
		return err
	}
	n += len(pkg.Name)

	if pkg.wide {
		if err := ch.WriteUint32(uint32(pkg.Options)); err != nil {
			return err
		}
		n += 4
	} else {
		if err := ch.WriteUint8(uint8(pkg.Options)); err != nil {
			return err
		}
		n++
	}

	if err := ch.WriteUint8(uint8(pkg.Status)); err != nil {
		return err
	}
	n++

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

	if err := ch.WriteUint16(uint16(len(pkg.columns))); err != nil {
		return err
	}
	n += 2

	for _, name := range pkg.columns {
		if err := ch.WriteUint8(uint8(len(name))); err != nil {
			return err
		}
		n++

		if err := ch.WriteString(name); err != nil {
			return err
		}
		n += len(name)
	}

	if n != int(totalLength) {
		return fmt.Errorf("expected to write %d bytes, wrote %d bytes instead",
			totalLength, n)
	}

	return nil
}

func (pkg CurDeclarePackage) String() string {
	strOpts := deBitmaskString(int(pkg.Options), int(TDS_CUR_DOPT_RELLOCKSONCLOSE),
		func(i int) string { return CursorOption(i).String() },
		TDS_CUR_DOPT_UNUSED.String(),
	)

	wide := "nowide"
	if pkg.wide {
		wide = "wide"
	}

	return fmt.Sprintf("%T(%s, %s, %s, %s, '%s')", pkg, wide, pkg.Name, strOpts, pkg.Status, pkg.Stmt)
}
