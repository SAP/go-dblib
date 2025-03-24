// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package main

import "log"

func ExampleDoMain() {
	if err := DoMain(); err != nil {
		log.Fatal(err)
	}
	// Output:
	//
	// Host: "a.host.name"
	// Port: "5555"
	// Username: ""
	// Password: ""
	// LogLevel: 5
}
