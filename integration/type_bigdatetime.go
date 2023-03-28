// SPDX-FileCopyrightText: 2023 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by "gen_type BigDateTime time.Time"; DO NOT EDIT.

package integration

import (
	"database/sql"

	"testing"

	"time"
)

// DoTestBigDateTime tests the handling of the BigDateTime.
func DoTestBigDateTime(t *testing.T) {
	TestForEachDB("TestBigDateTime", t, testBigDateTime)
	//
}

func testBigDateTime(t *testing.T, db *sql.DB, tableName string) {
	// insert is the amount of insertions (see fn SetupTableInsert)
	insert := 2

	pass := make([]interface{}, len(samplesBigDateTime))
	mySamples := make([]time.Time, len(samplesBigDateTime)*insert)

	for i, sample := range samplesBigDateTime {

		mySample := sample

		pass[i] = mySample

		// Add passed sample for the later validation (for every
		// insert)
		for j := 0; j < insert; j++ {
			mySamples[i+(len(samplesBigDateTime)*j)] = mySample
		}
	}

	rows, teardownFn, err := SetupTableInsert(db, tableName, "bigdatetime", pass...)
	if err != nil {
		t.Errorf("Error preparing table: %v", err)
		return
	}
	defer rows.Close()
	defer teardownFn()

	i := 0
	var recv time.Time
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
