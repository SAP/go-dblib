// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package term

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/chzyer/readline"
)

var (
	rl *readline.Instance
	// PromptDatabaseName contains the used database name when using
	// the prompt.
	PromptDatabaseName string
	promptMultiline    bool
)

// UpdatePrompt updates the displayed prompt in interactive use.
func UpdatePrompt() {
	prompt := "> "

	if promptMultiline {
		prompt = ">>> "
	}

	if PromptDatabaseName != "" {
		prompt = PromptDatabaseName + prompt
	}

	if rl != nil {
		rl.SetPrompt(prompt)
	}
}

// Repl is the interactive interface that reads, evaluates, and prints
// the passed queries.
func Repl(db *sql.DB) error {
	var err error
	rl, err = readline.New("")
	if err != nil {
		return fmt.Errorf("term: failed to initialize readline: %w", err)
	}
	defer rl.Close()

	cmds := []string{}
	for {
		UpdatePrompt()

		line, readlineErr := rl.Readline()
		line = strings.TrimSpace(line)

		// exit immediately on non-EOF errors
		if readlineErr != nil && !errors.Is(readlineErr, io.EOF) {
			return fmt.Errorf("term: received error from readline: %w", err)
		}

		// Only add non empty lines
		if len(line) > 0 {
			cmds = append(cmds, line)
		}

		// exit if no statements are present and EOF is given
		if errors.Is(readlineErr, io.EOF) && len(cmds) == 0 {
			return nil
		}

		// Start multiline if query is not finished
		if !strings.HasSuffix(line, ";") && len(line) > 0 {
			promptMultiline = true
			continue
		}

		// Command is finished, reset prompt and execute
		promptMultiline = false

		line = strings.Join(cmds, " ")
		cmds = []string{}

		if err := ParseAndExecQueries(db, line); err != nil {
			log.Println(err)
		}

		if errors.Is(readlineErr, io.EOF) {
			return nil
		}
	}
}
