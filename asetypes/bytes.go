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
	"time"
	"unicode/utf16"

	"github.com/SAP/go-dblib/asetime"
)

// Bytes returns a byte slice based on a given value-interface and depending
// on the ASE data type.
// TODO: Instead of parameter 'length', one could use a struct to store
// additional optional parameters, e.g. 'length', 'datastatus', ...
// (Problem: golang mod import cycle)
func (t DataType) Bytes(endian binary.ByteOrder, value interface{}, length int64) ([]byte, error) {
	// If value is nil, then immediately return an empty byteslice
	if value == nil {
		bs := make([]byte, 0)
		return bs, nil
	}

	switch t {
	case MONEY, SHORTMONEY, MONEYN:
		dec, ok := value.(*Decimal)
		if !ok {
			return nil, fmt.Errorf("expected *asetypes.Decimal for %s, received %T", t, value)
		}
		deci := dec.Int()

		bs := make([]byte, length)
		switch length {
		case 4: // SHORTMONEY, MONEYN(4)
			endian.PutUint32(bs, uint32(deci.Int64()))
		case 8: // MONEY, MONEYN(8)
			endian.PutUint32(bs[:4], uint32(deci.Int64()>>32))
			endian.PutUint32(bs[4:], uint32(deci.Int64()))
		}

		return bs, nil
	case DECN, NUMN:
		dec, ok := value.(*Decimal)
		if !ok {
			return nil, fmt.Errorf("expected *asetypes.Decimal for %s, received %T", t, value)
		}

		bs := make([]byte, dec.ByteSize())
		copy(bs[dec.ByteSize()-len(dec.Bytes()):], dec.Bytes())
		if dec.IsNegative() {
			bs[0] = 0x1
		}
		return bs, nil
	case DATE, DATEN:
		t := asetime.DurationFromDateTime(value.(time.Time))
		t -= asetime.DurationFromDateTime(asetime.Epoch1900())

		bs := make([]byte, length)
		endian.PutUint32(bs, uint32(t.Days()))
		return bs, nil
	case TIME, TIMEN:
		dur := asetime.DurationFromTime(value.(time.Time))
		fract := asetime.MillisecondToFractionalSecond(dur.Microseconds())

		bs := make([]byte, length)
		endian.PutUint32(bs, uint32(fract))
		return bs, nil
	case SHORTDATE, DATETIME, DATETIMEN:
		t := asetime.DurationFromDateTime(value.(time.Time))
		t -= asetime.DurationFromDateTime(asetime.Epoch1900())

		days := t.Days()

		bs := make([]byte, length)
		switch length {
		case 4: // SHORTDATE/DATETIME4, DATETIMEN(4)
			s := asetime.ASEDuration(t.Microseconds() - days*int(asetime.Day))
			binary.LittleEndian.PutUint16(bs[:2], uint16(days))
			binary.LittleEndian.PutUint16(bs[2:], uint16(s.Minutes()))
		case 8: // DATETIME, DATETIMEN(8)
			s := t.Microseconds() - days*int(asetime.Day)
			s = asetime.MillisecondToFractionalSecond(s)
			binary.LittleEndian.PutUint32(bs[:4], uint32(days))
			binary.LittleEndian.PutUint32(bs[4:], uint32(s))
		}
		return bs, nil
	case BIGDATETIMEN:
		dur := asetime.DurationFromDateTime(value.(time.Time))

		bs := make([]byte, length)
		binary.LittleEndian.PutUint64(bs, uint64(dur))
		return bs, nil
	case BIGTIMEN:
		dur := asetime.DurationFromTime(value.(time.Time))

		bs := make([]byte, length)
		binary.LittleEndian.PutUint64(bs, uint64(dur))
		return bs, nil
	case UNITEXT:
		// convert go string to utf16 code points
		runes := []rune(value.(string))
		utf16bytes := utf16.Encode(runes)

		// convert utf16 code points to bytes
		bs := make([]byte, len(utf16bytes)*2)
		for i := 0; i < len(utf16bytes); i++ {
			binary.LittleEndian.PutUint16(bs[i:], utf16bytes[i])
		}

		return bs, nil
	}

	switch typed := value.(type) {
	case string:
		value = []byte(typed)
	}

	buf := &bytes.Buffer{}
	if err := binary.Write(buf, endian, value); err != nil {
		return nil, fmt.Errorf("error writing value: %w", err)
	}

	bs := buf.Bytes()
	if t.ByteSize() != -1 && t.ByteSize() != len(bs) {
		return nil, fmt.Errorf("binary.Write returned a byteslice of length %d, expected %d for datatype %s",
			len(bs), t.ByteSize(), t)
	}

	return bs, nil
}
