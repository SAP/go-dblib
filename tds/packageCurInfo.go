// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import "fmt"

//go:generate stringer -type=CursorCommand

// CursorCommand is the type used to send commands on cursors.
type CursorCommand uint

const (
	TDS_CUR_CMD_SETCURROWS CursorCommand = iota + 1
	TDS_CUR_CMD_INQUIRE
	TDS_CUR_CMD_INFORM
	TDS_CUR_CMD_LISTALL
)

//go:generate stringer -type=CursorIStatus

// CursorIStatus is the type for a bitmask indicating the status of
// a cursor.
type CursorIStatus uint

const (
	TDS_CUR_ISTAT_UNUSED          CursorIStatus = 0x0
	TDS_CUR_ISTAT_DECLARED        CursorIStatus = 0x1
	TDS_CUR_ISTAT_OPEN            CursorIStatus = 0x2
	TDS_CUR_ISTAT_CLOSED          CursorIStatus = 0x4
	TDS_CUR_ISTAT_RDONLY          CursorIStatus = 0x8
	TDS_CUR_ISTAT_UPDATABLE       CursorIStatus = 0x10
	TDS_CUR_ISTAT_ROWCNT          CursorIStatus = 0x20
	TDS_CUR_ISTAT_DEALLOC         CursorIStatus = 0x40
	TDS_CUR_ISTAT_SCROLLABLE      CursorIStatus = 0x80
	TDS_CUR_ISTAT_IMPLICIT        CursorIStatus = 0x100
	TDS_CUR_ISTAT_SENSITIVE       CursorIStatus = 0x200
	TDS_CUR_ISTAT_INSENSITIVE     CursorIStatus = 0x400
	TDS_CUR_ISTAT_SEMISENSITIVE   CursorIStatus = 0x800
	TDS_CUR_ISTAT_KEYSETDRIVEN    CursorIStatus = 0x1000
	TDS_CUR_ISTAT_RELLOCKSONCLOSE CursorIStatus = 0x2000
)

var _ Package = (*CurInfoPackage)(nil)

// CurInfoPackage is used to send and receive information about
// a cursor.
type CurInfoPackage struct {
	CursorID int32
	Name     string

	Command CursorCommand
	Status  CursorIStatus

	RowNum    int32
	TotalRows int32
	RowCount  int32

	wide bool
}

// ReadFrom implements the tds.Package interface.
func (pkg *CurInfoPackage) ReadFrom(ch BytesChannel) error {
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

	command, err := ch.Uint8()
	if err != nil {
		return ErrNotEnoughBytes
	}
	pkg.Command = CursorCommand(command)
	n++

	if pkg.wide {
		status, err := ch.Uint32()
		if err != nil {
			return ErrNotEnoughBytes
		}
		pkg.Status = CursorIStatus(status)
		n += 4
	} else {
		status, err := ch.Uint16()
		if err != nil {
			return ErrNotEnoughBytes
		}
		pkg.Status = CursorIStatus(status)
		n += 2
	}

	if pkg.wide {
		pkg.RowNum, err = ch.Int32()
		if err != nil {
			return ErrNotEnoughBytes
		}
		n += 4

		pkg.TotalRows, err = ch.Int32()
		if err != nil {
			return ErrNotEnoughBytes
		}
		n += 4
	}

	if pkg.Status&TDS_CUR_ISTAT_ROWCNT == TDS_CUR_ISTAT_ROWCNT {
		pkg.RowCount, err = ch.Int32()
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
func (pkg CurInfoPackage) WriteTo(ch BytesChannel) error {
	token := TDS_CURINFO
	if pkg.wide {
		token = TDS_CURINFO3
	}
	if err := ch.WriteByte(byte(token)); err != nil {
		return err
	}

	// 4 cursorID
	// if cursorID == 0:
	//   1 name len
	//   x name
	// 1 command
	// 2/4 status
	// if wide:
	//   4 rownum
	//   4 totalrows
	// if status ROWCNT:
	//   4 rowcnt
	totalLength := 4 + 1 + 2
	if pkg.CursorID == 0 {
		totalLength += 1 + len(pkg.Name)
	}
	if pkg.Status&TDS_CUR_ISTAT_ROWCNT == TDS_CUR_ISTAT_ROWCNT {
		totalLength += 4
	}
	if pkg.wide {
		// +2 status
		// 4 rownum
		// 4 totalrows
		totalLength += 2 + 4 + 4
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

	if err := ch.WriteUint8(uint8(pkg.Command)); err != nil {
		return err
	}

	if pkg.wide {
		if err := ch.WriteUint32(uint32(pkg.Status)); err != nil {
			return err
		}
	} else {
		if err := ch.WriteUint16(uint16(pkg.Status)); err != nil {
			return err
		}
	}

	if pkg.wide {
		if err := ch.WriteInt32(pkg.RowNum); err != nil {
			return err
		}

		if err := ch.WriteInt32(pkg.TotalRows); err != nil {
			return err
		}
	}

	if pkg.Status&TDS_CUR_ISTAT_ROWCNT == TDS_CUR_ISTAT_ROWCNT {
		if err := ch.WriteInt32(pkg.RowCount); err != nil {
			return err
		}
	}

	return nil
}

// String implements the fmt.Stringer interface.
func (pkg CurInfoPackage) String() string {
	wide := "nowide"
	if pkg.wide {
		wide = "wide"
	}

	strStati := deBitmaskString(int(pkg.Status), int(TDS_CUR_ISTAT_RELLOCKSONCLOSE),
		func(i int) string { return CursorIStatus(i).String() },
		TDS_CUR_ISTAT_UNUSED.String(),
	)

	return fmt.Sprintf("%T(%s, %d, %s, %s, %s, RowNum=%d, TotalRows=%d, RowCount=%d)", pkg, wide, pkg.CursorID,
		pkg.Name, pkg.Command, strStati, pkg.RowNum, pkg.TotalRows, pkg.RowCount)
}
