// SPDX-FileCopyrightText: 2020-2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by "gen_type NullNChar sql.NullString -columndef nchar(13) null -compare compareNullChar"; DO NOT EDIT.

package integration

import (
	"database/sql"

	"testing"
)

// DoTestNullNChar tests the handling of the NullNChar.
func DoTestNullNChar(t *testing.T) {
	TestForEachDB("TestNullNChar", t, testNullNChar)
	//
}

func testNullNChar(t *testing.T, db *sql.DB, tableName string) {
	// insert is the amount of insertions (see fn SetupTableInsert)
	insert := 2

	pass := make([]interface{}, len(samplesNullNChar))
	mySamples := make([]sql.NullString, len(samplesNullNChar)*insert)

	for i, sample := range samplesNullNChar {

		mySample := sample

		pass[i] = mySample

		// Add passed sample for the later validation (for every
		// insert)
		for j := 0; j < insert; j++ {
			mySamples[i+(len(samplesNullNChar)*j)] = mySample
		}
	}

	rows, teardownFn, err := SetupTableInsert(db, tableName, "nchar(13) null", pass...)
	if err != nil {
		t.Errorf("Error preparing table: %v", err)
		return
	}
	defer rows.Close()
	defer teardownFn()

	i := 0
	var recv sql.NullString
	for rows.Next() {
		if err := rows.Scan(&recv); err != nil {
			t.Errorf("Scan failed on %dth scan: %v", i, err)
			continue
		}

		if compareNullChar(recv, mySamples[i]) {

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
