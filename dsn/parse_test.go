// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
// SPDX-FileCopyrightText: 2023 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package dsn

import (
	"testing"
)

func TestParseURI(t *testing.T) {
	type Embed struct {
		Info
		S  string `json:"s"`
		S2 string `json:"s2"`
	}

	cases := map[string]struct {
		uri            string
		target, expect interface{}
	}{
		"simple": {
			uri:    "ase://user1:pass1@host1:5599/db1",
			target: new(Info),
			expect: &Info{
				Host:     "host1",
				Port:     "5599",
				Username: "user1",
				Password: "pass1",
				Database: "db1",
			},
		},
		"embed": {
			uri:    "ase://user1:pass1@host1:5599/db1?s=aString",
			target: new(Embed),
			expect: &Embed{
				Info: Info{
					Host:     "host1",
					Port:     "5599",
					Username: "user1",
					Password: "pass1",
					Database: "db1",
				},
				S: "aString",
			},
		},
		"embed space": {
			uri:    "ase://user1:pass1@host1:5599/db1?s=a string",
			target: new(Embed),
			expect: &Embed{
				Info: Info{
					Host:     "host1",
					Port:     "5599",
					Username: "user1",
					Password: "pass1",
					Database: "db1",
				},
				S: "a string",
			},
		},
		"embed multiple query values": {
			uri:    "ase://user1:pass1@host1:5599/db1?s=a string&s2=another string",
			target: new(Embed),
			expect: &Embed{
				Info: Info{
					Host:     "host1",
					Port:     "5599",
					Username: "user1",
					Password: "pass1",
					Database: "db1",
				},
				S:  "a string",
				S2: "another string",
			},
		},
		"empty value resets to default": {
			uri: "ase://user1@host1:5599/db1?s=a string&s2=another string",
			target: &Embed{
				Info: Info{
					Password: "this should be gone",
				},
			},
			expect: &Embed{
				Info: Info{
					Host:     "host1",
					Port:     "5599",
					Username: "user1",
					Password: "",
					Database: "db1",
				},
				S:  "a string",
				S2: "another string",
			},
		},
		"unset query values do not override prsets": {
			uri: "ase://user1:pass1@host1:5599/db1?s=a string",
			target: &Embed{
				S2: "this should stay",
			},
			expect: &Embed{
				Info: Info{
					Host:     "host1",
					Port:     "5599",
					Username: "user1",
					Password: "pass1",
					Database: "db1",
				},
				S:  "a string",
				S2: "this should stay",
			},
		},
		"empty query values override prsets": {
			uri: "ase://user1:pass1@host1:5599/db1?s=a string&s2=",
			target: &Embed{
				S2: "this should be gone",
			},
			expect: &Embed{
				Info: Info{
					Host:     "host1",
					Port:     "5599",
					Username: "user1",
					Password: "pass1",
					Database: "db1",
				},
				S:  "a string",
				S2: "",
			},
		},
	}

	for title, cas := range cases {
		t.Run(title, func(t *testing.T) {
			if err := ParseURI(cas.uri, cas.target); err != nil {
				t.Errorf("Parsing URI failed: %v", err)
				return
			}

			checker(t, false,
				TagToField(cas.target, OnlyJSON),
				TagToField(cas.expect, OnlyJSON),
			)
		})
	}
}

func TestParseSimple(t *testing.T) {
	type Embed struct {
		Info
		S string `json:"s"`
		I int    `json:"i"`
	}

	cases := map[string]struct {
		dsn            string
		target, expect interface{}
	}{
		"simple": {
			dsn:    "database=db1 host=host1 user=user1 pass=pass1",
			target: new(Info),
			expect: &Info{
				Database: "db1",
				Host:     "host1",
				Username: "user1",
				Password: "pass1",
			},
		},
		"simple double quotes": {
			dsn:    "database=db1 host=host1 user=user1 pass=\"a password with spaces\"",
			target: new(Info),
			expect: &Info{
				Database: "db1",
				Host:     "host1",
				Username: "user1",
				Password: "a password with spaces",
			},
		},
		"simple single quotes": {
			dsn:    "database=db1 host=host1 user=user1 pass='a password with spaces'",
			target: new(Info),
			expect: &Info{
				Database: "db1",
				Host:     "host1",
				Username: "user1",
				Password: "a password with spaces",
			},
		},
		"embed": {
			dsn:    "database=db1 host=host1 user=user1 pass=pass1 s=aString i=5",
			target: new(Embed),
			expect: &Embed{
				Info: Info{
					Database: "db1",
					Host:     "host1",
					Username: "user1",
					Password: "pass1",
				},
				S: "aString",
				I: 5,
			},
		},
		"embed double quotes": {
			dsn:    "database=db1 host=host1 user=user1 pass=\"a password with spaces\" s=\"a string with spaces\" i=5",
			target: new(Embed),
			expect: &Embed{
				Info: Info{
					Database: "db1",
					Host:     "host1",
					Username: "user1",
					Password: "a password with spaces",
				},
				S: "a string with spaces",
				I: 5,
			},
		},
		"embed single quotes": {
			dsn:    "database=db1 host=host1 user=user1 pass='a password with spaces' s='a string with spaces' i=5",
			target: new(Embed),
			expect: &Embed{
				Info: Info{
					Database: "db1",
					Host:     "host1",
					Username: "user1",
					Password: "a password with spaces",
				},
				S: "a string with spaces",
				I: 5,
			},
		},
	}

	for title, cas := range cases {
		t.Run(title, func(t *testing.T) {
			if err := ParseSimple(cas.dsn, cas.target); err != nil {
				t.Errorf("Parsing simple DSN failed: %v", err)
				return
			}

			checker(t, false,
				TagToField(cas.target, OnlyJSON),
				TagToField(cas.expect, OnlyJSON),
			)
		})
	}
}
