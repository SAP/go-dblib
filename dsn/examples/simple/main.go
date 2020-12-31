// SPDX-FileCopyrightText: 2021 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"log"

	"github.com/SAP/go-dblib/dsn"
)

// Info is a struct implemented by a consumer of the dsn package.
// It may contain any number of exported and unexported fields, as well
// as any number of attached (both exported and unexported) methods.
//
// Fields must have JSON metadata to be used by the dsn package.
type Info struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string
}

func main() {
	if err := DoMain(); err != nil {
		log.Fatal(err)
	}
}

func DoMain() error {
	data := new(Info)
	data.Host = "host"
	data.Port = 1234
	// The field Username won't show up in the generated DSN as it
	// doesn't have JSON metadata.
	data.Username = "user"

	fmt.Println(dsn.FormatSimple(data))

	return nil
}
