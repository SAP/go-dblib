// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"bytes"
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/SAP/go-dblib/asetypes"
)

//go:generate go run ./gen_type.go BigInt int64
var samplesBigInt = []int64{
	math.MinInt64, math.MaxInt64,
	-5000, -100, 0, 100, 5000,
}

//go:generate go run ./gen_type.go Int int32
var samplesInt = []int32{
	math.MinInt32, math.MaxInt32,
	-5000, -100, 0, 100, 5000,
}

//go:generate go run ./gen_type.go SmallInt int16
var samplesSmallInt = []int16{-32768, 0, 32767}

//go:generate go run ./gen_type.go TinyInt uint8
var samplesTinyInt = []uint8{0, 255}

//go:generate go run ./gen_type.go NullBigInt sql.NullInt64 -columndef "bigint null"
var samplesNullBigInt = []sql.NullInt64{
	{Int64: math.MinInt64, Valid: true},
	{Int64: 0, Valid: true},
	{Valid: false},
	{Int64: math.MaxInt64, Valid: true},
}

//go:generate go run ./gen_type.go NullInt sql.NullInt32 -columndef "int null"
var samplesNullInt = []sql.NullInt32{
	{Int32: math.MinInt32, Valid: true},
	{Int32: 0, Valid: true},
	{Valid: false},
	{Int32: math.MaxInt32, Valid: true},
}

//go:generate go run ./gen_type.go NullSmallInt github.com/SAP/go-dblib/asetypes.NullInt16 -columndef "smallint null"
var samplesNullSmallInt = []asetypes.NullInt16{
	{Int16: math.MinInt16, Valid: true},
	{Int16: 0, Valid: true},
	{Valid: false},
	{Int16: math.MaxInt16, Valid: true},
}

//go:generate go run ./gen_type.go NullTinyInt github.com/SAP/go-dblib/asetypes.NullUint8 -columndef "tinyint null"
var samplesNullTinyInt = []asetypes.NullUint8{
	{Uint8: 0, Valid: true},
	{Valid: false},
	{Uint8: math.MaxUint8, Valid: true},
}

//go:generate go run ./gen_type.go UnsignedBigInt uint64 -columndef "unsigned bigint"
var samplesUnsignedBigInt = []uint64{0, 1000, 5000, 150000, 123456789, math.MaxUint64}

//go:generate go run ./gen_type.go UnsignedInt uint32 -columndef "unsigned int"
var samplesUnsignedInt = []uint32{0, 1000, 5000, 150000, 123456789, math.MaxUint32}

//go:generate go run ./gen_type.go UnsignedSmallInt uint16 -columndef "unsigned smallint"
var samplesUnsignedSmallInt = []uint16{0, 65535}

//go:generate go run ./gen_type.go NullUnsignedBigInt github.com/SAP/go-dblib/asetypes.NullUint64 -columndef "unsigned bigint null"
var samplesNullUnsignedBigInt = []asetypes.NullUint64{
	{Uint64: 0, Valid: true},
	{Valid: false},
	{Uint64: math.MaxUint64, Valid: true},
}

//go:generate go run ./gen_type.go NullUnsignedInt github.com/SAP/go-dblib/asetypes.NullUint32 -columndef "unsigned int null"
var samplesNullUnsignedInt = []asetypes.NullUint32{
	{Uint32: 0, Valid: true},
	{Valid: false},
	{Uint32: math.MaxUint32, Valid: true},
}

//go:generate go run ./gen_type.go NullUnsignedSmallInt github.com/SAP/go-dblib/asetypes.NullUint16 -columndef "unsigned smallint null"
var samplesNullUnsignedSmallInt = []asetypes.NullUint16{
	{Uint16: 0, Valid: true},
	{Valid: false},
	{Uint16: math.MaxUint16, Valid: true},
}

func convertDecimal10(sample string) (*asetypes.Decimal, error) {
	return asetypes.NewDecimalString(1, 0, sample)
}

//go:generate go run ./gen_type.go Decimal10 github.com/SAP/go-dblib/*asetypes.Decimal -columndef decimal(1,0) -convert convertDecimal10 -compare compareDecimal
var samplesDecimal10 = []string{"0", "1", "9"}

func convertDecimal380(sample string) (*asetypes.Decimal, error) {
	return asetypes.NewDecimalString(38, 0, sample)
}

//go:generate go run ./gen_type.go Decimal380 github.com/SAP/go-dblib/*asetypes.Decimal -columndef decimal(38,0) -convert convertDecimal380 -compare compareDecimal
var samplesDecimal380 = []string{"99999999999999999999999999999999999999"}

