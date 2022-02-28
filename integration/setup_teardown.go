// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/SAP/go-dblib/dsn"
)

var (
	// ASE doesn't handle creating multiple databases concurrently well.
	// To prevent spurious test errors the DBCreateLock is used to
	// synchronise the goroutines creating databases.
	DBCreateLock = new(sync.Mutex)
)

const timeout = 10 * time.Second

// Waits until no CREATE DATABASE or DROP DATABASE processes are listed
// in master..sysprocesses.
func waitForDBCreateDropProcess(ctx context.Context, db *sql.DB) error {
	query := `select count(cmd)
		from master..sysprocesses
		where cmd like "CREATE DATABASE"
			or cmd like "DROP DATABASE"`

	queries := 0
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		row := db.QueryRowContext(timeoutCtx, query)

		var count int
		err := row.Scan(&count)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("error scanning query result: %w", err)
		}

		if !errors.Is(err, sql.ErrNoRows) && count == 0 {
			return nil
		}

		queries++
		time.Sleep(
			time.Duration(
				int(math.Pow(2, float64(queries))) + rand.Intn(10),
			),
		)
	}
}

func waitForDBInSysdatabases(ctx context.Context, db *sql.DB, dbname string) error {
	query := `select count(name)
		from master..sysdatabases
		where name = "` + dbname + `"`

	queries := 0
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		row := db.QueryRowContext(timeoutCtx, query)

		var count int
		err := row.Scan(&count)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("error scanning query result: %w", err)
		}

		if !errors.Is(err, sql.ErrNoRows) && count == 1 {
			return nil
		}

		queries++

		time.Sleep(
			time.Duration(
				int(math.Pow(2, float64(queries))) + rand.Intn(10),
			),
		)
	}
}

// SetupDB safely creates a database with the given name and sets
// .Database on the passed info.
func SetupDB(ctx context.Context, info interface{}, dbname string) error {
	ttf := dsn.TagToField(info, dsn.OnlyJSON)
	field, ok := ttf["database"]
	if !ok {
		return fmt.Errorf("integration: provided info does not have the 'database' field")
	}

	DBCreateLock.Lock()
	defer DBCreateLock.Unlock()

	db, err := sql.Open("ase", dsn.FormatSimple(info))
	if err != nil {
		return fmt.Errorf("integration: failed to open database: %w", err)
	}
	defer db.Close()

	if err := waitForDBCreateDropProcess(ctx, db); err != nil {
		return fmt.Errorf("integration: error waiting for CREATE/DROP DATABASE processes: %w", err)
	}

	conn, err := db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("integration: failed to open connection: %w", err)
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, "use master"); err != nil {
		return fmt.Errorf("integration: failed to switch context to master: %w", err)
	}

	if _, err := conn.ExecContext(ctx, fmt.Sprintf("if db_id('%s') is not null drop database %s", dbname, dbname)); err != nil {
		return fmt.Errorf("integration: error on conditional drop of database: %w", err)
	}

	if _, err := conn.ExecContext(ctx, "create database "+dbname); err != nil {
		return fmt.Errorf("integration: failed to create database: %w", err)
	}

	// Wait for this CREATE process
	if err := waitForDBCreateDropProcess(ctx, db); err != nil {
		return fmt.Errorf("integration: error waiting for CREATE/DROP DATABASE processes: %w", err)
	}

	if err := waitForDBInSysdatabases(ctx, db, dbname); err != nil {
		return fmt.Errorf("integration: error waiting for database %q to be available in master..sysdatabases: %w", dbname, err)
	}

	if _, err := conn.ExecContext(ctx, "use "+dbname); err != nil {
		return fmt.Errorf("integration: failed to switch context to %s: %w", dbname, err)
	}

	field.SetString(dbname)
	return nil
}

// TeardownDB safely deletes the database indicated by .Database of the
// passed info and unsets the member.
func TeardownDB(ctx context.Context, info interface{}) error {
	ttf := dsn.TagToField(info, dsn.OnlyJSON)
	field, ok := ttf["database"]
	if !ok {
		return errors.New("integration: provided info does not have the 'database' field")
	}

	// Field must be unset beforehand to prevent issues with queries
	// erroring because new connections cannot connect to the deleted
	// database.
	dbname := field.String()
	field.SetString("")

	DBCreateLock.Lock()
	defer DBCreateLock.Unlock()

	db, err := sql.Open("ase", dsn.FormatSimple(info))
	if err != nil {
		return fmt.Errorf("integration: failed to open database: %w", err)
	}
	defer db.Close()

	// Wait for current CRETAE/DROP processes
	if err := waitForDBCreateDropProcess(ctx, db); err != nil {
		return fmt.Errorf("integration: error waiting for CREATE/DROP DATABASE processes: %w", err)
	}

	conn, err := db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("integration: error getting db.Conn: %w", err)
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, "use master"); err != nil {
		return fmt.Errorf("integration: error switching to master: %w", err)
	}

	if _, err := conn.ExecContext(ctx, fmt.Sprintf("if db_id('%s') is not null drop database %s", dbname, dbname)); err != nil {
		return fmt.Errorf("integration: error dropping database %q: %w", dbname, err)
	}

	// Wait for this DROP process
	if err := waitForDBCreateDropProcess(ctx, db); err != nil {
		return fmt.Errorf("integration: error waiting for CREATE/DROP DATABASE processes: %w", err)
	}

	return nil
}

// SetupTableInsert creates a table with the passed type and inserts all
// passed samples as rows.
func SetupTableInsert(db *sql.DB, tableName, aseType string, samples ...interface{}) (*sql.Rows, func() error, error) {
	if _, err := db.Exec(fmt.Sprintf("create table %s (a %s)", tableName, aseType)); err != nil {
		return nil, nil, fmt.Errorf("failed to create table: %w", err)
	}
	stmt, err := db.Prepare(fmt.Sprintf("insert into %s (a) values (?)", tableName))
	if err != nil {
		return nil, nil, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	for _, sample := range samples {
		if _, err := stmt.Exec(sample); err != nil {
			return nil, nil, fmt.Errorf("failed to execute prepared statement with %v: %w", sample, err)
		}
	}

	for _, sample := range samples {
		if _, err := db.Exec(fmt.Sprintf("insert into %s (a) values (?)", tableName), sample); err != nil {
			return nil, nil, fmt.Errorf("failed to insert value with %v: %w", sample, err)
		}
	}

	rows, err := db.Query("select * from " + tableName)
	if err != nil {
		return nil, nil, fmt.Errorf("error selecting from %s: %w", tableName, err)
	}

	teardownFn := func() error {
		_, err := db.Exec("drop table " + tableName)
		return err
	}

	return rows, teardownFn, nil
}
