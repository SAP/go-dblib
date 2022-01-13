// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
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
