// SPDX-FileCopyrightText: 2020-2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by "stringer -type=PacketHeaderStatus"; DO NOT EDIT.

package tds

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TDS_BUFSTAT_EOM-1]
	_ = x[TDS_BUFSTAT_ATTNACK-2]
	_ = x[TDS_BUFSTAT_ATTN-4]
	_ = x[TDS_BUFSTAT_EVENT-8]
	_ = x[TDS_BUFSTAT_SEAL-16]
	_ = x[TDS_BUFSTAT_ENCRYPT-32]
	_ = x[TDS_BUFSTAT_SYMENCRYPT-64]
}

const (
	_PacketHeaderStatus_name_0 = "TDS_BUFSTAT_EOMTDS_BUFSTAT_ATTNACK"
	_PacketHeaderStatus_name_1 = "TDS_BUFSTAT_ATTN"
	_PacketHeaderStatus_name_2 = "TDS_BUFSTAT_EVENT"
	_PacketHeaderStatus_name_3 = "TDS_BUFSTAT_SEAL"
	_PacketHeaderStatus_name_4 = "TDS_BUFSTAT_ENCRYPT"
	_PacketHeaderStatus_name_5 = "TDS_BUFSTAT_SYMENCRYPT"
)

var (
	_PacketHeaderStatus_index_0 = [...]uint8{0, 15, 34}
)

func (i PacketHeaderStatus) String() string {
	switch {
	case 1 <= i && i <= 2:
		i -= 1
		return _PacketHeaderStatus_name_0[_PacketHeaderStatus_index_0[i]:_PacketHeaderStatus_index_0[i+1]]
	case i == 4:
		return _PacketHeaderStatus_name_1
	case i == 8:
		return _PacketHeaderStatus_name_2
	case i == 16:
		return _PacketHeaderStatus_name_3
	case i == 32:
		return _PacketHeaderStatus_name_4
	case i == 64:
		return _PacketHeaderStatus_name_5
	default:
		return "PacketHeaderStatus(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
