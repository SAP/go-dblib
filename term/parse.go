// SPDX-FileCopyrightText: 2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package term

import (
	"database/sql"
	"fmt"
	"strings"
)

func ParseAndExecQueries(db *sql.DB, line string) error {
	builder := strings.Builder{}
	currentlyQuoted := false

	for _, chr := range line {
		switch chr {
		case '"', '\'':
			if currentlyQuoted {
				currentlyQuoted = false
				builder.WriteRune(chr)
			} else {
				currentlyQuoted = true
				builder.WriteRune(chr)
			}
		case ';':
			if currentlyQuoted {
				builder.WriteRune(chr)
			} else {
				if err := process(db, builder.String()); err != nil {
					return fmt.Errorf("term: failed to process query: %w", err)
				}
				builder.Reset()
			}
		default:
			builder.WriteRune(chr)
		}
	}

	if builder.String() != "" {
		if err := process(db, builder.String()); err != nil {
			return fmt.Errorf("term: failed to process query: %w", err)
		}
	}

	return nil
}
