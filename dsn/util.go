// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package dsn

import (
	"fmt"
	"reflect"
	"strconv"
)

func setValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("error parsing %q as bool: %w", value, err)
		}
		field.SetBool(b)
	case reflect.Int:
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing %q as int: %w", value, err)
		}
		field.SetInt(n)
	default:
		return fmt.Errorf("unhandled field kind: %s", field.Kind())
	}

	return nil
}
