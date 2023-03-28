// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
// SPDX-FileCopyrightText: 2023 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"database/sql"
	"fmt"
	"testing"
)

// DoTestSQLTx runs tests for sql.Tx.
func DoTestSQLTx(t *testing.T) {
	t.Run("Commit",
		func(t *testing.T) {
			TestForEachDB("TestSQLTxCommit", t, testSQLTxCommit)
		},
	)

	t.Run("Rollback",
		func(t *testing.T) {
			TestForEachDB("TestSQLTxRollback", t, testSQLTxRollback)
		},
	)
}

func testSQLTxCommit(t *testing.T, db *sql.DB, tableName string) {
	if _, err := db.Exec(fmt.Sprintf("create table %s (a bigint)", tableName)); err != nil {
		t.Errorf("Error creating table %s: %v", tableName, err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		t.Errorf("Failed to initialize transaction: %v", err)
		return
	}

	sample := 5
	if _, err := tx.Exec(fmt.Sprintf("insert into %s (a) values (?)", tableName), sample); err != nil {
		t.Errorf("Error inserting value %d in transaction: %v", sample, err)
		return
	}

	if err := tx.Commit(); err != nil {
		t.Errorf("Error committing transaction: %v", err)
		return
	}

	rows, err := db.Query(fmt.Sprintf("select * from %s", tableName))
	if err != nil {
		t.Errorf("Error selecting from table created in transaction: %v", err)
		return
	}
	defer rows.Close()

	var recv int
	for rows.Next() {
		if err := rows.Scan(&recv); err != nil {
			t.Errorf("Scan failed: %v", err)
			return
		}

		if recv != sample {
			t.Errorf("Scanned value does not match inserted value")
			t.Errorf("Expected: %d", sample)
			t.Errorf("Received: %d", recv)
		}
	}

	if err := rows.Err(); err != nil {
		t.Errorf("Error preparing rows: %v", err)
	}
}

func testSQLTxRollback(t *testing.T, db *sql.DB, tableName string) {
	if _, err := db.Exec(fmt.Sprintf("create table %s (a bigint)", tableName)); err != nil {
		t.Errorf("Error creating table %s: %v", tableName, err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		t.Errorf("Failed to initialize transaction: %v", err)
		return
	}

	sample := 5
	if _, err := tx.Exec(fmt.Sprintf("insert into %s (a) values (?)", tableName), sample); err != nil {
		t.Errorf("Error inserting value %d in transaction: %v", sample, err)
		return
	}

	if err := tx.Rollback(); err != nil {
		t.Errorf("Error while rolling back transaction: %v", err)
		return
	}

	rows, err := db.Query(fmt.Sprintf("select * from %s", tableName))
	if err != nil {
		t.Errorf("Error selecting from table created in transaction: %v", err)
		return
	}
	defer rows.Close()

	if rows.Next() {
		var recv int
		err := rows.Scan(&recv)
		t.Errorf("Insert was rolled back, still received value '%v', scan error: %v", recv, err)
	}

	if err := rows.Err(); err != nil {
		t.Errorf("Error preparing rows: %v", err)
	}
}
