// SPDX-FileCopyrightText: 2022 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by "gen_type VarBinary []byte -columndef varbinary(13) -compare compareBinary"; DO NOT EDIT.

package integration

import (
	"database/sql"

	"testing"
)

// DoTestVarBinary tests the handling of the VarBinary.
func DoTestVarBinary(t *testing.T) {
	TestForEachDB("TestVarBinary", t, testVarBinary)
	//
}

func testVarBinary(t *testing.T, db *sql.DB, tableName string) {
	pass := make([]interface{}, len(samplesVarBinary))
	mySamples := make([][]byte, len(samplesVarBinary))

	for i, sample := range samplesVarBinary {

		mySample := sample

		pass[i] = mySample
		mySamples[i] = mySample
	}

	rows, teardownFn, err := SetupTableInsert(db, tableName, "varbinary(13)", pass...)
	if err != nil {
		t.Errorf("Error preparing table: %v", err)
		return
	}
	defer rows.Close()
	defer teardownFn()

	i := 0
	var recv []byte
	for rows.Next() {
		if err := rows.Scan(&recv); err != nil {
			t.Errorf("Scan failed on %dth scan: %v", i, err)
			continue
		}

		if compareBinary(recv, mySamples[i]) {

			t.Errorf("Received value does not match passed parameter")
			t.Errorf("Expected: %v", mySamples[i])
			t.Errorf("Received: %v", recv)
		}

		i++
	}

	if err := rows.Err(); err != nil {
		t.Errorf("Error preparing rows: %v", err)
	}

	if i != len(pass) {
		t.Errorf("Only read %d values from database, expected to read %d", i, len(pass))
	}
}