func convertDecimal3838(sample string) (*asetypes.Decimal, error) {
	return asetypes.NewDecimalString(38, 38, sample)
}

//go:generate go run ./gen_type.go Decimal3838 github.com/SAP/go-dblib/*asetypes.Decimal -columndef decimal(38,38) -convert convertDecimal3838 -compare compareDecimal
var samplesDecimal3838 = []string{".99999999999999999999999999999999999999"}

func convertDecimal3819(sample string) (*asetypes.Decimal, error) {
	return asetypes.NewDecimalString(38, 19, sample)
}

//go:generate go run ./gen_type.go Decimal github.com/SAP/go-dblib/*asetypes.Decimal -columndef "decimal(38,19)" -convert convertDecimal3819 -compare compareDecimal
var samplesDecimal = []string{
	// ASE max
	"1234567890123456789",
	"9999999999999999999",
	"-1234567890123456789",
	"-9999999999999999999",
	// ASE min
	".1234567890123456789",
	".9999999999999999999",
	"-.1234567890123456789",
	"-.9999999999999999999",
	// default
	"0",
	// arbitrary
	"1234.5678",
}

func compareDecimal(recv, expect *asetypes.Decimal) bool {
	return !expect.Cmp(*recv)
}

func convertNullDecimal(sample sql.NullString) (asetypes.NullDecimal, error) {
	var nd asetypes.NullDecimal
	if !sample.Valid {
		nd.Valid = false
		return nd, nil
	}

	dec, err := asetypes.NewDecimalString(38, 19, sample.String)
	if err != nil {
		return nd, fmt.Errorf("integration: error creating new decimal from string '%s': %w", sample.String, err)
	}
	if err := nd.Scan(dec); err != nil {
		return nd, fmt.Errorf("integration: error scanning decimal %+v into NullDecimal: %w", dec, err)
	}

	return nd, nil
}

func compareNullDecimal(recv, expect asetypes.NullDecimal) bool {
	if recv.Valid != expect.Valid {
		return true
	}

	// If recv or expect is !Valid, the potential decimal-value could be
	// nil, which confuses math/big.Int.Cmp().
	if !recv.Valid {
		return false
	}

	return !expect.Dec.Cmp(*recv.Dec)
}

//go:generate go run ./gen_type.go NullDecimal github.com/SAP/go-dblib/asetypes.NullDecimal -columndef "decimal(38,19) null" -convert convertNullDecimal -compare compareNullDecimal
var samplesNullDecimal = []sql.NullString{
	{String: "-9999999999999999999", Valid: true},
	{String: "-.9999999999999999999", Valid: true},
	{String: "0", Valid: true},
	{String: ".9999999999999999999", Valid: true},
	{String: "9999999999999999999", Valid: true},
	{Valid: false},
	{String: "1234.56789", Valid: true},
}

//go:generate go run ./gen_type.go Float float64
var samplesFloat = []float64{
	-math.SmallestNonzeroFloat64,
	math.SmallestNonzeroFloat64,
	-1000,
	1000,
	-math.MaxFloat64,
	math.MaxFloat64,
}

//go:generate go run ./gen_type.go Real float32
var samplesReal = []float32{
	-math.SmallestNonzeroFloat32,
	math.SmallestNonzeroFloat32,
	-1000,
	1000,
	-math.MaxFloat32,
	math.MaxFloat32,
}

//go:generate go run ./gen_type.go NullFloat sql.NullFloat64 -columndef "float null"
var samplesNullFloat = []sql.NullFloat64{
	{Float64: -math.SmallestNonzeroFloat64, Valid: true},
	{Float64: math.SmallestNonzeroFloat64, Valid: true},
	{Float64: -1000, Valid: true},
	{Float64: 0, Valid: true},
	{Valid: false},
	{Float64: 1000, Valid: true},
	{Float64: -math.MaxFloat64, Valid: true},
	{Float64: math.MaxFloat64, Valid: true},
}

//go:generate go run ./gen_type.go NullReal github.com/SAP/go-dblib/asetypes.NullFloat32 -columndef "real null"
var samplesNullReal = []asetypes.NullFloat32{
	{Float32: -math.SmallestNonzeroFloat32, Valid: true},
	{Float32: math.SmallestNonzeroFloat32, Valid: true},
	{Float32: -1000, Valid: true},
	{Float32: 0, Valid: true},
	{Valid: false},
	{Float32: 1000, Valid: true},
	{Float32: -math.MaxFloat32, Valid: true},
	{Float32: math.MaxFloat32, Valid: true},
}

