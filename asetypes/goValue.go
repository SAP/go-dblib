// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
// SPDX-FileCopyrightText: 2023 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package asetypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
	"unicode/utf16"

	"github.com/SAP/go-dblib/asetime"
)

// GoValue returns a value-interface based on a given byte slice and
// depending on the ASE data type.
func (t DataType) GoValue(endian binary.ByteOrder, bs []byte) (interface{}, error) {
	if t.ByteSize() != -1 && len(bs) != t.ByteSize() {
		return nil, fmt.Errorf("byte slice has invalid length of %d, expected %d bytes", len(bs), t.ByteSize())
	}

	val, err := t.goValue(endian, bs)
	if err != nil {
		return nil, fmt.Errorf("error converting %v into value of type %s: %w", bs, t, err)
	}

	return val, nil
}

func (t DataType) goValue(endian binary.ByteOrder, bs []byte) (interface{}, error) {
	buffer := bytes.NewBuffer(bs)

	switch t {
	case INT1:
		var x uint8
		err := binary.Read(buffer, endian, &x)
		return x, err
	case INT2:
		var x int16
		err := binary.Read(buffer, endian, &x)
		return x, err
	case INT4:
		var x int32
		err := binary.Read(buffer, endian, &x)
		return x, err
	case INT8:
		var x int64
		err := binary.Read(buffer, endian, &x)
		return x, err
	case INTN:
		switch len(bs) {
		case 0:
			return nil, nil
		case 1:
			return INT1.GoValue(endian, bs)
		case 2:
			return INT2.GoValue(endian, bs)
		case 4:
			return INT4.GoValue(endian, bs)
		case 8:
			return INT8.GoValue(endian, bs)
		default:
			return nil, fmt.Errorf("invalid length for INTN: %d", len(bs))
		}
	case UINT2:
		var x uint16
		err := binary.Read(buffer, endian, &x)
		return x, err
	case UINT4:
		var x uint32
		err := binary.Read(buffer, endian, &x)
		return x, err
	case UINT8:
		var x uint64
		err := binary.Read(buffer, endian, &x)
		return x, err
	case UINTN:
		switch len(bs) {
		case 0:
			return nil, nil
		case 1:
			return INT1.GoValue(endian, bs)
		case 2:
			return UINT2.GoValue(endian, bs)
		case 4:
			return UINT4.GoValue(endian, bs)
		case 8:
			return UINT8.GoValue(endian, bs)
		default:
			return nil, fmt.Errorf("invalid length for UINTN: %d", len(bs))
		}
	case FLT4:
		var x float32
		err := binary.Read(buffer, endian, &x)
		return x, err
	case FLT8:
		var x float64
		err := binary.Read(buffer, endian, &x)
		return x, err
	case FLTN:
		switch len(bs) {
		case 0:
			return nil, nil
		case 4:
			return FLT4.GoValue(endian, bs)
		case 8:
			return FLT8.GoValue(endian, bs)
		default:
			return nil, fmt.Errorf("invalid length for FLTN: %d", len(bs))
		}
	case BIT:
		if bs[0] == 0x1 {
			return true, nil
		}
		return false, nil
	case LONGBINARY, BINARY, VARBINARY, IMAGE:
		if len(bs) == 0 {
			return nil, nil
		}
		// Noop
		return bs, nil
	case CHAR, VARCHAR, TEXT, LONGCHAR:
		if len(bs) == 0 {
			return nil, nil
		}
		return string(bs), nil
	case UNITEXT:
		runes := []rune{}

		for i := 0; i < len(bs); i++ {
			// Determine if byte is a utf16 surrogate - if so two
			// bytes must be consumed to form one utf16 code point
			if utf16.IsSurrogate(rune(bs[i])) {
				r := utf16.DecodeRune(rune(bs[i]), rune(bs[i+1]))
				runes = append(runes, r)
				i++
			} else {
				runes = append(runes, rune(bs[i]))
			}
		}

		s := string(runes)
		// Trim null bytes from the right - ASE always sends the
		// maximum bytes for the TEXT datatype, causing the string
		// to have a couple thousand null bytes. These are also
		// carried over in a string() conversion and cause
		// false-negatives in comparisons.
		s = strings.TrimRight(s, "\x00")

		return s, nil
	case SHORTMONEY, MONEY, MONEYN:
		dec, err := NewDecimal(0, 0)
		if err != nil {
			return nil, fmt.Errorf("error creating decimal: %w", err)
		}

		switch len(bs) {
		case 0:
			return &Decimal{i: nil}, nil
		case 4:
			dec.Precision = ASEShortMoneyPrecision
			dec.Scale = ASEShortMoneyScale

			dec.SetInt64(int64(int32(endian.Uint32(bs))))
		case 8:
			dec.Precision = ASEMoneyPrecision
			dec.Scale = ASEMoneyScale

			mnyhigh := endian.Uint32(bs[:4])
			mnylow := endian.Uint32(bs[4:])

			mny := int64(int64(mnyhigh)<<32 + int64(mnylow))
			dec.SetInt64(mny)
		}

		return dec, nil
	case DECN, NUMN:
		if len(bs) == 0 {
			return &Decimal{i: nil}, nil
		}

		dec, err := NewDecimal(ASEDecimalDefaultPrecision, ASEDecimalDefaultScale)
		if err != nil {
			return nil, fmt.Errorf("error creating decimal: %w", err)
		}

		dec.SetBytes(bs[1:])
		if bs[0] == 0x1 {
			dec.Negate()
		}

		// User must set precision and scale
		return dec, nil
	case DATE, DATEN:
		if len(bs) == 0 {
			return nil, nil
		}

		x := int32(endian.Uint32(bs))
		days := asetime.ASEDuration(x) * asetime.Day
		return asetime.Epoch1900().AddDate(0, 0, days.Days()), nil
	case TIME, BIGTIMEN, TIMEN:
		switch len(bs) {
		case 0: // Null
			return nil, nil
		case 4: // TIME
			x := int(int32(endian.Uint32(bs)))
			dur := asetime.FractionalSecondToMillisecond(x)
			t := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
			return t.Add(time.Duration(dur.Milliseconds()) * time.Millisecond), nil
		case 8: // BIGTIME
			dur := asetime.ASEDuration(endian.Uint64(bs))

			t := asetime.EpochRataDie()
			t = t.Add(time.Duration(dur) * time.Microsecond)

			return t, nil
		default:
			return nil, fmt.Errorf("invalid length for %v: %d", t, len(bs))
		}
	case SHORTDATE, DATETIME, DATETIMEN:
		switch len(bs) {
		case 0:
			return nil, nil
		case 4:
			days := endian.Uint16(bs[:2])
			mins := endian.Uint16(bs[2:])

			t := asetime.Epoch1900()
			t = t.AddDate(0, 0, int(days))
			t = t.Add(time.Duration(int(mins)) * time.Minute)
			return t, nil
		case 8:
			days := asetime.ASEDuration(int32(endian.Uint32(bs[:4]))) * asetime.Day
			ms := asetime.FractionalSecondToMillisecond(int(endian.Uint32(bs[4:])))

			t := asetime.Epoch1900()
			t = t.AddDate(0, 0, days.Days())
			t = t.Add(time.Duration(ms.Microseconds()) * time.Microsecond)

			return t, nil
		default:
			return nil, fmt.Errorf("invalid length for %v: %d", t, len(bs))
		}
	case BIGDATETIMEN:
		if len(bs) == 0 {
			return nil, nil
		}

		dur := asetime.ASEDuration(endian.Uint64(bs))

		t := time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC)
		t = t.AddDate(0, 0, dur.Days())
		ms := dur.Microseconds() - (dur.Days() * int(asetime.Day))
		t = t.Add(time.Duration(ms) * time.Microsecond)

		return t, nil
	default:
		return nil, fmt.Errorf("unhandled data type %s", t)
	}
}
