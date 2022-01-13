// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"log"

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
	data := new(Info)

	data.Host = "host"
	data.Port = "1234"
	data.Username = "user"
	data.Password = "pass"
	data.S = "extra string"
	data.I = 5
	data.B = true

	fmt.Println(dsn.FormatSimple(data))

	return nil
}