func convertMoney(sample string) (*asetypes.Decimal, error) {
	return asetypes.NewDecimalString(asetypes.ASEMoneyPrecision, asetypes.ASEMoneyScale, sample)
}

//go:generate go run ./gen_type.go Money github.com/SAP/go-dblib/*asetypes.Decimal -convert convertMoney -compare compareDecimal
var samplesMoney = []string{
	// ASE min
	"-922337203685477.5808",
	// ASE max
	"922337203685477.5807",
	// default
	"0.0",
	// arbitrary
	"1234.5678",
}

func convertSmallMoney(sample string) (*asetypes.Decimal, error) {
	return asetypes.NewDecimalString(asetypes.ASEShortMoneyPrecision, asetypes.ASEShortMoneyScale, sample)
}

//go:generate go run ./gen_type.go Money4 github.com/SAP/go-dblib/*asetypes.Decimal -columndef smallmoney -convert convertSmallMoney -compare compareDecimal
var samplesMoney4 = []string{
	// ASE min
	"-214748.3648",
	// ASE max
	"214748.3647",
	// default
	"0.0",
	// arbitrary
	"1234.5678",
}

func convertNullMoney(sample sql.NullString) (asetypes.NullDecimal, error) {
	var nd asetypes.NullDecimal
	if !sample.Valid {
		nd.Valid = false
		return nd, nil
	}

	dec, err := asetypes.NewDecimalString(asetypes.ASEMoneyPrecision, asetypes.ASEMoneyScale, sample.String)
	if err != nil {
		return nd, fmt.Errorf("integration: error creating new decimal from string '%s': %w", sample.String, err)
	}
	if err := nd.Scan(dec); err != nil {
		return nd, fmt.Errorf("integration: error scanning decimal %+v into NullDecimal: %w", dec, err)
	}

	return nd, nil
}

//go:generate go run ./gen_type.go NullMoney github.com/SAP/go-dblib/asetypes.NullDecimal -columndef "money null" -convert convertNullMoney -compare compareNullDecimal
var samplesNullMoney = []sql.NullString{
	// ASE min
	{String: "-922337203685477.5808", Valid: true},
	// ASE max
	{String: "922337203685477.5807", Valid: true},
	// default
	{String: "0.0", Valid: true},
	// Null
	{Valid: false},
	// arbitrary
	{String: "1234.5678", Valid: true},
}

func convertNullMoney4(sample sql.NullString) (asetypes.NullDecimal, error) {
	var nd asetypes.NullDecimal
	if !sample.Valid {
		nd.Valid = false
		return nd, nil
	}

	dec, err := asetypes.NewDecimalString(asetypes.ASEShortMoneyPrecision, asetypes.ASEShortMoneyScale, sample.String)
	if err != nil {
		return nd, fmt.Errorf("integration: error creating new decimal from string '%s': %w", sample.String, err)
	}
	if err := nd.Scan(dec); err != nil {
		return nd, fmt.Errorf("integration: error scanning decimal %+v into NullDecimal: %w", dec, err)
	}

	return nd, nil
}

//go:generate go run ./gen_type.go NullMoney4 github.com/SAP/go-dblib/asetypes.NullDecimal -columndef "smallmoney null" -convert convertNullMoney4 -compare compareNullDecimal
var samplesNullMoney4 = []sql.NullString{
	// ASE min
	{String: "-214748.3648", Valid: true},
	// ASE max
	{String: "214748.3647", Valid: true},
	// default
	{String: "0.0", Valid: true},
	// Null
	{Valid: false},
	// arbitrary
	{String: "1234.5678", Valid: true},
}

//go:generate go run ./gen_type.go Date time.Time
var samplesDate = []time.Time{
	// Sybase & Golang zero value
	time.Time{},
	// Sybase max
	time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC),
}

//go:generate go run ./gen_type.go Time time.Time
var samplesTime = []time.Time{
	// Sybase & Golang zero-value; 00:00:00.00
	time.Time{},
	// 13:15:55.123
	time.Date(1, time.January, 1, 13, 15, 55, 123000000, time.UTC),
	// Sybase max: 23:59:59.990
	time.Date(1, time.January, 1, 23, 59, 59, 996000000, time.UTC),
}

//go:generate go run ./gen_type.go BigTime time.Time
var samplesBigTime = []time.Time{
	// Sybase & Golang zero-value; 00:00:00.00
	time.Time{},
	// Sybase max: 23:59:59.999999
	time.Date(1, time.January, 1, 23, 59, 59, 999999000, time.UTC),
}

