// SPDX-FileCopyrightText: 2020-2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by "stringer -type=CursorOption"; DO NOT EDIT.

package tds

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TDS_CUR_DOPT_UNUSED-0]
	_ = x[TDS_CUR_DOPT_RDONLY-1]
	_ = x[TDS_CUR_DOPT_UPDATABLE-2]
	_ = x[TDS_CUR_DOPT_SENSITIVE-4]
	_ = x[TDS_CUR_DOPT_DYNAMIC-8]
	_ = x[TDS_CUR_DOPT_IMPLICIT-16]
	_ = x[TDS_CUR_DOPT_INSENSITIVE-32]
	_ = x[TDS_CUR_DOPT_SEMISENSITIVE-64]
	_ = x[TDS_CUR_DOPT_KEYSETDRIVEN-128]
	_ = x[TDS_CUR_DOPT_SCROLLABLE-256]
	_ = x[TDS_CUR_DOPT_RELLOCKSONCLOSE-512]
}

const (
	_CursorOption_name_0 = "TDS_CUR_DOPT_UNUSEDTDS_CUR_DOPT_RDONLYTDS_CUR_DOPT_UPDATABLE"
	_CursorOption_name_1 = "TDS_CUR_DOPT_SENSITIVE"
	_CursorOption_name_2 = "TDS_CUR_DOPT_DYNAMIC"
	_CursorOption_name_3 = "TDS_CUR_DOPT_IMPLICIT"
	_CursorOption_name_4 = "TDS_CUR_DOPT_INSENSITIVE"
	_CursorOption_name_5 = "TDS_CUR_DOPT_SEMISENSITIVE"
	_CursorOption_name_6 = "TDS_CUR_DOPT_KEYSETDRIVEN"
	_CursorOption_name_7 = "TDS_CUR_DOPT_SCROLLABLE"
	_CursorOption_name_8 = "TDS_CUR_DOPT_RELLOCKSONCLOSE"
)

var (
	_CursorOption_index_0 = [...]uint8{0, 19, 38, 60}
)

func (i CursorOption) String() string {
	switch {
	case i <= 2:
		return _CursorOption_name_0[_CursorOption_index_0[i]:_CursorOption_index_0[i+1]]
	case i == 4:
		return _CursorOption_name_1
	case i == 8:
		return _CursorOption_name_2
	case i == 16:
		return _CursorOption_name_3
	case i == 32:
		return _CursorOption_name_4
	case i == 64:
		return _CursorOption_name_5
	case i == 128:
		return _CursorOption_name_6
	case i == 256:
		return _CursorOption_name_7
	case i == 512:
		return _CursorOption_name_8
	default:
		return "CursorOption(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
