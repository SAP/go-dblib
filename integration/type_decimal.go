// SPDX-FileCopyrightText: 2021 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by "gen_type Decimal github.com/SAP/go-dblib/*asetypes.Decimal -columndef decimal(38,19) -convert convertDecimal3819 -compare compareDecimal"; DO NOT EDIT.

package integration

import (
	"database/sql"

	"github.com/SAP/go-dblib/asetypes"

	"testing"
)

// DoTestDecimal tests the handling of the Decimal.
func DoTestDecimal(t *testing.T) {
	TestForEachDB("TestDecimal", t, testDecimal)
	//
}

func testDecimal(t *testing.T, db *sql.DB, tableName string) {
	pass := make([]interface{}, len(samplesDecimal))
	mySamples := make([]*asetypes.Decimal, len(samplesDecimal))

	for i, sample := range samplesDecimal {

		// Convert sample with passed function before proceeding
		mySample, err := convertDecimal3819(sample)
		if err != nil {
			t.Errorf("Failed to convert sample %v: %v", sample, err)
			return
		}

		pass[i] = mySample
		mySamples[i] = mySample
	}

	rows, teardownFn, err := SetupTableInsert(db, tableName, "decimal(38,19)", pass...)
	if err != nil {
		t.Errorf("Error preparing table: %v", err)
		return
	}
	defer rows.Close()
	defer teardownFn()

	i := 0
	var recv *asetypes.Decimal
	for rows.Next() {
		if err := rows.Scan(&recv); err != nil {
			t.Errorf("Scan failed on %dth scan: %v", i, err)
			continue
		}

		if compareDecimal(recv, mySamples[i]) {

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
