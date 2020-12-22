// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import (
	"fmt"
)

//go:generate stringer -type=DoneState

// DoneState is the type for bitmask values of a done-package state.
type DoneState uint16

// States of a done-package.
const (
	TDS_DONE_FINAL      DoneState = 0x0
	TDS_DONE_MORE       DoneState = 0x1
	TDS_DONE_ERROR      DoneState = 0x2
	TDS_DONE_INXACT     DoneState = 0x4
	TDS_DONE_PROC       DoneState = 0x8
	TDS_DONE_COUNT      DoneState = 0x10
	TDS_DONE_ATTN       DoneState = 0x20
	TDS_DONE_EVENT      DoneState = 0x40
	TDS_DONE_CUMULATIVE DoneState = 0x80
)

//go:generate stringer -type=TransState

// TransState is the type for bitmask values of a trans-package state.
type TransState uint16

// States of a trans-package.
const (
	TDS_NOT_IN_TRAN      TransState = 0x0
	TDS_TRAN_IN_PROGRESS TransState = 0x1
	TDS_TRAN_COMPLETED   TransState = 0x2
	TDS_TRAN_FAIL        TransState = 0x3
	TDS_TRAN_STMT_FAIL   TransState = 0x4
)

// DonePackage is used to finalize a package stream.
type DonePackage struct {
	Status    DoneState
	TranState TransState
	Count     int32
}

type DoneProcPackage = DonePackage
type DoneInProcPackage = DonePackage

// ReadFrom implements the tds.Package interface.
func (pkg *DonePackage) ReadFrom(ch BytesChannel) error {
	status, err := ch.Uint16()
	if err != nil {
		return ErrNotEnoughBytes
	}
	pkg.Status = DoneState(status)

	tranState, err := ch.Uint16()
	if err != nil {
		return ErrNotEnoughBytes
	}
	pkg.TranState = TransState(tranState)

	pkg.Count, err = ch.Int32()
	if err != nil {
		return ErrNotEnoughBytes
	}

	return nil
}

// WriteTo implements the tds.Package interface.
func (pkg DonePackage) WriteTo(ch BytesChannel) error {
	if err := ch.WriteByte(byte(TDS_DONE)); err != nil {
		return err
	}

	if err := ch.WriteUint16(uint16(pkg.Status)); err != nil {
		return err
	}

	if err := ch.WriteUint16(uint16(pkg.TranState)); err != nil {
		return err
	}

	return ch.WriteInt32(pkg.Count)
}

// String implements the fmt.Stringer interface.
func (pkg DonePackage) String() string {
	strStati := deBitmaskString(int(pkg.Status), int(TDS_DONE_CUMULATIVE),
		func(i int) string { return DoneState(i).String() },
		TDS_DONE_FINAL.String(),
	)

	strTransi := deBitmaskString(int(pkg.TranState), int(TDS_TRAN_STMT_FAIL),
		func(i int) string { return TransState(i).String() },
		TDS_NOT_IN_TRAN.String(),
	)

	return fmt.Sprintf("%T(%s, %s, Count=%d)", pkg, strStati, strTransi, pkg.Count)
}
