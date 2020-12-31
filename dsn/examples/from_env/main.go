// SPDX-FileCopyrightText: 2021 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SAP/go-dblib/dsn"
)

type Info struct {
	dsn.Info
	S string `json:"s"`
	I int    `json:"i"`
	B bool   `json:"b"`
}

func main() {
	if err := DoMain(); err != nil {
		log.Fatal(err)
	}
}

func DoMain() error {
	os.Setenv("TDS_HOST", "host")
	os.Setenv("TDS_PORT", "2222")
	os.Setenv("TDS_USER", "user")
	os.Setenv("TDS_PASS", "pass")
	os.Setenv("TDS_DATABASE", "dbname")
	os.Setenv("TDS_S", "a string")
	os.Setenv("TDS_I", "5")

	data := new(Info)

	if err := dsn.FromEnv("TDS", data); err != nil {
		return err
	}

	fmt.Println(dsn.FormatSimple(data))

	return nil
}