//go:generate go run ./gen_type.go SmallDateTime time.Time
var samplesSmallDateTime = []time.Time{
	// Sybase zero-value; January 1, 1900 Midnight
	time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC),
	// Sybase max: 06.06.2079 23:59
	time.Date(2079, time.June, 6, 23, 59, 0, 0, time.UTC),
}

//go:generate go run ./gen_type.go DateTime time.Time
var samplesDateTime = []time.Time{
	// Sybase min: January 1, 1753 Midnight
	time.Date(1753, time.January, 1, 0, 0, 0, 0, time.UTC),
	// Sybase zero-value; January 1, 1900 Midnight
	time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC),
	// Sybase max: 31.12.9999 23:59:59.996
	time.Date(9999, time.December, 31, 23, 59, 59, 996000000, time.UTC),
}

//go:generate go run ./gen_type.go BigDateTime time.Time
var samplesBigDateTime = []time.Time{
	// Sybase & Golang zero-value; January 1, 0001 Midnight
	time.Time{},
	time.Date(2019, time.March, 29, 9, 26, 0, 0, time.UTC),
	// Sybase max
	time.Date(9999, time.December, 31, 23, 59, 59, 999999000, time.UTC),
}

//go:generate go run ./gen_type.go NullDate sql.NullTime -columndef "date null"
var samplesNullDate = []sql.NullTime{
	// Sybase & Golang zero value
	{Time: time.Time{}, Valid: true},
	// Null
	{Valid: false},
	// Sybase max
	{Time: time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC), Valid: true},
}

//go:generate go run ./gen_type.go NullTime sql.NullTime -columndef "time null"
var samplesNullTime = []sql.NullTime{
	// Sybase & Golang zero-value; 00:00:00.00
	{Time: time.Time{}, Valid: true},
	// 13:15:55.123
	{Time: time.Date(1, time.January, 1, 13, 15, 55, 123000000, time.UTC), Valid: true},
	// Null
	{Valid: false},
	// Sybase max: 23:59:59.990
	{Time: time.Date(1, time.January, 1, 23, 59, 59, 996000000, time.UTC), Valid: true},
}

//go:generate go run ./gen_type.go NullBigTime sql.NullTime -columndef "bigtime null"
var samplesNullBigTime = []sql.NullTime{
	// Sybase & Golang zero-value; 00:00:00.00
	{Time: time.Time{}, Valid: true},
	// Null
	{Valid: false},
	// Sybase max: 23:59:59.999999
	{Time: time.Date(1, time.January, 1, 23, 59, 59, 999999000, time.UTC), Valid: true},
}

//go:generate go run ./gen_type.go NullSmallDateTime sql.NullTime -columndef "smalldatetime null"
var samplesNullSmallDateTime = []sql.NullTime{
	// Sybase zero-value; January 1, 1900 Midnight
	{Time: time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: true},
	// Null
	{Valid: false},
	// Sybase max: 06.06.2079 23:59
	{Time: time.Date(2079, time.June, 6, 23, 59, 0, 0, time.UTC), Valid: true},
}

