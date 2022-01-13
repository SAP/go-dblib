// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package dsn

import (
	"errors"
	"testing"
)

func TestFormatURI(t *testing.T) {
	type infoExpanded struct {
		Scheme       string `json:"scheme"`
		Key          string `json:"userstorekey"`
		User         string `json:"user"`
		Pass         string `json:"pass"`
		Host         string `json:"host"`
		Port         int    `json:"port"`
		Database     string `json:"database"`
		ConnectProp1 bool   `json:"connectProp1"`
		ConnectProp2 string `json:"connectProp2"`
	}
	type infoError struct {
		Scheme byte `json:"scheme"`
	}

	cases := map[string]struct {
		info         interface{}
		expectString string
		expectErr    error
	}{
		"infoSimple": {
			info: &Info{
				Username: "user",
				Password: "pass",
				Host:     "host",
				Port:     "12345",
				Database: "db",
			},
			expectString: "//user:pass@host:12345/?database=db",
			expectErr:    nil,
		},
		"infoExpandedSimple1": {
			info: &infoExpanded{
				Scheme: "db",
				User:   "user",
				Pass:   "pass",
				Host:   "host",
				Port:   12345,
			},
			expectString: "db://user:pass@host:12345/?connectProp1=false",
			expectErr:    nil,
		},
		"infoExpandedSimple2": {
			info: &infoExpanded{
				Scheme:       "db",
				Host:         "host",
				Port:         12345,
				User:         "user",
				Pass:         "pass",
				ConnectProp1: true,
				ConnectProp2: "connectionProperty",
			},
			expectString: "db://user:pass@host:12345/?connectProp1=true&connectProp2=connectionProperty",
			expectErr:    nil,
		},
		"infoExpandedSimple3": {
			info: &infoExpanded{
				Scheme:       "db",
				Host:         "host",
				Port:         12345,
				User:         "user",
				Pass:         "pass",
				Database:     "db1",
				ConnectProp1: true,
				ConnectProp2: "connectionProperty",
			},
			expectString: "db://user:pass@host:12345/?connectProp1=true&connectProp2=connectionProperty&database=db1",
			expectErr:    nil,
		},
		"infoExpandedKey1": {
			info: &infoExpanded{
				Scheme: "db",
				Key:    "userkey",
			},
			expectString: "db://?KEY=userkey&connectProp1=false",
			expectErr:    nil,
		},
		"infoExpandedKey2": {
			info: &infoExpanded{
				Scheme:       "db",
				Key:          "userkey",
				ConnectProp1: true,
				ConnectProp2: "connectionProperty",
			},
			expectString: "db://?KEY=userkey&connectProp1=true&connectProp2=connectionProperty",
			expectErr:    nil,
		},
		"infoError": {
			info: &infoError{
				Scheme: uint8(1),
			},
			expectString: "",
			expectErr:    errors.New("FormatURI: failed to transform the value of <scheme=<uint8 Value>> to string"),
		},
	}

	for title, cas := range cases {
		t.Run(title, func(t *testing.T) {
			connStr, err := FormatURI(cas.info)
			if err != nil {
				if cas.expectErr == nil {
					t.Errorf("Expected no error, but got error: %s", err)
				}
			}
			if connStr != cas.expectString {
				t.Errorf("Expected <%s> but got <%s>", cas.expectString, connStr)
			}
		})
	}
}
