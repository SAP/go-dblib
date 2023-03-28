// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
// SPDX-FileCopyrightText: 2023 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"database/sql"
	"testing"
)

// DoTestSQLQueryRow runs tests for sql.QueryRow.
func DoTestSQLQueryRow(t *testing.T) {
	TestForEachDB("TestSQLQueryRowSingleRow", t, testSQLQueryRowSingleRow)
	TestForEachDB("TestSQLQueryRowMultipleRows", t, testSQLQueryRowMultipleRows)
}

func testSQLQueryRowSingleRow(t *testing.T, db *sql.DB, tableName string) {
	if _, err := db.Exec("create table " + tableName + " (a bigint, b varchar(30))"); err != nil {
		t.Errorf("error creating table %s: %v", tableName, err)
		return
	}

	if _, err := db.Exec("insert into "+tableName+" values (?, ?)", 1, "one"); err != nil {
		t.Errorf("error inserting values into table %s: %v", tableName, err)
		return
	}

	row := db.QueryRow("select * from " + tableName)
	if row == nil {
		t.Errorf("db.QueryRow returned nil row")
		return
	}

	var a int
	var b string
	if err := row.Scan(&a, &b); err != nil {
		t.Errorf("error scanning row: %v", err)
		return
	}

	if a != 1 {
		t.Errorf("scanned integer is %q, expected %q", a, 1)
	}

	if b != "one" {
		t.Errorf("scanned string is %q, expected %q", b, "one")
	}
}

func testSQLQueryRowMultipleRows(t *testing.T, db *sql.DB, tableName string) {
	if _, err := db.Exec("create table " + tableName + " (a bigint, b varchar(30))"); err != nil {
		t.Errorf("error creating table %s: %v", tableName, err)
		return
	}

	if _, err := db.Exec("insert into "+tableName+" values (?, ?)", 1, "one"); err != nil {
		t.Errorf("error inserting values into table %s: %v", tableName, err)
		return
	}

	if _, err := db.Exec("insert into "+tableName+" values (?, ?)", 2, "two"); err != nil {
		t.Errorf("error inserting values into table %s: %v", tableName, err)
		return
	}

	row := db.QueryRow("select * from " + tableName)
	if row == nil {
		t.Errorf("db.QueryRow returned nil row")
		return
	}

	var a int
	var b string
	if err := row.Scan(&a, &b); err != nil {
		t.Errorf("error scanning row: %v", err)
		return
	}

	if a != 1 {
		t.Errorf("scanned integer is %q, expected %q", a, 1)
	}

	if b != "one" {
		t.Errorf("scanned string is %q, expected %q", b, "one")
	}
}
