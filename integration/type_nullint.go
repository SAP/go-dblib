// SPDX-FileCopyrightText: 2023 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by "gen_type NullInt sql.NullInt32 -columndef int null"; DO NOT EDIT.

package integration

import (
	"database/sql"

	"testing"
)

// DoTestNullInt tests the handling of the NullInt.
func DoTestNullInt(t *testing.T) {
	TestForEachDB("TestNullInt", t, testNullInt)
	//
}

func testNullInt(t *testing.T, db *sql.DB, tableName string) {
	// insert is the amount of insertions (see fn SetupTableInsert)
	insert := 2

	pass := make([]interface{}, len(samplesNullInt))
	mySamples := make([]sql.NullInt32, len(samplesNullInt)*insert)

	for i, sample := range samplesNullInt {

		mySample := sample

		pass[i] = mySample

		// Add passed sample for the later validation (for every
		// insert)
		for j := 0; j < insert; j++ {
			mySamples[i+(len(samplesNullInt)*j)] = mySample
		}
	}

	rows, teardownFn, err := SetupTableInsert(db, tableName, "int null", pass...)
	if err != nil {
		t.Errorf("Error preparing table: %v", err)
		return
	}
	defer rows.Close()
	defer teardownFn()

	i := 0
	var recv sql.NullInt32
	for rows.Next() {
		if err := rows.Scan(&recv); err != nil {
			t.Errorf("Scan failed on %dth scan: %v", i, err)
			continue
		}

		if recv != mySamples[i] {

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
