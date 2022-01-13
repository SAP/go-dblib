// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package dsn

import (
	"reflect"
	"strings"
)

type TagType string

const (
	OnlyJSON TagType = "json"
	Multiref TagType = "multiref"
	Doc      TagType = "doc"
)

// TagToField returns a mapping from json metadata tags to
// reflect.Values.
//
// If TagType OnlyJSON is passed only the json tag will be mapped.
// If the TagType MultiRef is passed the tags from the `multiref`
// metadata tag will also be mapped to their field.Value.
//
// Example:
//   type Example struct {
//   	Host `json:"host" multiref:"hostname"`
//   }
//
//   example := Example{Host: "host.name.tld"}
//   TagToField(example, OnlyJSON)
//   -> map[string]reflect.Value{
//       "host": example.Host,
//   }
//   TagToField(example, Multiref)
//   -> map[string]reflect.Value{
//       "host": example.Host,
//       "hostname": example.Host,
//   }
//
// If the TagType Doc is passed the json metadata tag value will be
// mapped to a reflect.Value of Kind string containing the value of the
// doc metadata tag..
//
// Example:
//    type Example struct {
//        Host `json:"hostname" doc:"Hostname to connect to"`
//    }
//
//    example := new(Example)
//    TagToField(example, Doc)
//    -> map[string]reflect.Value{
//          "hostname": reflect.ValueOf("Hostname to connect to"),
//       }
//
// If the input is a pointer the reflect.Values will be addressable and
// settable - allowing to modify the fields of the passed structure.
//
// Example:
//   type Example struct {
//       Host `json:"host"`
//   }
//
//   example := Example{Host: "old.host.name"}
//   ttf := TagToField(&example, OnlyJSON)
//   ttf["host"].SetString("new.host.name")
//   fmt.Println(example.Host)
//   -> "new.host.name"
func TagToField(input interface{}, tagType TagType) map[string]reflect.Value {
	return tagToField(reflect.ValueOf(input), tagType)
}

func tagToField(input reflect.Value, tagType TagType) map[string]reflect.Value {
	ttf := map[string]reflect.Value{}

	if input.Kind() == reflect.Ptr || input.Kind() == reflect.Interface {
		// Pointers and Interfaces need to be dereferenced to access
		// their fields.
		input = input.Elem()
	}

	inputT := input.Type()

	for i := 0; i < input.NumField(); i++ {
		field := input.Field(i)

		if field.Kind() == reflect.Struct {
			// Field is embedded, retrieve all of its members and merge them
			// into tTF.
			// This can result in overridden fields, depending on the order
			// of embedded structs.
			for key, value := range tagToField(field, tagType) {
				ttf[key] = value
			}
			continue
		}

		// Grab json tag
		names := strings.Split(inputT.Field(i).Tag.Get(string(OnlyJSON)), ",")
		if names[0] == "" {
			// JSON tag is not set, skip field
			continue
		}

		switch tagType {
		case OnlyJSON:
			ttf[names[0]] = field
		case Multiref:
			names = []string{names[0]}

			// Grab multiref tags
			multirefs := strings.Split(inputT.Field(i).Tag.Get("multiref"), ",")
			names = append(names, multirefs...)

			for _, name := range names {
				ttf[name] = field
			}
		case Doc:
			ttf[names[0]] = reflect.ValueOf(inputT.Field(i).Tag.Get(string(tagType)))
		}
	}

	return ttf
}
