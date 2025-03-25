// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package asetypes

import (
	"database/sql/driver"
	"fmt"
)

// NullableTypes maps certain non-nullable types to their
// respective nullable type.
var nullableTypes = map[DataType]DataType{
	DATE:       DATEN,
	DATETIME:   DATETIMEN,
	SHORTDATE:  DATETIMEN,
	FLT4:       FLTN,
	FLT8:       FLTN,
	INT1:       INTN,
	INT2:       INTN,
	INT4:       INTN,
	INT8:       INTN,
	MONEY:      MONEYN,
	SHORTMONEY: MONEYN,
	TIME:       TIMEN,
	UINT2:      UINTN,
	UINT4:      UINTN,
	UINT8:      UINTN,
}

// NullableType returns the respective nullable type of t, according to
// the NullableTypes map.
func (t DataType) NullableType() (DataType, error) {
	// If true, the passed DataType is already a nullable type
	if t.ByteSize() == -1 {
		return t, nil
	}
	nullableType, ok := nullableTypes[t]
	if !ok {
		return t, fmt.Errorf("datatype %q has no nullable equivalent", t)
	}
	return nullableType, nil
}

// Below are several custom nullable datatypes implemented that are not
// provided by database/sql.Null<type>. The datatypes implement the
// Scanner and driver Valuer interface provided by database/sql.

type NullInt16 struct {
	Int16 int16
	Valid bool
}

// Scan implements the Scanner interface.
func (ni *NullInt16) Scan(value interface{}) error {
	if value == nil {
		ni.Int16 = 0
		ni.Valid = false
		return nil
	}

	v, ok := value.(int16)
	if !ok {
		return fmt.Errorf("asetypes: %T is not of type int16", value)
	}
	ni.Int16 = v
	ni.Valid = true

	return nil
}

// Value implements the driver Valuer interface.
func (ni NullInt16) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Int16, nil
}

type NullUint64 struct {
	Uint64 uint64
	Valid  bool
}

// Scan implements the Scanner interface.
func (nu *NullUint64) Scan(value interface{}) error {
	if value == nil {
		nu.Uint64 = 0
		nu.Valid = false
		return nil
	}

	v, ok := value.(uint64)
	if !ok {
		return fmt.Errorf("asetypes: %T is not of type uint64", value)
	}
	nu.Uint64 = v
	nu.Valid = true

	return nil
}

// Value implements the driver Valuer interface.
func (ni NullUint64) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Uint64, nil
}

type NullUint32 struct {
	Uint32 uint32
	Valid  bool
}

// Scan implements the Scanner interface.
func (nu *NullUint32) Scan(value interface{}) error {
	if value == nil {
		nu.Uint32 = 0
		nu.Valid = false
		return nil
	}

	v, ok := value.(uint32)
	if !ok {
		return fmt.Errorf("asetypes: %T is not of type uint32", value)
	}
	nu.Uint32 = v
	nu.Valid = true

	return nil
}

// Value implements the driver Valuer interface.
func (ni NullUint32) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Uint32, nil
}

type NullUint16 struct {
	Uint16 uint16
	Valid  bool
}

// Scan implements the Scanner interface.
func (nu *NullUint16) Scan(value interface{}) error {
	if value == nil {
		nu.Uint16 = 0
		nu.Valid = false
		return nil
	}

	v, ok := value.(uint16)
	if !ok {
		return fmt.Errorf("asetypes: %T is not of type uint16", value)
	}
	nu.Uint16 = v
	nu.Valid = true

	return nil
}

// Value implements the driver Valuer interface.
func (ni NullUint16) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Uint16, nil
}

type NullUint8 struct {
	Uint8 uint8
	Valid bool
}

// Scan implements the Scanner interface.
func (nu *NullUint8) Scan(value interface{}) error {
	if value == nil {
		nu.Uint8 = 0
		nu.Valid = false
		return nil
	}

	v, ok := value.(uint8)
	if !ok {
		return fmt.Errorf("asetypes: %T is not of type uint8", value)
	}
	nu.Uint8 = v
	nu.Valid = true

	return nil
}

// Value implements the driver Valuer interface.
func (nu NullUint8) Value() (driver.Value, error) {
	if !nu.Valid {
		return nil, nil
	}
	return nu.Uint8, nil
}

type NullDecimal struct {
	Dec   *Decimal
	Valid bool
}

type NullFloat32 struct {
	Float32 float32
	Valid   bool
}

// Scan implements the Scanner interface.
func (nf *NullFloat32) Scan(value interface{}) error {
	if value == nil {
		nf.Float32 = 0
		nf.Valid = false
		return nil
	}

	v, ok := value.(float32)
	if !ok {
		return fmt.Errorf("asetypes: %T is not of type float32", value)
	}
	nf.Float32 = v
	nf.Valid = true

	return nil
}

// Value implements the driver Valuer interface.
func (nf NullFloat32) Value() (driver.Value, error) {
	if !nf.Valid {
		return nil, nil
	}
	return nf.Float32, nil
}

// Scan implements the Scanner interface.
func (nd *NullDecimal) Scan(value interface{}) error {
	if value == nil {
		nd.Dec = &Decimal{}
		nd.Valid = false
		return nil
	}

	v, ok := value.(*Decimal)
	if !ok {
		return fmt.Errorf("asetypes: %T is not of type decimal", value)
	}
	nd.Dec = v

	if nd.Dec.i == nil {
		nd.Valid = false
		return nil
	}
	nd.Valid = true

	return nil
}

// Value implements the driver Valuer interface.
func (nd NullDecimal) Value() (driver.Value, error) {
	if !nd.Valid {
		return nil, nil
	}
	return nd.Dec, nil
}

type NullBinary struct {
	ByteSlice []byte
	Valid     bool
}

// Scan implements the Scanner interface.
func (nb *NullBinary) Scan(value interface{}) error {
	if value == nil {
		nb.ByteSlice = []byte{}
		nb.Valid = false
		return nil
	}

	v, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("asetypes: %T is not of type []byte", value)
	}
	nb.ByteSlice = v
	nb.Valid = true

	return nil
}

// Value implements the driver Valuer interface.
func (nb NullBinary) Value() (driver.Value, error) {
	if !nb.Valid {
		return nil, nil
	}
	return nb.ByteSlice, nil
}
