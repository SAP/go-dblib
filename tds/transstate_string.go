// SPDX-FileCopyrightText: 2020-2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by "stringer -type=TransState"; DO NOT EDIT.

package tds

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TDS_NOT_IN_TRAN-0]
	_ = x[TDS_TRAN_IN_PROGRESS-1]
	_ = x[TDS_TRAN_COMPLETED-2]
	_ = x[TDS_TRAN_FAIL-3]
	_ = x[TDS_TRAN_STMT_FAIL-4]
}

const _TransState_name = "TDS_NOT_IN_TRANTDS_TRAN_IN_PROGRESSTDS_TRAN_COMPLETEDTDS_TRAN_FAILTDS_TRAN_STMT_FAIL"

var _TransState_index = [...]uint8{0, 15, 35, 53, 66, 84}

func (i TransState) String() string {
	if i >= TransState(len(_TransState_index)-1) {
		return "TransState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TransState_name[_TransState_index[i]:_TransState_index[i+1]]
}
