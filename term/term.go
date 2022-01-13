// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package term

import (
	"flag"
	"fmt"
	"io/ioutil"
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
		bs, err := ioutil.ReadFile(*fInputFile)
		if err != nil {
			return fmt.Errorf("term: error reading file '%s': %w", *fInputFile, err)
		}
		query = string(bs)
	}

	return ParseAndExecQueries(db, query)
}
