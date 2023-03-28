// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
// SPDX-FileCopyrightText: 2023 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/SAP/go-dblib/dsn"
)

type MyDSN struct {
	dsn.Info
	LogLevel int `json:"log-level"`
}

func main() {
	if err := DoMain(); err != nil {
		log.Fatal(err)
	}
}

func DoMain() error {
	// Init info
	info := new(MyDSN)

	// Bind info members to flags
	flagset, err := dsn.FlagSet("", flag.ContinueOnError, info)
	if err != nil {
		return err
	}

	arguments := []string{
		// Two arguments from dsn.Info
		"-host=a.host.name",
		"-port=5555",
		// One argument from MyDSN
		"-log-level=5",
	}

	// Parse arguments
	if err := flagset.Parse(arguments); err != nil {
		return err
	}

	// Host and Port are filled
	fmt.Printf("Host: %q\n", info.Host)
	fmt.Printf("Port: %q\n", info.Port)
	// Username and Password are empty
	fmt.Printf("Username: %q\n", info.Username)
	fmt.Printf("Password: %q\n", info.Password)
	// LogLevel is filled
	fmt.Printf("LogLevel: %d\n", info.LogLevel)

	return nil
}
