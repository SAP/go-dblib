// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
// SPDX-FileCopyrightText: 2023 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package asetypes

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// Default properties for ASE data type 'decimal'.
const (
	ASEDecimalDefaultPrecision = 18
	ASEDecimalDefaultScale     = 0

	ASEMoneyPrecision = 20
	ASEMoneyScale     = 4

	ASEShortMoneyPrecision = 10
	ASEShortMoneyScale     = 4

	aseMaxDecimalDigits = 38
)

// Errors of ASE data type 'decimal' operations.
var (
	ErrDecimalPrecisionTooHigh         = fmt.Errorf("precision is set to more than %d digits", aseMaxDecimalDigits)
	ErrDecimalPrecisionTooLow          = fmt.Errorf("precision is set to less than 0 digits")
	ErrDecimalScaleTooHigh             = fmt.Errorf("scale is set to more than %d digits", aseMaxDecimalDigits)
	ErrDecimalScaleBiggerThanPrecision = fmt.Errorf("scale is bigger then precision")
)

// Decimal only carries the information of Decimal, Numeric and Money
// ASE datatypes. This is only sufficient for displaying, not
// calculations.
type Decimal struct {
	Precision, Scale int
	i                *big.Int
}

// NewDecimal creates a new decimal with the passed precision and scale
// and returns it.
// An error is returned if the precision/scale combination is not valid.
func NewDecimal(precision, scale int) (*Decimal, error) {
	dec := &Decimal{
		Precision: precision,
		Scale:     scale,
		i:         new(big.Int),
	}

	if err := dec.sanity(); err != nil {
		return nil, err
	}

	return dec, nil
}

// NewDecimalString creates a new decimal based on the passed string.
// If the string contains an invalid precision/scale combination an
// error is returned.
func NewDecimalString(precision, scale int, s string) (*Decimal, error) {
	dec, err := NewDecimal(precision, scale)
	if err != nil {
		return nil, fmt.Errorf("error creating decimal: %w", err)
	}

	if err := dec.SetString(s); err != nil {
		return nil, fmt.Errorf("error setting string: %w", err)
	}

	return dec, nil
}

func (dec Decimal) sanity() error {
	if dec.Precision > aseMaxDecimalDigits {
		return ErrDecimalPrecisionTooHigh
	}

	if dec.Precision < 0 {
		return ErrDecimalPrecisionTooLow
	}

	if dec.Scale > aseMaxDecimalDigits {
		return ErrDecimalScaleTooHigh
	}

	if dec.Scale > dec.Precision {
		return ErrDecimalScaleBiggerThanPrecision
	}

	return nil
}

// Cmp compares precision, scale, and value between two decimals.
func (dec Decimal) Cmp(other Decimal) bool {
	if dec.Precision != other.Precision {
		return false
	}

	if dec.Scale != other.Scale {
		return false
	}

	return dec.i.Cmp(other.i) == 0
}

// IsNegative returns true if dec is negative.
func (dec Decimal) IsNegative() bool {
	return dec.i.Sign() < 0
}

// Negate sets dec to the value of it with its sign negated.
func (dec *Decimal) Negate() {
	dec.i.Neg(dec.i)
}

// Bytes returns the byte-slice of dec.
func (dec Decimal) Bytes() []byte {
	return dec.i.Bytes()
}

// ByteSize calls DecimalByteSize.
func (dec Decimal) ByteSize() int {
	// Source: https://github.com/thda/tds/blob/master/num.go#L248
	return int(math.Ceil(float64(dec.i.BitLen())/8) + 1)
}

// SetInt64 sets dec.i to i.
func (dec *Decimal) SetInt64(i int64) {
	dec.i.SetInt64(i)
}

// Int returns a copy of the underlying big.Int.
func (dec Decimal) Int() *big.Int {
	i := &big.Int{}
	i.Add(i, dec.i)
	return i
}

// SetBytes interprets b as the bytes of a big-endian unsigned integer
// and sets dec to that values.
func (dec *Decimal) SetBytes(b []byte) {
	dec.i.SetBytes(b)
}

func (dec *Decimal) String() string {
	if dec.i == nil {
		return "<nil>"
	}

	s := fmt.Sprintf("%0"+strconv.Itoa(dec.Precision)+"s", big.NewInt(0).Abs(dec.i))

	neg := ""
	if dec.IsNegative() {
		neg = "-"
	}

	right := strings.TrimRight(s[dec.Precision-dec.Scale:], "0")
	if len(right) == 0 {
		right = "0"
	}

	left := strings.TrimLeft(s[:dec.Precision-dec.Scale], "0")
	if len(left) == 0 {
		left = "0"
	}

	ret := fmt.Sprintf("%s%s.%s", neg, left, right)

	return ret
}

// Set decimal to the passed string value.
// Precision and scale are untouched.
//
// If an error is returned dec is untouched.
func (dec *Decimal) SetString(s string) error {
	// Trim spaces to avoid errors with "+0.0 " etc.pp.
	s = strings.TrimSpace(s)

	split := strings.Split(s, ".")
	left := split[0]
	right := ""
	if len(split) > 1 {
		right = split[1]
	}

	// Set underlying big.Int structure to the whole number
	i := &big.Int{}
	if _, ok := i.SetString(left+right, 10); !ok {
		return fmt.Errorf("failed to parse number %s%s", left, right)
	}

	// Multiply underlying big.Int to fit to the scale of the decimal
	if dec.Scale-len(right) > 0 {
		mul := big.NewInt(10)
		mul.Exp(mul, big.NewInt(int64(dec.Scale-len(right))), nil)
		i.Mul(i, mul)
	}

	dec.i = i
	return nil
}
