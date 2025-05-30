// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package asetime

import (
	"math"
	"time"
)

// ASEDuration is the reference time-unit as microseconds.
type ASEDuration int

// Base time-units in reference to ASEDuration.
const (
	Microsecond ASEDuration = 1
	Millisecond             = 1000 * Microsecond
	Second                  = 1000 * Millisecond
	Minute                  = 60 * Second
	Hour                    = 60 * Minute
	Day                     = 24 * Hour
)

// Days returns duration in days.
func (d ASEDuration) Days() int { return int(d / Day) }

// Hours returns duration in hours.
func (d ASEDuration) Hours() int { return int(d / Hour) }

// Minutes returns duration in minutes.
func (d ASEDuration) Minutes() int { return int(d / Minute) }

// Seconds returns duration in seconds.
func (d ASEDuration) Seconds() int { return int(d / Second) }

// Milliseconds returns duration in milliseconds.
func (d ASEDuration) Milliseconds() int { return int(d / Millisecond) }

// Microseconds returns duration in microseconds.
func (d ASEDuration) Microseconds() int { return int(d) }

// DurationFromDateTime returns duration based on passed time.Time.
func DurationFromDateTime(t time.Time) ASEDuration {
	y := int64(t.Year())
	m := int64(t.Month())
	d := int64(t.Day())
	// Calculate JDN, formula from Calendars by Doggett
	jd := (1461*(y+4800+(int64(m)-14)/12))/4 + (367*(m-2-12*((m-14)/12)))/12 - (3*((y+4900+(m-14)/12)/100))/4 + d - 32075

	// Calculate Rata Die from JDN and convert to microseconds
	rataDie := int64(jd-1721425) * (int64(time.Duration(24)*time.Hour) / 1000)
	// While Sybase uses Rata Die it still seems to count year 0
	rataDie += int64((time.Hour * 24) * 365 / 1000)

	return ASEDuration(rataDie) + DurationFromTime(t)
}

// DurationFromTime returns duration based on passed time.Time.
func DurationFromTime(t time.Time) ASEDuration {
	hours := int64(time.Duration(t.Hour())*time.Hour) / 1000
	minutes := int64(time.Duration(t.Minute())*time.Minute) / 1000
	seconds := int64(time.Duration(t.Second())*time.Second) / 1000
	nanoseconds := int64(t.Nanosecond()) / 1000

	return ASEDuration(hours + minutes + seconds + nanoseconds)
}

// DurationAsASEDuration returns type time.Duration as type ASEDuration.
func DurationAsASEDuration(d time.Duration) ASEDuration {
	return ASEDuration(d / 1000)
}

// FractionalSecondToMillisecond returns a fractional second as
// millisecond.
func FractionalSecondToMillisecond(s int) ASEDuration {
	return ASEDuration(float64(s)*1000/300) * Millisecond
}

// MillisecondToFractionalSecond returns milliseconds as fractional
// second.
func MillisecondToFractionalSecond(s int) int {
	return int(math.Round(float64(s) * 300 / 1000 / float64(Millisecond)))
}
