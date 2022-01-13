// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package term

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/SAP/go-dblib/asetypes"
)

var (
	fMaxColPrintLength = flag.Int("max-col-length", 50, "Maximum number of characters to print for column")
	fPrintColType      = flag.Bool("print-col-type", false, "Display the column type next to the column name")
)

// GenericExecer is the interface providing the GenericExec method.
type GenericExecer interface {
	// GenericExec is the central method through which SQL statements
	// are sent to ASE.
	GenericExec(context.Context, string, []driver.NamedValue) (driver.Rows, driver.Result, error)
}

func process(db *sql.DB, query string) error {
	conn, err := db.Conn(context.Background())
	if err != nil {
		return fmt.Errorf("error getting sql.Conn: %w", err)
	}
	defer conn.Close()

	return conn.Raw(func(driverConn interface{}) error {
		return rawProcess(driverConn, query)
	})
}

func rawProcess(driverConn interface{}, query string) error {
	execer, ok := driverConn.(GenericExecer)
	if !ok {
		return fmt.Errorf("invalid driver, must support GenericExecer")
	}

	rows, result, err := execer.GenericExec(context.Background(), query, nil)
	if err != nil {
		return fmt.Errorf("GenericExec failed: %w", err)
	}

	if rows != nil && !reflect.ValueOf(rows).IsNil() {
		defer rows.Close()

		if err := processRows(rows); err != nil {
			return fmt.Errorf("error processing rows: %w", err)
		}
	}

	if result != nil && !reflect.ValueOf(result).IsNil() {
		if err := processResult(result); err != nil {
			return fmt.Errorf("error processing result: %w", err)
		}
	}

	return nil
}

type rowsColumnTypeDisplayLengther interface {
	ColumnTypeDisplayLength(int) (int64, bool)
}

func processRows(rows driver.Rows) error {
	rowsColumnTypeDisplayLength, _ := rows.(rowsColumnTypeDisplayLengther)

	rowsColumnTypeLength, ok := rows.(driver.RowsColumnTypeLength)
	if !ok {
		return errors.New("rows does not support driver.RowsColumnTypLength")
	}

	rowsColumnTypeName, ok := rows.(driver.RowsColumnTypeDatabaseTypeName)
	if !ok {
		return errors.New("rows does not support driver.RowsColumnTypesDatabaseTypeName")
	}

	colNames := rows.Columns()
	// Check if rows is empty
	if len(colNames) == 0 {
		return nil
	}

	colLengths := map[int]int{}

	fmt.Printf("|")
	for i, colName := range colNames {
		cellLen := len(colName)
		typeName := rowsColumnTypeName.ColumnTypeDatabaseTypeName(i)
		if *fPrintColType {
			cellLen += 1 + len(typeName)
		}

		colTypeLen, ok := rowsColumnTypeLength.ColumnTypeLength(i)
		if ok && int(colTypeLen) > cellLen {
			cellLen = int(colTypeLen)
		}

		if rowsColumnTypeDisplayLength != nil {
			colTypeLen, ok := rowsColumnTypeDisplayLength.ColumnTypeDisplayLength(i)
			if ok && int(colTypeLen) > cellLen {
				cellLen = int(colTypeLen)
			}
		}

		if cellLen > *fMaxColPrintLength {
			cellLen = *fMaxColPrintLength
		}

		s := " %-" + strconv.Itoa(cellLen) + "s |"

		if *fPrintColType {
			fmt.Printf(s, colName+" "+typeName)
		} else {
			fmt.Printf(s, colName)
		}

		colLengths[i] = cellLen
	}
	fmt.Printf("\n")

	cells := make([]driver.Value, len(colNames))

	for {
		if err := rows.Next(cells); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return fmt.Errorf("scanning cells failed: %w", err)
		}

		fmt.Printf("|")

		for i, cell := range cells {
			var cellS string
			switch rowsColumnTypeName.ColumnTypeDatabaseTypeName(i) {
			case "DECIMAL", "DECN", "NUMN":
				cellS = cell.(*asetypes.Decimal).String()
			case "IMAGE", "BINARY", "LONGBINARY", "VARBINARY":
				cellS = hex.EncodeToString(cell.([]byte))
			default:
				cellS = fmt.Sprintf("%v", (interface{})(cell))
			}

			if len(cellS) > colLengths[i] {
				cellS = cellS[:colLengths[i]-3] + "..."
			}

			fmt.Printf(" %-"+strconv.Itoa(colLengths[i])+"v |", cellS)
		}

		fmt.Printf("\n")
	}

	if nextResultSetter, ok := rows.(driver.RowsNextResultSet); ok && nextResultSetter.HasNextResultSet() {
		return processRows(rows)
	}

	return nil
}

func processResult(result sql.Result) error {
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Retrieving the affected rows failed: %w", err)
	}

	if affectedRows >= 0 {
		fmt.Printf("Rows affected: %d\n", affectedRows)
	}
	return nil
}
