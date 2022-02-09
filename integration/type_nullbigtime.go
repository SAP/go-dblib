// SPDX-FileCopyrightText: 2022 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by "gen_type NullBigTime sql.NullTime -columndef bigtime null"; DO NOT EDIT.

package integration

import (
	"database/sql"

	"testing"
)

// DoTestNullBigTime tests the handling of the NullBigTime.
func DoTestNullBigTime(t *testing.T) {
	TestForEachDB("TestNullBigTime", t, testNullBigTime)
	//
}

func testNullBigTime(t *testing.T, db *sql.DB, tableName string) {
	// insert is the amount of insertions (see fn SetupTableInsert)
	insert := 2

	pass := make([]interface{}, len(samplesNullBigTime))
	mySamples := make([]sql.NullTime, len(samplesNullBigTime)*insert)

	for i, sample := range samplesNullBigTime {

		mySample := sample

		pass[i] = mySample

		// Add passed sample for the later validation (for every
		// insert)
		for j := 0; j < insert; j++ {
			mySamples[i+(len(samplesNullBigTime)*j)] = mySample
		}
	}

	rows, teardownFn, err := SetupTableInsert(db, tableName, "bigtime null", pass...)
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
