// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package dblib

import "database/sql/driver"

// ValuesToNamedValues translates a slice of driver.Values into
// driver.NamedValues.
func ValuesToNamedValues(values []driver.Value) []driver.NamedValue {
	ret := make([]driver.NamedValue, len(values))

	for i, value := range values {
		ret[i] = driver.NamedValue{
			Name:    "",
			Ordinal: i + 1,
			Value:   value,
		}
	}

	return ret
}
