// SPDX-FileCopyrightText: 2022 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by "gen_type NullDateTime sql.NullTime -columndef datetime null"; DO NOT EDIT.

package integration

import (
	"database/sql"

	"testing"
)

// DoTestNullDateTime tests the handling of the NullDateTime.
func DoTestNullDateTime(t *testing.T) {
	TestForEachDB("TestNullDateTime", t, testNullDateTime)
	//
}

func testNullDateTime(t *testing.T, db *sql.DB, tableName string) {
	// insert is the amount of insertions (see fn SetupTableInsert)
	insert := 2

	pass := make([]interface{}, len(samplesNullDateTime))
	mySamples := make([]sql.NullTime, len(samplesNullDateTime)*insert)

	for i, sample := range samplesNullDateTime {

		mySample := sample

		pass[i] = mySample

		// Add passed sample for the later validation (for every
		// insert)
		for j := 0; j < insert; j++ {
			mySamples[i+(len(samplesNullDateTime)*j)] = mySample
		}
	}

	rows, teardownFn, err := SetupTableInsert(db, tableName, "datetime null", pass...)
	if err != nil {
		t.Errorf("Error preparing table: %v", err)
		return
	}
	defer rows.Close()
	defer teardownFn()

	i := 0
	var recv sql.NullTime
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
