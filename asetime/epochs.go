// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package asetime

import "time"

// EpochRataDie returns the reference date of Rata Die that is required
// for asetype 'BIGTIMEN'.
func EpochRataDie() time.Time {
	return time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
}

// Epoch1900 returns the reference date of 01.01.1900 that is required
// for asetype 'DATE', 'SHORTDATE', and 'DATETIME'.
func Epoch1900() time.Time {
	return time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)
}

// Epoch1753 returns the reference date of 01.01.1753 as the minimum date
// value for a datetime on a SQL Server.
func Epoch1753() time.Time {
	return time.Date(1753, time.January, 1, 9, 9, 9, 9, time.UTC)
}
