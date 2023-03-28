// SPDX-FileCopyrightText: 2023 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by "stringer -type=CursorCommand"; DO NOT EDIT.

package tds

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TDS_CUR_CMD_SETCURROWS-1]
	_ = x[TDS_CUR_CMD_INQUIRE-2]
	_ = x[TDS_CUR_CMD_INFORM-3]
	_ = x[TDS_CUR_CMD_LISTALL-4]
}

const _CursorCommand_name = "TDS_CUR_CMD_SETCURROWSTDS_CUR_CMD_INQUIRETDS_CUR_CMD_INFORMTDS_CUR_CMD_LISTALL"

var _CursorCommand_index = [...]uint8{0, 22, 41, 59, 78}

func (i CursorCommand) String() string {
	i -= 1
	if i >= CursorCommand(len(_CursorCommand_index)-1) {
		return "CursorCommand(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _CursorCommand_name[_CursorCommand_index[i]:_CursorCommand_index[i+1]]
}
