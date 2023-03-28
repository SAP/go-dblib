// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
// SPDX-FileCopyrightText: 2023 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package dsn

import (
	"flag"
	"fmt"
	"reflect"
)

// FlagSet creates a flag.FlagSet to be used with the stdlib flag
// package or flag-compatible third party packages.
//
// Flags will be created for members based on their json metadata. If no
// json metadata is present no flag will be created.
// If the member has a doc metadata tag its value will be used as the
// usage argument for the flag.
//
// The resulting FlagSet can be used with e.g. github.com/spf13/pflag to
// merge multiple FlagSets.
func FlagSet(name string, errorHandling flag.ErrorHandling, info interface{}) (*flag.FlagSet, error) {
	ttf := TagToField(info, OnlyJSON)
	docs := TagToField(info, Doc)
	flagset := flag.NewFlagSet(name, errorHandling)

	for key, field := range ttf {
		usage := ""
		docField, ok := docs[key]
		if ok {
			usage = docField.String()
		}

		switch field.Kind() {
		case reflect.String:
			flagset.StringVar(field.Addr().Interface().(*string),
				key, field.String(), usage)
		case reflect.Bool:
			flagset.BoolVar(field.Addr().Interface().(*bool),
				key, field.Interface().(bool), usage)
		case reflect.Int:
			flagset.IntVar(field.Addr().Interface().(*int),
				key, field.Interface().(int), usage)
		default:
			return nil, fmt.Errorf("dsn: unhandled reflect kind %q", field.Kind())
		}
	}

	return flagset, nil
}
