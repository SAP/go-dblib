// SPDX-FileCopyrightText: 2020-2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by "gen_type NVarChar string -columndef nvarchar(13) -compare compareVarChar"; DO NOT EDIT.

package integration

import (
	"database/sql"

	"testing"
)

// DoTestNVarChar tests the handling of the NVarChar.
func DoTestNVarChar(t *testing.T) {
	TestForEachDB("TestNVarChar", t, testNVarChar)
	//
}

func testNVarChar(t *testing.T, db *sql.DB, tableName string) {
	// insert is the amount of insertions (see fn SetupTableInsert)
	insert := 2

	pass := make([]interface{}, len(samplesNVarChar))
	mySamples := make([]string, len(samplesNVarChar)*insert)

	for i, sample := range samplesNVarChar {

		mySample := sample

		pass[i] = mySample

		// Add passed sample for the later validation (for every
		// insert)
		for j := 0; j < insert; j++ {
			mySamples[i+(len(samplesNVarChar)*j)] = mySample
		}
	}

	rows, teardownFn, err := SetupTableInsert(db, tableName, "nvarchar(13)", pass...)
	if err != nil {
		t.Errorf("Error preparing table: %v", err)
		return
	}
	defer rows.Close()
	defer teardownFn()

	i := 0
	var recv string
	for rows.Next() {
		if err := rows.Scan(&recv); err != nil {
			t.Errorf("Scan failed on %dth scan: %v", i, err)
			continue
		}

		if compareVarChar(recv, mySamples[i]) {

			t.Errorf("Received value does not match passed parameter")
			t.Errorf("Expected: %v", mySamples[i])
			t.Errorf("Received: %v", recv)
		}

		i++
	}

	if err := rows.Err(); err != nil {
		t.Errorf("Error preparing rows: %v", err)
	}

	if i != len(pass)*insert {
		t.Errorf("Only read %d values from database, expected to read %d", i, len(pass)*insert)
	}
}
