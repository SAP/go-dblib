// SPDX-FileCopyrightText: 2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import "fmt"

//go:generate stringer -type=TDSMsgStatus

// TDSMsgStatus is the type for bitmask values of a TDS-message status.
type TDSMsgStatus uint8

// States of a TDS-message.
const (
	TDS_MSG_HASNOARGS TDSMsgStatus = iota
	TDS_MSG_HASARGS
)

//go:generate stringer -type=TDSMsgId

// TDSMsgId is the type for bitmask values of a TDS-message id.
type TDSMsgId uint16

// Ids of a TDS-message.
const (
	TDS_MSG_SEC_ENCRYPT TDSMsgId = iota + 1
	TDS_MSG_SEC_LOGPWD
	TDS_MSG_SEC_REMPWD
	TDS_MSG_SEC_CHALLENGE
	TDS_MSG_SEC_RESPONSE
	TDS_MSG_SEC_GETLABEL
	TDS_MSG_SEC_LABEL
	TDS_MSG_SQL_TBLNAME
	TDS_MSG_GW_RESERVED
	TDS_MSG_OMNI_CAPABILITIES
	TDS_MSG_SEC_OPAQUE
	TDS_MSG_HAFAILOVER
	TDS_MSG_EMPTY
	TDS_MSG_SEC_ENCRYPT2
	TDS_MSG_SEC_LOGPWD2
	TDS_MSG_SEC_SUP_CIPHER2
	TDS_MSG_MIG_REQ
	TDS_MSG_MIG_SYNC
	TDS_MSG_MIG_CONT
	TDS_MSG_MIG_IGN
	TDS_MSG_MIG_FAIL
	TDS_MSG_SEC_REMPWD2
	TDS_MSG_MIG_RESUME
	TDS_MSG_HELLO
	TDS_MSG_LOGINPARAMS
	TDS_MSG_GRID_MIGREQ
	TDS_MSG_GRID_QUIESCE
	TDS_MSG_GRID_UNQUIESCE
	TDS_MSG_GRID_EVENT
	TDS_MSG_SEC_ENCRYPT3
	TDS_MSG_SEC_LOGPWD3
	TDS_MSG_SEC_REMPWD3
	TDS_MSG_DR_MAP
	TDS_MSG_SEC_SYMKEY
	TDS_MSG_SEC_ENCRYPT4
)

// TDSOpaqueSecurityToken is the type for bitmask values of
// a TDS-message security-token.
type TDSOpaqueSecurityToken uint8

// Security-tokens of a TDS-message.
const (
	/*
	 ** TDS_MSG_SEC_OPAQUE message types.
	 */
	TDS_SEC_SECSESS TDSOpaqueSecurityToken = iota
	TDS_SEC_FORWARD
	TDS_SEC_SIGN
	TDS_SEC_OTHER
)

// MsgPackage contains the status and id of a TDS-message.
type MsgPackage struct {
	Status TDSMsgStatus
	MsgId  TDSMsgId
}

// NewMsgPackage creates a TDS-message with status and id.
func NewMsgPackage(status TDSMsgStatus, msgId TDSMsgId) *MsgPackage {
	return &MsgPackage{
		Status: status,
		MsgId:  msgId,
	}
}

// ReadFrom implements the tds.Package interface.
func (pkg *MsgPackage) ReadFrom(ch BytesChannel) error {
	var err error

	if _, err := ch.Uint8(); err != nil {
		return ErrNotEnoughBytes
	}

	var status uint8
	status, err = ch.Uint8()
	if err != nil {
		return ErrNotEnoughBytes
	}
	pkg.Status = (TDSMsgStatus)(status)

	msgId, err := ch.Uint16()
	if err != nil {
		return ErrNotEnoughBytes
	}
	pkg.MsgId = TDSMsgId(msgId)

	return nil
}

// Write to implements the tds.Package interface.
func (pkg MsgPackage) WriteTo(ch BytesChannel) error {
	if err := ch.WriteByte(byte(TDS_MSG)); err != nil {
		return err
	}

	// Length
	if err := ch.WriteUint8(3); err != nil {
		return err
	}

	if err := ch.WriteUint8(uint8(pkg.Status)); err != nil {
		return err
	}

	return ch.WriteUint16(uint16(pkg.MsgId))
}

func (pkg MsgPackage) String() string {
	return fmt.Sprintf("%T(%s, %s)", pkg, pkg.Status, pkg.MsgId)
}
