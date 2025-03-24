// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package term

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"database/sql"
)

var (
	fInputFile = flag.String("f", "", "Read SQL commands from file")
)

// Entrypoint controls the execution of the program by starting the
// interactive command-line or executing the passed query or input-file.
func Entrypoint(db *sql.DB, args []string) error {
	if len(args) == 0 && *fInputFile == "" {
		return Repl(db)
	}

	query := strings.Join(args, " ") + ";"

	if *fInputFile != "" {
		bs, err := os.ReadFile(*fInputFile)
		if err != nil {
			return fmt.Errorf("term: error reading file '%s': %w", *fInputFile, err)
		}
		query = string(bs)
	}

	return ParseAndExecQueries(db, query)
}