//go:generate go run ./gen_type.go NullDateTime sql.NullTime -columndef "datetime null"
var samplesNullDateTime = []sql.NullTime{
	// Sybase min: January 1, 1753 Midnight
	{Time: time.Date(1753, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: true},
	// Sybase zero-value; January 1, 1900 Midnight
	{Time: time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: true},
	// Null
	{Valid: false},
	// Sybase max: 31.12.9999 23:59:59.996
	{Time: time.Date(9999, time.December, 31, 23, 59, 59, 996000000, time.UTC), Valid: true},
}

//go:generate go run ./gen_type.go NullBigDateTime sql.NullTime -columndef "bigdatetime null"
var samplesNullBigDateTime = []sql.NullTime{
	// Sybase & Golang zero-value; January 1, 0001 Midnight
	{Time: time.Time{}, Valid: true},
	{Time: time.Date(2019, time.March, 29, 9, 26, 0, 0, time.UTC), Valid: true},
	// Null
	{Valid: false},
	// Sybase max
	{Time: time.Date(9999, time.December, 31, 23, 59, 59, 999999000, time.UTC), Valid: true},
}

//go:generate go run ./gen_type.go Char string -columndef "char(13)" -compare compareChar
var samplesChar = []string{" ", "test", "a longer test"}

// Pad with whitespaces to 13 chars (-> char(13))
func compareChar(recv, expect string) bool {
	return strings.Compare(recv, fmt.Sprintf("%-13s", expect)) != 0
}

//go:generate go run ./gen_type.go NChar string -columndef "nchar(13)"  -compare compareNChar
var samplesNChar = samplesChar

// Pad with whitespaces to 39 chars (nchar(13) x @@ncharsize == 3 -> 13 x 3 = 39)
// (@@ncharsize is depending on the average national character length)
func compareNChar(recv, expect string) bool {
	return strings.Compare(recv, fmt.Sprintf("%-39s", expect)) != 0
}

//go:generate go run ./gen_type.go VarChar string -columndef "varchar(13)" -compare compareVarChar
var samplesVarChar = samplesChar

//go:generate go run ./gen_type.go NVarChar string -columndef "nvarchar(13)" -compare compareVarChar
var samplesNVarChar = samplesChar

func compareVarChar(recv, expect string) bool {
	return strings.Compare(recv, expect) != 0
}

//go:generate go run ./gen_type.go NullChar sql.NullString -columndef "char(13) null" -compare compareNullChar
var samplesNullChar = []sql.NullString{
	{String: "", Valid: true},
	{String: "test", Valid: true},
	{Valid: false},
	{String: "a longer test", Valid: true},
}

func compareNullChar(recv, expect sql.NullString) bool {
	// Special case: Non-null, zero-length strings
	// ASE: '""' -> '" "'
	// go-ase: '""' -> NULL (default)
	//   Example: sql.NullString{String: "", Valid: true} -> sql.NullString{Valid: false}
	if expect.String == "" && expect.Valid {
		return recv.Valid
	}

	if recv.Valid != expect.Valid {
		return true
	}

	// No whitespace-padding required, since values are inserted as
	// non-fixed length datatype
	return strings.Compare(recv.String, expect.String) != 0
}

//go:generate go run ./gen_type.go NullNChar sql.NullString -columndef "nchar(13) null" -compare compareNullChar
var samplesNullNChar = samplesNullChar

//go:generate go run ./gen_type.go NullVarChar sql.NullString -columndef "varchar(13) null" -compare compareNullChar
var samplesNullVarChar = samplesNullChar

//go:generate go run ./gen_type.go NullNVarChar sql.NullString -columndef "nvarchar(13) null" -compare compareNullChar
var samplesNullNVarChar = samplesNullChar

//go:generate go run ./gen_type.go Binary []byte -columndef binary(13) -compare compareBinary
var samplesBinary = [][]byte{
	[]byte("test"),
	[]byte("a longer test"),
}

func compareBinary(recv, expect []byte) bool {
	return !bytes.Equal(bytes.Trim(recv, "\x00"), expect)
}

//go:generate go run ./gen_type.go NullBinary github.com/SAP/go-dblib/asetypes.NullBinary -columndef "binary(13) null" -compare compareNullBinary
var samplesNullBinary = []asetypes.NullBinary{
	{ByteSlice: []byte(""), Valid: false},
	{ByteSlice: []byte(" "), Valid: true},
	{ByteSlice: []byte("test"), Valid: true},
	{Valid: false},
	{ByteSlice: []byte("a longer test"), Valid: true},
}

func compareNullBinary(recv, expect asetypes.NullBinary) bool {
	if recv.Valid != expect.Valid {
		return true
	}

	return !bytes.Equal(recv.ByteSlice, expect.ByteSlice)
}

//go:generate go run ./gen_type.go VarBinary []byte -columndef varbinary(13) -compare compareVarBinary
var samplesVarBinary = samplesBinary

func compareVarBinary(recv, expect []byte) bool {
	return !bytes.Equal(recv, expect)
}

//go:generate go run ./gen_type.go NullVarBinary github.com/SAP/go-dblib/asetypes.NullBinary -columndef "varbinary(13) null" -compare compareNullBinary
var samplesNullVarBinary = samplesNullBinary

//go:generate go run ./gen_type.go Bit bool
// Cannot be nulled
var samplesBit = []bool{true, false}

//go:generate go run ./gen_type.go Image []byte -compare compareBinary
// TODO: -null github.com/SAP/go-dblib/asetypes.NullBinary
var samplesImage = [][]byte{[]byte("test"), []byte("a longer test")}

// TODO: Separate null test, ctlib transforms empty value to null
//go:generate go run ./gen_type.go UniChar string -columndef "unichar(30) null" -compare compareChar
// TODO: -null database/sql.NullString
var samplesUniChar = []string{"", "not a unicode example"}

// TODO: Separate null test, ctlib transforms empty value to null
//go:generate go run ./gen_type.go Text string -columndef "text null" -compare compareChar
// TODO: -null database/sql.NullString
var samplesText = []string{"", "a long text"}

//go:generate go run ./gen_type.go UniText string -columndef unitext -compare compareChar
// TODO: -null database/sql.NullString
var samplesUniText = []string{"not a unicode example", "another not unicode example"}
