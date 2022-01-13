// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package dsn

import (
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
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

// FormatURI formats the passed values as a simple URI.
//
// Only members with a json metadata tag are added with the json
// metadata value used as key. This key should reflect an URI-resource
// name, e.g. scheme, user, password, host, port. The associated values
// of these keys are passed to an URL-struct.
//
// Additionally, the keys referring to a userstorekey or database will
// be added as connection property using their respective key, e.g.
// "KEY=<userstorekey>" or "database=<database>".
//
// Any key that does not match such URI-resource name will be added as
// connection property following the scheme: "<key>=<value>".
// Thus, non-string struct-members that do not match any URI-resource
// name and default to "0" (int) or "false" (boolean) are added as
// connection property with their default as value.
// To avoid possible issues, use string-members with an empty string.
//
// Example:
//   type Example struct {
//       Scheme       string `json:"scheme"`
//       User         string `json:"user"`
//       Pass         string `json:"pass"`
//       Host         string `json:"host"`
//       Port         int    `json:"port"`
//       Database     string `json:"database"`
//       ConnectProp1 bool   `json:"connectProp1"`
//       ConnectProp2 string `json:"connectProp2"`
//   }
//
//   ex := new(Example)
//   ex.Scheme       = "db"
//   ex.User         = "username"
//   ex.Pass         = "password"
//   ex.Host         = "hostname"
//   ex.Port         = 12345
//   ex.Database     = "db1"
//   ex.ConnectProp1 = false
//   ex.ConnectProp2 = "connectionProperty"
//
//   uri, err := dsn.FormatURI(ex)
//   if err != nil {
//     fmt.Errorf("Something failed: %w", err)
//   }
//
//   fmt.Println(uri)
//
// Will print:
//   db://username:password@hostname:12345/?database=db1&connectProp1=false&connectProp2=connectionProperty
//
// For more examples see format_test.go
func FormatURI(input interface{}) (string, error) {
	// Initialize URL-struct to store connection-values
	urlValues := new(url.URL)
	// Initialize user, password, host, and port variables since they
	// cannot be stored directly in the URL-struct due to varying
	// key-orders
	var user, passwd, host, port string
	// Initialize url.Values to store additional connection-properties
	connectProp := url.Values{}

	// Loop over passed input and store values to respective connection
	// value
	for key, field := range TagToField(input, OnlyJSON) {

		// Store all values as string
		var v string
		switch field.Kind() {
		case reflect.String:
			v = field.String()
		case reflect.Int:
			v = strconv.Itoa(int(field.Int()))
		case reflect.Bool:
			v = strconv.FormatBool(field.Bool())
		default:
			return "", fmt.Errorf("dsn: failed to transform the value of <%s=%v> to string", key, field.String())
		}

		// continue if value is not set
		// e.g. if the passed input contains a member/key without value
		if v == "" {
			continue
		}

		// Fill the URL-struct accordingly to dsn-keywords
		switch key {
		case "scheme":
			urlValues.Scheme = v
		case "user", "username":
			user = v
		case "password", "pass", "passwd":
			passwd = v
		case "host", "hostname":
			host = v
		case "port":
			port = v
		case "userstorekey", "key":
			connectProp.Add("KEY", v)
		case "database", "db":
			connectProp.Add("database", v)

		// If no URI keyword is matched, the key=value must be
		// a connection property
		default:
			connectProp.Add(key, v)
		}

	}

	// If a userstorekey was given, encode the connection properties
	// (that contain the key) and return the resulting string
	if strings.Contains(connectProp.Encode(), "KEY") {
		return fmt.Sprintf("%s://?%s", urlValues.Scheme, connectProp.Encode()), nil
	}

	// Since the key-order can vary the Userinfo and Host/Port-members must be
	// set after the passed input was processed.
	urlValues.User = url.UserPassword(user, passwd)
	urlValues.Host = fmt.Sprintf("%s:%s", host, port)

	// Return URI-string by using url.URL.String() and encoding the
	// connection properties.
	return fmt.Sprintf("%s/?%s", urlValues.String(), connectProp.Encode()), nil
}
