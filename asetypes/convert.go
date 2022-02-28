// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package asetypes

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
)

// ValueConverter implements the driver.types.ValueConverter interface.
type ValueConverter struct{}

// DefaultValueConverter implements the driver.types.ValueConverter
// interface.
var DefaultValueConverter ValueConverter

// ConvertValue implements the driver.types.ValueConverter interface.
func (conv ValueConverter) ConvertValue(v interface{}) (driver.Value, error) {
	// Check the default value converter
	if driver.IsValue(v) {
		return v, nil
	}

	// Check for driver.Valuer interface, e.g. database/sql.Null<types>
	if vv, ok := v.(driver.Valuer); ok {
		// Get value
		val, err := vv.Value()
		if err != nil {
			return nil, err
		}
		// Return nil, if value is nil
		if val == nil {
			return nil, nil
		}
		// Catch return of sql.NullInt32, since it's returned as int64
		switch vv.(type) {
		case sql.NullInt32:
			v = int32(val.(int64))
		default:
			v = val
		}
	}

	// Convert any values that can be handled as another type
	switch value := v.(type) {
	case int:
		return int64(value), nil
	case uint:
		return uint64(value), nil
	}

	// Check the reflect types if the value is handled.
	sv := reflect.TypeOf(v)
	for _, kind := range ReflectTypes {
		if kind == sv {
			return v, nil
		}
	}

	return nil, fmt.Errorf("unsupported type %T, a %s", v, sv.Kind())
}
