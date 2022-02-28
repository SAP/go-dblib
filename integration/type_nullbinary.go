// SPDX-FileCopyrightText: 2022 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by "gen_type NullBinary github.com/SAP/go-dblib/asetypes.NullBinary -columndef binary(13) null -compare compareNullBinary"; DO NOT EDIT.

package integration

import (
	"database/sql"

	"github.com/SAP/go-dblib/asetypes"

	"testing"
)

// DoTestNullBinary tests the handling of the NullBinary.
func DoTestNullBinary(t *testing.T) {
	TestForEachDB("TestNullBinary", t, testNullBinary)
	//
}

func testNullBinary(t *testing.T, db *sql.DB, tableName string) {
	// insert is the amount of insertions (see fn SetupTableInsert)
	insert := 2

	pass := make([]interface{}, len(samplesNullBinary))
	mySamples := make([]asetypes.NullBinary, len(samplesNullBinary)*insert)

	for i, sample := range samplesNullBinary {

		mySample := sample

		pass[i] = mySample

		// Add passed sample for the later validation (for every
		// insert)
		for j := 0; j < insert; j++ {
			mySamples[i+(len(samplesNullBinary)*j)] = mySample
		}
	}

	rows, teardownFn, err := SetupTableInsert(db, tableName, "binary(13) null", pass...)
	if err != nil {
		t.Errorf("Error preparing table: %v", err)
		return
	}
	defer rows.Close()
	defer teardownFn()

	i := 0
	var recv asetypes.NullBinary
	for rows.Next() {
		if err := rows.Scan(&recv); err != nil {
			t.Errorf("Scan failed on %dth scan: %v", i, err)
			continue
		}

		if compareNullBinary(recv, mySamples[i]) {

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
