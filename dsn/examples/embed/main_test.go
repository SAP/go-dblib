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
	// b=true database="" host="host" i=5 password="pass" port="1234" s="extra string" username="user"
}
