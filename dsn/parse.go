// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package dsn

import (
	"fmt"
	"net/url"
	"strings"
)

// Parse uses ParseURI or ParseSimple and returns the respective error.
// The decision is made based on the existence of "://" in the passed
// string.
func Parse(dsn string, target interface{}) error {
	if strings.Contains(dsn, "://") {
		return ParseURI(dsn, target)
	}

	return ParseSimple(dsn, target)
}

// ParseURI uses url.Parse to parse the passed string.
//
// Only members with json metadata tags will be filled. Additionally
// multiref metadata tags are recognized.
//
// Hostname, Port, Username and Password are hardcoded to be set to
// "hostname", "port", "username" and "password" respectively due to
// technical limitations.
//
// If the "database" tag is set its member is set to the path of the
// URI, sans the leading "/".
func ParseURI(dsn string, target interface{}) error {
	url, err := url.Parse(dsn)
	if err != nil {
		return fmt.Errorf("dsn: error parsing DSN using url.Parse: %w", err)
	}

	ttf := TagToField(target, Multiref)

	ttf["hostname"].SetString(url.Hostname())
	ttf["port"].SetString(url.Port())

	if url.User != nil {
		ttf["username"].SetString(url.User.Username())
		pass, _ := url.User.Password()
		ttf["password"].SetString(pass)
	}

	if database, ok := ttf["database"]; ok {
		database.SetString(strings.TrimPrefix(url.Path, "/"))
	}

	props := url.Query()
	for key, values := range props {
		field, ok := ttf[key]
		if !ok {
			return fmt.Errorf("dsn: query value %q has no matching field", key)
		}

		if err := setValue(field, values[len(values)-1]); err != nil {
			return fmt.Errorf("dsn: error setting field %s of kind %s to %q",
				key, field.Kind(), values[len(values)-1])
		}
	}

	return nil
}

// ParseSimple parses a simple DSN in the form of "key=value k2=v2".
//
// Only members with a json or multiref metadata tag will be filled.
// Should be member have both json and multiref tags which occur
// multiple times in the passed string the last occurrence will take
// priority.
//
// ParseSimple supports whitespaces in values, but not in keys - given
// that values are quoted with either double or single quotes.
//
// Example:
//   type Example struct {
//       // Recognized as "hostname", "host" and "remote"
//       Host string `json:"hostname" multiref:"host,remote"`
//       // Only recognized as "port"
//       Port string `json:"port"`
//       // Recognized as "database" and "db"
//       Database string `json:"database" multiref:"db"`
//       // Not recognized due to missing metadata
//       Username string
//   }
//
//   ex := new(Example)
//   simple := `host="a.b.c.d" remote='w.x.y.z' port=ssl username="user"`
//   if err := dsn.ParseSimple(simple, ex); err != nil {
//       return err
//   }
//
// Will result in :
//   ex.Host being set to "w.x.y.z" as it's multiref tag 'remote' came last.
//   ex.Port being set to "ssl".
//   ex.Database not being set as no values were provided.
//   ex.Username not being set as it has no metadata.
func ParseSimple(dsn string, target interface{}) error {
	ttf := TagToField(target, Multiref)

	// Valid quotation marks to detect values with whitespaces
	quotations := []byte{'\'', '"'}

	// Split the DSN on whitespace - any quoted values containing
	// whitespaces will be concatenated in the first step in the loop.
	dsnS := strings.Split(dsn, " ")

	for len(dsnS) > 0 {
		var part string
		part, dsnS = dsnS[0], dsnS[1:]

		// If the value starts with a quotation mark consume more parts
		// until the quotation is finished.
		for _, quot := range quotations {
			if !strings.Contains(part, "="+string(quot)) {
				continue
			}

			for part[len(part)-1] != quot {
				part = strings.Join([]string{part, dsnS[0]}, " ")
				dsnS = dsnS[1:]
			}
			break
		}

		partS := strings.SplitN(part, "=", 2)
		if len(partS) != 2 {
			return fmt.Errorf("dsn: recognized DSN part does not contain key/value parts: %q", partS)
		}

		key, value := partS[0], partS[1]

		// Remove quotation from value
		if value != "" {
			for _, quot := range quotations {
				if value[0] == quot && value[len(value)-1] == quot {
					value = value[1 : len(value)-1]
				}
			}
		}

		field, ok := ttf[key]
		if !ok {
			return fmt.Errorf("no field for key %q", key)
		}

		if err := setValue(field, value); err != nil {
			return fmt.Errorf("dsn: error setting field %s of kind %s to %q",
				key, field.Kind(), value)
		}
	}

	return nil
}
