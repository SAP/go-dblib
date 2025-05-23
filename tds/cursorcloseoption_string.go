// SPDX-FileCopyrightText: 2020-2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by "stringer -type=CursorCloseOption"; DO NOT EDIT.

package tds

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TDS_CUR_COPT_UNUSED-0]
	_ = x[TDS_CUR_COPT_DEALLOC-1]
}

const _CursorCloseOption_name = "TDS_CUR_COPT_UNUSEDTDS_CUR_COPT_DEALLOC"

var _CursorCloseOption_index = [...]uint8{0, 19, 39}

func (i CursorCloseOption) String() string {
	if i >= CursorCloseOption(len(_CursorCloseOption_index)-1) {
		return "CursorCloseOption(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _CursorCloseOption_name[_CursorCloseOption_index[i]:_CursorCloseOption_index[i+1]]
}
