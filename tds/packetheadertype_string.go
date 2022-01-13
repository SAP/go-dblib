// SPDX-FileCopyrightText: 2022 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by "stringer -type=PacketHeaderType"; DO NOT EDIT.

package tds

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TDS_BUF_LANG-1]
	_ = x[TDS_BUF_LOGIN-2]
	_ = x[TDS_BUF_RPC-3]
	_ = x[TDS_BUF_RESPONSE-4]
	_ = x[TDS_BUF_UNFMT-5]
	_ = x[TDS_BUF_ATTN-6]
	_ = x[TDS_BUF_BULK-7]
	_ = x[TDS_BUF_SETUP-8]
	_ = x[TDS_BUF_CLOSE-9]
	_ = x[TDS_BUF_ERROR-10]
	_ = x[TDS_BUF_PROTACK-11]
	_ = x[TDS_BUF_ECHO-12]
	_ = x[TDS_BUF_LOGOUT-13]
	_ = x[TDS_BUF_ENDPARAM-14]
	_ = x[TDS_BUF_NORMAL-15]
	_ = x[TDS_BUF_URGENT-16]
	_ = x[TDS_BUF_MIGRATE-17]
	_ = x[TDS_BUF_HELLO-18]
	_ = x[TDS_BUF_CMDSEQ_NORMAL-19]
	_ = x[TDS_BUF_CMDSEQ_LOGIN-20]
	_ = x[TDS_BUF_CMDSEQ_LIVENESS-21]
	_ = x[TDS_BUF_CMDSEQ_RESERVED1-22]
	_ = x[TDS_BUF_CMDSEQ_RESERVED2-23]
}

const _PacketHeaderType_name = "TDS_BUF_LANGTDS_BUF_LOGINTDS_BUF_RPCTDS_BUF_RESPONSETDS_BUF_UNFMTTDS_BUF_ATTNTDS_BUF_BULKTDS_BUF_SETUPTDS_BUF_CLOSETDS_BUF_ERRORTDS_BUF_PROTACKTDS_BUF_ECHOTDS_BUF_LOGOUTTDS_BUF_ENDPARAMTDS_BUF_NORMALTDS_BUF_URGENTTDS_BUF_MIGRATETDS_BUF_HELLOTDS_BUF_CMDSEQ_NORMALTDS_BUF_CMDSEQ_LOGINTDS_BUF_CMDSEQ_LIVENESSTDS_BUF_CMDSEQ_RESERVED1TDS_BUF_CMDSEQ_RESERVED2"

var _PacketHeaderType_index = [...]uint16{0, 12, 25, 36, 52, 65, 77, 89, 102, 115, 128, 143, 155, 169, 185, 199, 213, 228, 241, 262, 282, 305, 329, 353}

func (i PacketHeaderType) String() string {
	i -= 1
	if i >= PacketHeaderType(len(_PacketHeaderType_index)-1) {
		return "PacketHeaderType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _PacketHeaderType_name[_PacketHeaderType_index[i]:_PacketHeaderType_index[i+1]]
}
