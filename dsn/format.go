// SPDX-FileCopyrightText: 2021 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package dsn

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// FormatSimple formats the passed value as a simple DSN.
//
// Only members with a json metadata tag are added with the json
// metadata value used as the key.
//
// Example:
//   type Example struct {
//       StringA string `json:"a"`
//       IntB int `json:"b"`
//   }
//
//   ex := new(Example)
//   ex.StringA = "a string"
//   ex.IntB = 5
//
//   fmt.Println(dsn.FormatSimple(ex))
//
// Will print:
//   a="a string" b=5
func FormatSimple(input interface{}) string {
	ret := []string{}

	for key, field := range TagToField(input, OnlyJSON) {
		var v interface{}
		switch field.Kind() {
		case reflect.String:
			v = fmt.Sprintf("%q", field.String())
		default:
			v = field
		}

		ret = append(ret, fmt.Sprintf("%s=%v", key, v))
	}

	// Sort for deterministic output
	sort.Strings(ret)

	return strings.Join(ret, " ")
}
