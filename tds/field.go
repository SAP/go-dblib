// SPDX-FileCopyrightText: 2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// The types in this file are format fields for TDS_PARAMFMT,
// TDS_PARAMFMT2, TDS_ROWFMT, TDS_ROWFMT2 and data fields for
// TDS_PARAMS and TDS_ROW.
//
// To make handling and identification of the different data types
// easier in Go each data type has their own field structure, created by
// embedding abstract types.
//
// Note: The data types are embedded instead of aliased to provide the
// methods of the embedded types - aliasing does not provide access to
// the methods of the source type.
//
// File-Structure:
// - Format and data status
// - Interfaces
// - Base structs and methods: fieldFmtBase, fieldDataBase
// - Implementations (Derived structs and methods):
//   - fieldFmtLength, fieldData
//   - fieldFmtLengthScale, fieldData
//   - fieldFmtLengthPrecisionScale, fieldDataPrecisionScale
//   - fieldFmtBlob, fieldDataBlob
//   - fieldFmtTxtPtr, fieldDataTxtPtr
// - Utility-functions

package tds

import (
	"fmt"

	"github.com/SAP/go-dblib/asetypes"
)

// Format and data status

// Both Param- and RowFmtStatus are uints communicated using
// TDS_PARAMFMT* and TDS_ROWFMT*. Depending on the token they have
// slightly different meanings.
//
// fmtStatus is only used for legibility in the code when e.g. checking
// for column status or if null types are allowed. The methods setting
// and getting status convert it from and to fmtStatus.
type fmtStatus uint

const (
	tdsFmtColumnStatus fmtStatus = 0x8
	// TODO: Nulltypes
	// tdsFmtNullAllowed  fmtStatus = 0x20
)

//go:generate stringer -type=ParamFmtStatus

// ParamFmtStatus is the type for bitmask values of a ParamFmt status.
type ParamFmtStatus uint

const (
	TDS_PARAM_NOSTATUS     ParamFmtStatus = 0x0
	TDS_PARAM_RETURN       ParamFmtStatus = 0x1
	TDS_PARAM_COLUMNSTATUS ParamFmtStatus = 0x8
	TDS_PARAM_NULLALLOWED  ParamFmtStatus = 0x20
)

//go:generate stringer -type=RowFmtStatus

// RowFmtStatus is the type for bitmask values of a RowFmt status.
type RowFmtStatus uint

const (
	TDS_ROW_NOSTATUS     RowFmtStatus = 0x0
	TDS_ROW_HIDDEN       RowFmtStatus = 0x1
	TDS_ROW_KEY          RowFmtStatus = 0x2
	TDS_ROW_VERSION      RowFmtStatus = 0x4
	TDS_ROW_COLUMNSTATUS RowFmtStatus = 0x8
	TDS_ROW_UPDATEABLE   RowFmtStatus = 0x10
	TDS_ROW_NULLALLOWED  RowFmtStatus = 0x20
	TDS_ROW_IDENTITY     RowFmtStatus = 0x40
	TDS_ROW_PADCHAR      RowFmtStatus = 0x80
)

// DataStatus is the type for bitmask values of a data status.
type DataStatus uint

const (
	TDS_DATA_NONNULL           DataStatus = 0x0
	TDS_DATA_NULL              DataStatus = 0x1
	TDS_DATA_ZEROLENGTHNONNULL DataStatus = 0x2
	TDS_DATA_RESERVED          DataStatus = 0xfc
)

// Interfaces

// FieldFmt is the interface providing the ReadFrom- and WriteTo-method
// to read from or write to a field by communicating format information
// through the tds-protocol from or to a server.
//
// Additionally, several Setter- and Getter-methods are provided.
type FieldFmt interface {
	// Format information as sent to or received from TDS server
	DataType() asetypes.DataType
	setDataType(asetypes.DataType)
	SetName(string)
	Name() string

	// specific to TDS_ROWFMT2
	SetColumnLabel(string)
	ColumnLabel() string
	SetCatalogue(string)
	Catalogue() string
	SetSchema(string)
	Schema() string
	SetTable(string)
	Table() string

	SetStatus(uint)
	Status() uint

	SetUserType(int32)
	UserType() int32
	SetLocaleInfo(string)
	LocaleInfo() string

	// Interface methods for go-ase

	// Returns true if the data type has a fixed length.
	IsFixedLength() bool
	// The return value of LengthBytes depends on IsFixedLength.
	// If the data type has a fixed length LengthBytes returns the
	// total number of bytes of the data portion (not the entire data
	// field - only the actual data).
	// If the data type has a variable length LengthBytes returns the
	// number of bytes to be read from the data stream for the length in
	// bytes of the data portion.
	LengthBytes() int
	// Length returns the maximum length of the column
	// TODO: is this actually required when sending from client?
	MaxLength() int64
	setMaxLength(int64)

	// ReadFrom reads bytes from the passed channel until either the
	// channel is closed or the package has all required information.
	// The read bytes are usually stored in fieldFmt* structs.
	ReadFrom(BytesChannel) (int, error)
	// WriteTo writes bytes to the passed channel until either the
	// channel is closed or the package has written all required
	// information. The information are usually based on fieldFmt*
	// structs.
	WriteTo(BytesChannel) (int, error)

	FormatByteLength() int
}

// FieldData is the interface providing the ReadFrom- and WriteTo-method
// to read from or write to a field by communicating data through the tds-
// protocol from or to a server.
//
// Additionally, several Setter- and Getter-methods are provided.
type FieldData interface {
	// Format information send by TDS server
	Status() DataStatus

	// Interface methods for go-ase
	Format() FieldFmt
	setFormat(FieldFmt)

	// ReadFrom reads bytes from the passed channel until either the
	// channel is closed or the package has all required information.
	// The read bytes are usually stored in fieldData* structs.
	ReadFrom(BytesChannel) (int, error)
	// WriteTo writes bytes to the passed channel until either the
	// channel is closed or the package has written all required
	// information. The information are usually based on fieldData*
	// structs.
	WriteTo(BytesChannel) (int, error)

	Value() interface{}
	SetValue(interface{})
}

// Base structs and methods

type fieldFmtBase struct {
	dataType asetypes.DataType
	name     string

	// specific to TDS_ROWFMT2
	// wide_row controls if the TDS_ROWFMT2 specific members are filled
	// and written. It is set by TDS_ROWFMT2 when creating a field.
	wide_row    bool
	columnLabel string
	catalogue   string
	schema      string
	table       string

	status     fmtStatus
	userType   int32
	localeInfo string

	// length is the maximum length of the data type
	maxLength int64
}

// DataType implements the tds.FieldFmt interface.
func (field fieldFmtBase) DataType() asetypes.DataType {
	return field.dataType
}

func (field *fieldFmtBase) setDataType(t asetypes.DataType) {
	field.dataType = t
}

// SetName implements the tds.FieldFmt interface.
func (field *fieldFmtBase) SetName(name string) {
	field.name = name
}

// Name implements the tds.FieldFmt interface.
func (field fieldFmtBase) Name() string {
	return field.name
}

// SetColumnLabel implements the tds.FieldFmt interface.
func (field *fieldFmtBase) SetColumnLabel(columnLabel string) {
	field.columnLabel = columnLabel
}

// ColumnLabel implements the tds.FieldFmt interface.
func (field fieldFmtBase) ColumnLabel() string {
	return field.columnLabel
}

// SetCatalogue implements the tds.FieldFmt interface.
func (field *fieldFmtBase) SetCatalogue(catalogue string) {
	field.catalogue = catalogue
}

// Catalogue implements the tds.FieldFmt interface.
func (field fieldFmtBase) Catalogue() string {
	return field.catalogue
}

// SetSchema implements the tds.FieldFmt interface.
func (field *fieldFmtBase) SetSchema(schema string) {
	field.schema = schema
}

// Schema implements the tds.FieldFmt interface.
func (field fieldFmtBase) Schema() string {
	return field.schema
}

// SetTable implements the tds.FieldFmt interface.
func (field *fieldFmtBase) SetTable(table string) {
	field.table = table
}

// Table implements the tds.FieldFmt interface.
func (field fieldFmtBase) Table() string {
	return field.table
}

// SetStatus implements the tds.FieldFmt interface.
func (field *fieldFmtBase) SetStatus(status uint) {
	field.status = fmtStatus(status)
}

// Status implements the tds.FieldFmt interface.
func (field fieldFmtBase) Status() uint {
	return uint(field.status)
}

// SetUserType implements the tds.FieldFmt interface.
func (field *fieldFmtBase) SetUserType(userType int32) {
	field.userType = userType
}

// UserType implements the tds.FieldFmt interface.
func (field fieldFmtBase) UserType() int32 {
	return field.userType
}

// SetLocaleInfo implements the tds.FieldFmt interface.
func (field *fieldFmtBase) SetLocaleInfo(localeInfo string) {
	field.localeInfo = localeInfo
}

// LocaleInfo implements the tds.FieldFmt interface.
func (field fieldFmtBase) LocaleInfo() string {
	return field.localeInfo
}

// IsFixedLength implements the tds.FieldFmt interface.
func (field fieldFmtBase) IsFixedLength() bool {
	return field.DataType().ByteSize() != -1
}

// LengthBytes implements the tds.FieldFmt interface.
func (field fieldFmtBase) LengthBytes() int {
	if field.IsFixedLength() {
		return field.DataType().ByteSize()
	}

	return field.DataType().LengthBytes()
}

func (field fieldFmtBase) setMaxLength(i int64) {
	field.maxLength = i
}

// MaxLength implements the tds.FieldFmt interface.
func (field fieldFmtBase) MaxLength() int64 {
	return field.maxLength
}

func (field *fieldFmtBase) readFromBase(ch BytesChannel) (int, error) {
	if field.IsFixedLength() {
		return 0, nil
	}

	length, err := readLengthBytes(ch, field.LengthBytes())
	if err != nil {
		return 0, ErrNotEnoughBytes
	}
	field.maxLength = int64(length)

	return field.LengthBytes(), nil
}

func (field fieldFmtBase) writeToBase(ch BytesChannel) (int, error) {
	if field.IsFixedLength() {
		return 0, nil
	}

	return field.LengthBytes(), writeLengthBytes(ch, field.LengthBytes(), field.MaxLength())
}

type fieldDataBase struct {
	fmt    FieldFmt
	status DataStatus
	value  interface{}
}

func (field *fieldDataBase) setFormat(f FieldFmt) {
	field.fmt = f
}

// Format implements the tds.FieldData interface.
func (field fieldDataBase) Format() FieldFmt {
	return field.fmt
}

// Status implements the tds.FieldData interface.
func (field fieldDataBase) Status() DataStatus {
	return field.status
}

// Value implements the tds.FieldData interface.
func (field *fieldDataBase) Value() interface{} {
	// TODO set a default?
	return field.value
}

// SetValue implements the tds.FieldData interface.
func (field *fieldDataBase) SetValue(value interface{}) {
	field.value = value
}

func (field *fieldDataBase) readFromStatus(ch BytesChannel) (int, error) {
	if fmtStatus(field.fmt.Status())&tdsFmtColumnStatus != tdsFmtColumnStatus {
		return 0, nil
	}

	status, err := ch.Uint8()
	if err != nil {
		return 0, ErrNotEnoughBytes
	}
	field.status = DataStatus(status)
	return 1, nil
}

func (field fieldDataBase) writeToStatus(ch BytesChannel) (int, error) {
	if fmtStatus(field.fmt.Status())&tdsFmtColumnStatus != tdsFmtColumnStatus {
		return 0, nil
	}

	// TODO depends on wide
	if err := ch.WriteUint8(uint8(field.status)); err != nil {
		return 0, fmt.Errorf("failed to write status: %w", err)
	}
	return 1, nil
}

func (field *fieldDataBase) readFrom(ch BytesChannel) (int, error) {
	n, err := field.readFromStatus(ch)
	if err != nil {
		return n, err
	}

	length := field.fmt.LengthBytes()
	if !field.fmt.IsFixedLength() {
		var err error
		length, err = readLengthBytes(ch, field.fmt.LengthBytes())
		if err != nil {
			return n, fmt.Errorf("failed to read %d bytes of length: %w", field.fmt.LengthBytes(), err)
		}
		n += field.fmt.LengthBytes()
	}

	bs, err := ch.Bytes(length)
	if err != nil {
		return n, fmt.Errorf("failed to read %d bytes: %w", length, err)
	}

	field.value, err = field.fmt.DataType().GoValue(endian, bs)
	if err != nil {
		return n, fmt.Errorf("failed to parse field data: %w", err)
	}

	if len(bs) != length {
		return n, fmt.Errorf("expected to read %d bytes of data, read %d", length, len(bs))
	}

	return n, nil
}

func (field fieldDataBase) writeTo(ch BytesChannel) (int, error) {
	n, err := field.writeToStatus(ch)
	if err != nil {
		return n, err
	}

	bs, err := field.fmt.DataType().Bytes(endian, field.value)
	if err != nil {
		return n, fmt.Errorf("error converting field value to bytes: %w", err)
	}

	if !field.fmt.IsFixedLength() {
		if err := writeLengthBytes(ch, field.fmt.LengthBytes(), int64(len(bs))); err != nil {
			return n, fmt.Errorf("failed to write data length: %w", err)
		}
		n += field.fmt.LengthBytes()
	}

	if err = ch.WriteBytes(bs); err != nil {
		return n, fmt.Errorf("failed to write field data: %w", err)
	}
	n += len(bs)

	return n, nil
}

// Implementations: fieldFmtLength, fieldData

type fieldFmtLength struct {
	fieldFmtBase
}

// FormatByteLength implements the tds.FieldFmt interface.
func (field fieldFmtLength) FormatByteLength() int {
	if field.IsFixedLength() {
		return 0
	}

	return field.LengthBytes()
}

// ReadFrom implements the tds.FieldFmt interface.
func (field *fieldFmtLength) ReadFrom(ch BytesChannel) (int, error) {
	return field.readFromBase(ch)
}

// WriteTo implements the tds.FieldFmt interface.
func (field fieldFmtLength) WriteTo(ch BytesChannel) (int, error) {
	return field.writeToBase(ch)
}

type BitFieldFmt struct{ fieldFmtLength }
type DateTimeFieldFmt struct{ fieldFmtLength }
type DateFieldFmt struct{ fieldFmtLength }
type ShortDateFieldFmt struct{ fieldFmtLength }
type Flt4FieldFmt struct{ fieldFmtLength }
type Flt8FieldFmt struct{ fieldFmtLength }
type Int1FieldFmt struct{ fieldFmtLength }
type Int2FieldFmt struct{ fieldFmtLength }
type Int4FieldFmt struct{ fieldFmtLength }
type Int8FieldFmt struct{ fieldFmtLength }
type IntervalFieldFmt struct{ fieldFmtLength }
type Sint1FieldFmt struct{ fieldFmtLength }
type Uint2FieldFmt struct{ fieldFmtLength }
type Uint4FieldFmt struct{ fieldFmtLength }
type Uint8FieldFmt struct{ fieldFmtLength }
type MoneyFieldFmt struct{ fieldFmtLength }
type ShortMoneyFieldFmt struct{ fieldFmtLength }
type TimeFieldFmt struct{ fieldFmtLength }
type BinaryFieldFmt struct{ fieldFmtLength }
type BoundaryFieldFmt struct{ fieldFmtLength }
type CharFieldFmt struct{ fieldFmtLength }
type DateNFieldFmt struct{ fieldFmtLength }
type DateTimeNFieldFmt struct{ fieldFmtLength }
type FltNFieldFmt struct{ fieldFmtLength }
type IntNFieldFmt struct{ fieldFmtLength }
type UintNFieldFmt struct{ fieldFmtLength }
type LongBinaryFieldFmt struct{ fieldFmtLength }
type LongCharFieldFmt struct{ fieldFmtLength }
type MoneyNFieldFmt struct{ fieldFmtLength }
type SensitivityFieldFmt struct{ fieldFmtLength }
type TimeNFieldFmt struct{ fieldFmtLength }
type VarBinaryFieldFmt struct{ fieldFmtLength }
type VarCharFieldFmt struct{ fieldFmtLength }

type fieldData struct{ fieldDataBase }

// ReadFrom implements the tds.FieldData interface.
func (field *fieldData) ReadFrom(ch BytesChannel) (int, error) {
	return field.readFrom(ch)
}

// WriteTo implements the tds.FieldData interface.
func (field fieldData) WriteTo(ch BytesChannel) (int, error) {
	return field.writeTo(ch)
}

type BitFieldData struct{ fieldData }
type DateTimeFieldData struct{ fieldData }
type DateFieldData struct{ fieldData }
type ShortDateFieldData struct{ fieldData }
type Flt4FieldData struct{ fieldData }
type Flt8FieldData struct{ fieldData }
type Int1FieldData struct{ fieldData }
type Int2FieldData struct{ fieldData }
type Int4FieldData struct{ fieldData }
type Int8FieldData struct{ fieldData }
type IntervalFieldData struct{ fieldData }
type Sint1FieldData struct{ fieldData }
type Uint2FieldData struct{ fieldData }
type Uint4FieldData struct{ fieldData }
type Uint8FieldData struct{ fieldData }
type MoneyFieldData struct{ fieldData }
type ShortMoneyFieldData struct{ fieldData }
type TimeFieldData struct{ fieldData }
type BinaryFieldData struct{ fieldData }
type BoundaryFieldData struct{ fieldData }
type CharFieldData struct{ fieldData }
type DateNFieldData struct{ fieldData }
type DateTimeNFieldData struct{ fieldData }
type FltNFieldData struct{ fieldData }
type IntNFieldData struct{ fieldData }
type UintNFieldData struct{ fieldData }
type LongBinaryFieldData struct{ fieldData }
type LongCharFieldData struct{ fieldData }
type MoneyNFieldData struct{ fieldData }
type SensitivityFieldData struct{ fieldData }
type TimeNFieldData struct{ fieldData }
type VarBinaryFieldData struct{ fieldData }
type VarCharFieldData struct{ fieldData }

// Implementations: fieldFmtLengthScale, fieldData

type fieldFmtBaseScale struct {
	scale uint8
}

// Scale returns the scale stored in a field format.
func (field fieldFmtBaseScale) Scale() uint8 {
	return field.scale
}

func (field *fieldFmtBaseScale) readFromScale(ch BytesChannel) (int, error) {
	var err error
	field.scale, err = ch.Uint8()
	if err != nil {
		return 0, ErrNotEnoughBytes
	}
	return 1, nil
}

func (field fieldFmtBaseScale) writeToScale(ch BytesChannel) (int, error) {
	if err := ch.WriteUint8(field.scale); err != nil {
		return 0, fmt.Errorf("failed to write scale: %w", err)
	}
	return 1, nil
}

type fieldFmtLengthScale struct {
	fieldFmtBase
	fieldFmtBaseScale
}

// FormatByteLength implements the tds.FieldFmt interface.
func (field fieldFmtLengthScale) FormatByteLength() int {
	// 1 byte scale
	// 1 to 4 bytes length
	return 1 + field.LengthBytes()
}

// ReadFrom implements the tds.FieldFmt interface.
func (field *fieldFmtLengthScale) ReadFrom(ch BytesChannel) (int, error) {
	n, err := field.readFromBase(ch)
	if err != nil {
		return n, err
	}

	n2, err := field.readFromScale(ch)
	return n + n2, err
}

// WriteTo implements the tds.FieldFmt interface.
func (field fieldFmtLengthScale) WriteTo(ch BytesChannel) (int, error) {
	n, err := field.writeToBase(ch)
	if err != nil {
		return n, err
	}

	n2, err := field.writeToScale(ch)
	return n + n2, err
}

type BigDateTimeNFieldFmt struct{ fieldFmtLengthScale }
type BigTimeNFieldFmt struct{ fieldFmtLengthScale }

type BigDateTimeNFieldData struct{ fieldData }
type BigTimeNFieldData struct{ fieldData }

// Implementations: fieldFmtLengthPrecisionScale,
//                  fieldDataPrecisionScale

type fieldFmtBasePrecision struct {
	precision uint8
}

// Precision returns the precision stored in a field format.
func (field fieldFmtBasePrecision) Precision() uint8 {
	return field.precision
}

func (field *fieldFmtBasePrecision) readFromPrecision(ch BytesChannel) (int, error) {
	var err error
	field.precision, err = ch.Uint8()
	if err != nil {
		return 0, ErrNotEnoughBytes
	}
	return 1, nil
}

func (field fieldFmtBasePrecision) writeToPrecision(ch BytesChannel) (int, error) {
	if err := ch.WriteUint8(field.precision); err != nil {
		return 0, fmt.Errorf("failed to write precision: %w", err)
	}
	return 1, nil
}

type fieldFmtLengthPrecisionScale struct {
	fieldFmtBase
	fieldFmtBasePrecision
	fieldFmtBaseScale
}

// FormatByteLength implements the tds.FieldFmt interface.
func (field fieldFmtLengthPrecisionScale) FormatByteLength() int {
	return 2 + field.LengthBytes()
}

// ReadFrom implements the tds.FieldFmt interface.
func (field *fieldFmtLengthPrecisionScale) ReadFrom(ch BytesChannel) (int, error) {
	n, err := field.readFromBase(ch)
	if err != nil {
		return n, err
	}

	n2, err := field.readFromPrecision(ch)
	if err != nil {
		return n + n2, err
	}

	n3, err := field.readFromScale(ch)
	return n + n2 + n3, err
}

// WriteTo implements the tds.FieldFmt interface.
func (field fieldFmtLengthPrecisionScale) WriteTo(ch BytesChannel) (int, error) {
	n, err := field.writeToBase(ch)
	if err != nil {
		return n, err
	}

	n2, err := field.writeToPrecision(ch)
	if err != nil {
		return n + n2, err
	}

	n3, err := field.writeToScale(ch)
	return n + n2 + n3, err
}

type DecNFieldFmt struct{ fieldFmtLengthPrecisionScale }
type NumNFieldFmt struct{ fieldFmtLengthPrecisionScale }

type fieldDataPrecisionScale struct {
	fieldData
}

// ReadFrom implements the tds.FieldData interface.
func (field *fieldDataPrecisionScale) ReadFrom(ch BytesChannel) (int, error) {
	n, err := field.readFrom(ch)
	if err != nil {
		return n, err
	}

	dec, ok := field.value.(*asetypes.Decimal)
	if !ok {
		return n, fmt.Errorf("%T is not of type decimal", field.value)
	}

	switch fieldFmt := field.fmt.(type) {
	case *DecNFieldFmt:
		dec.Precision = int(fieldFmt.precision)
		dec.Scale = int(fieldFmt.scale)
	case *NumNFieldFmt:
		dec.Precision = int(fieldFmt.precision)
		dec.Scale = int(fieldFmt.scale)
	default:
		return n, fmt.Errorf("%T is neither of type DecNFieldFmt nor NumNFieldFmt", field.value)
	}

	return n, nil
}

type DecNFieldData struct{ fieldDataPrecisionScale }
type NumNFieldData struct{ fieldDataPrecisionScale }

// Implementations: fieldFmtBlob, fieldDataBlob

//go:generate stringer -type=BlobType

// BlobType is the type for bitmask values of a blob types.
type BlobType uint8

const (
	TDS_BLOB_FULLCLASSNAME BlobType = 0x01
	TDS_BLOB_DBID_CLASSDEF BlobType = 0x02
	TDS_BLOB_CHAR          BlobType = 0x03
	TDS_BLOB_BINARY        BlobType = 0x04
	TDS_BLOB_UNICHAR       BlobType = 0x05
	TDS_LOBLOC_CHAR        BlobType = 0x06
	TDS_LOBLOC_BINARY      BlobType = 0x07
	TDS_LOBLOC_UNICHAR     BlobType = 0x08
)

//go:generate stringer -type=BlobSerializationType

// BlobSerializationType is the type for bitmask values of a blob
// serialization-types.
type BlobSerializationType uint8

const (
	NativeJavaSerialization BlobSerializationType = iota
	NativeCharacterFormat
	BinaryData
	UnicharUTF16
	UnicharUTF8
	UnicharSCSU
)

type fieldFmtBlob struct {
	fieldFmtBase
	blobType BlobType
	classID  string
}

// FormatByteLength implements the tds.FieldFmt interface.
func (field fieldFmtBlob) FormatByteLength() int {
	return 1 + 1 + len(field.classID) + field.LengthBytes()
}

// ReadFrom implements the tds.FieldFmt interface.
func (field *fieldFmtBlob) ReadFrom(ch BytesChannel) (int, error) {
	n, err := field.readFromBase(ch)
	if err != nil {
		return n, err
	}

	blobType, err := ch.Uint8()
	if err != nil {
		return 0, ErrNotEnoughBytes
	}
	field.blobType = BlobType(blobType)
	n++

	if field.blobType == TDS_BLOB_FULLCLASSNAME || field.blobType == TDS_BLOB_DBID_CLASSDEF {
		classIdLength, err := ch.Uint16()
		if err != nil {
			return 0, ErrNotEnoughBytes
		}
		n += 2

		field.classID, err = ch.String(int(classIdLength))
		if err != nil {
			return 0, ErrNotEnoughBytes
		}
		n += int(classIdLength)
	}

	return n, nil
}

// WriteTo implements the tds.FieldFmt interface.
func (field fieldFmtBlob) WriteTo(ch BytesChannel) (int, error) {
	n, err := field.writeToBase(ch)
	if err != nil {
		return n, err
	}

	if err := ch.WriteUint8(uint8(field.blobType)); err != nil {
		return n, fmt.Errorf("failed to write blobtype: %w", err)
	}
	n++

	if field.blobType == TDS_BLOB_FULLCLASSNAME || field.blobType == TDS_BLOB_DBID_CLASSDEF {
		if err := ch.WriteUint16(uint16(len(field.classID))); err != nil {
			return n, fmt.Errorf("failed to write ClassID length: %w", err)
		}
		n += 2

		if len(field.classID) > 0 {
			if err := ch.WriteString(field.classID); err != nil {
				return n, fmt.Errorf("failed to write ClassID: %w", err)
			}
			n += len(field.classID)
		}
	}

	return n, nil
}

type BlobFieldFmt struct{ fieldFmtBlob }

type fieldDataBlob struct {
	fieldData
	serializationType BlobSerializationType
	subClassID        string
	locator           string
}

const fieldDataBlobHighBit uint32 = 0x80000000

// ReadFrom implements the tds.FieldData interface.
func (field *fieldDataBlob) ReadFrom(ch BytesChannel) (int, error) {
	fieldFmt, ok := field.fmt.(*BlobFieldFmt)
	if !ok {
		return 0, fmt.Errorf("field.fmt is not of type BlobFieldfmt")
	}

	n, err := field.readFromStatus(ch)
	if err != nil {
		return n, err
	}

	serialization, err := ch.Uint8()
	if err != nil {
		return 0, ErrNotEnoughBytes
	}
	n++

	switch serialization {
	case 0:
		switch fieldFmt.blobType {
		case TDS_BLOB_FULLCLASSNAME, TDS_BLOB_DBID_CLASSDEF:
			field.serializationType = NativeJavaSerialization
		case TDS_BLOB_CHAR:
			field.serializationType = NativeCharacterFormat
		case TDS_BLOB_BINARY:
			field.serializationType = BinaryData
		case TDS_BLOB_UNICHAR:
			field.serializationType = UnicharUTF16
		}
	case 1:
		if fieldFmt.blobType != TDS_BLOB_UNICHAR {
			return n, fmt.Errorf("invalid blob (%s) and serialization (%d) type combination",
				fieldFmt.blobType, serialization)
		}
		field.serializationType = UnicharUTF8
	case 2:
		if fieldFmt.blobType != TDS_BLOB_UNICHAR {
			return n, fmt.Errorf("invalid blob (%s) and serialization (%d) type combination",
				fieldFmt.blobType, serialization)
		}
		field.serializationType = UnicharSCSU
	default:
		return n, fmt.Errorf("unhandled serialization type %d", serialization)
	}

	switch fieldFmt.blobType {
	case TDS_BLOB_FULLCLASSNAME, TDS_BLOB_DBID_CLASSDEF:
		subClassIdLength, err := ch.Uint16()
		if err != nil {
			return 0, ErrNotEnoughBytes
		}
		n += 2

		if subClassIdLength > 0 {
			field.subClassID, err = ch.String(int(subClassIdLength))
			if err != nil {
				return 0, ErrNotEnoughBytes
			}
			n += int(subClassIdLength)
		}
	case TDS_LOBLOC_CHAR, TDS_LOBLOC_BINARY, TDS_LOBLOC_UNICHAR:
		locatorLength, err := ch.Uint16()
		if err != nil {
			return 0, ErrNotEnoughBytes
		}
		n += 2

		field.locator, err = ch.String(int(locatorLength))
		if err != nil {
			return 0, ErrNotEnoughBytes
		}
		n += int(locatorLength)
	}

	// TODO better data type
	data := []byte{}

	for {
		dataLen, err := ch.Uint32()
		if err != nil {
			return 0, ErrNotEnoughBytes
		}
		n += 4

		// extract high bit:
		// 0 -> last data set
		// 1 -> another data set follows
		highBitSet := dataLen&fieldDataBlobHighBit == fieldDataBlobHighBit
		dataLen = dataLen &^ fieldDataBlobHighBit

		// if high bit is set and dataLen is zero no data array follows,
		// instead read the next data length immediately
		if highBitSet {
			break
		}

		if dataLen == 0 {
			continue
		}

		dataPart, err := ch.Bytes(int(dataLen))
		if err != nil {
			return 0, ErrNotEnoughBytes
		}
		n += int(dataLen)

		// TODO this is inefficient for large datasets - must be
		// replaced by a low-overhead extensible byte storage (so - not
		// a slice)
		data = append(data, dataPart...)
	}

	field.value = data

	return n, nil
}

// WriteTo implements the tds.FieldData interface.
func (field fieldDataBlob) WriteTo(ch BytesChannel) (int, error) {
	fieldFmt, ok := field.fmt.(*BlobFieldFmt)
	if !ok {
		return 0, fmt.Errorf("field.fmt is not of type BlobFieldFmt")
	}

	n, err := field.writeToStatus(ch)
	if err != nil {
		return n, err
	}

	var serialization uint8
	switch field.serializationType {
	case NativeJavaSerialization, NativeCharacterFormat, BinaryData, UnicharUTF16:
		serialization = 0
	case UnicharUTF8:
		serialization = 1
	case UnicharSCSU:
		serialization = 2
	}
	if err := ch.WriteUint8(serialization); err != nil {
		return n, fmt.Errorf("failed to write SerializationType: %w", err)
	}
	n++

	switch fieldFmt.blobType {
	case TDS_BLOB_FULLCLASSNAME, TDS_BLOB_DBID_CLASSDEF:
		if err := ch.WriteUint16(uint16(len(field.subClassID))); err != nil {
			return n, fmt.Errorf("failed to write SubClassID length: %w", err)
		}
		n += 2

		if err := ch.WriteString(field.subClassID); err != nil {
			return n, fmt.Errorf("failed to write SubClassID: %w", err)
		}
		n += len(field.subClassID)
	case TDS_LOBLOC_CHAR, TDS_LOBLOC_BINARY, TDS_LOBLOC_UNICHAR:
		if err := ch.WriteUint16(uint16(len(field.locator))); err != nil {
			return n, fmt.Errorf("failed to write Locator length: %w", err)
		}
		n += 2

		if err := ch.WriteString(field.locator); err != nil {
			return n, fmt.Errorf("failed to write Locator: %w", err)
		}
		n += len(field.locator)
	}

	data, ok := field.value.([]byte)
	if !ok {
		return n, fmt.Errorf("field.value is not of type []byte, but type %T", field.value)
	}

	dataLen := 1024
	if dataLen > len(data) {
		dataLen = len(data)
	}

	start, end := 0, dataLen
	for {
		passLen := uint32(dataLen)
		if end == len(data) {
			passLen |= fieldDataBlobHighBit
		}

		if err := ch.WriteUint32(uint32(passLen)); err != nil {
			return n, fmt.Errorf("failed to write data chunk length: %w", err)
		}
		n += 4

		if err := ch.WriteBytes(data[start:end]); err != nil {
			return n, fmt.Errorf("failed to write %d bytes of data: %w", dataLen, err)
		}
		n += end - start

		if end == len(data) {
			break
		}

		start = end
		end += dataLen
	}

	return n, nil
}

type BlobFieldData struct{ fieldDataBlob }

// Implementations: fieldFmtTxtPtr, fieldDataTxtPtr

type fieldFmtTxtPtr struct {
	fieldFmtBase

	tableName string
}

// FormatByteLength implements the tds.FieldFmt interface.
func (field fieldFmtTxtPtr) FormatByteLength() int {
	return 2 + len(field.tableName) + field.LengthBytes()
}

// ReadFrom implements the tds.FieldFmt interface.
func (field *fieldFmtTxtPtr) ReadFrom(ch BytesChannel) (int, error) {
	n, err := field.readFromBase(ch)
	if err != nil {
		return n, err
	}

	tableNameLength, err := ch.Uint16()
	if err != nil {
		return 0, ErrNotEnoughBytes
	}
	n += 2

	field.tableName, err = ch.String(int(tableNameLength))
	if err != nil {
		return 0, ErrNotEnoughBytes
	}
	n += int(tableNameLength)

	return n, nil
}

// WriteTo implements the tds.FieldFmt interface.
func (field fieldFmtTxtPtr) WriteTo(ch BytesChannel) (int, error) {
	n, err := field.writeToBase(ch)
	if err != nil {
		return n, err
	}

	if err := ch.WriteUint16(uint16(len(field.tableName))); err != nil {
		return n, fmt.Errorf("failed to write TableName length: %w", err)
	}
	n += 2

	if err := ch.WriteString(field.tableName); err != nil {
		return n, fmt.Errorf("failed to write TableName: %w", err)
	}
	n += len(field.tableName)

	return n, nil
}

type ImageFieldFmt struct{ fieldFmtTxtPtr }
type TextFieldFmt struct{ fieldFmtTxtPtr }
type UniTextFieldFmt struct{ fieldFmtTxtPtr }
type XMLFieldFmt struct{ fieldFmtTxtPtr }

type fieldDataTxtPtr struct {
	fieldData

	txtPtr    []byte
	timeStamp []byte
}

// ReadFrom implements the tds.FieldData interface.
func (field *fieldDataTxtPtr) ReadFrom(ch BytesChannel) (int, error) {
	n, err := field.readFromStatus(ch)
	if err != nil {
		return n, err
	}

	txtPtrLen, err := ch.Uint8()
	if err != nil {
		return 0, ErrNotEnoughBytes
	}
	n++

	field.txtPtr, err = ch.Bytes(int(txtPtrLen))
	if err != nil {
		return 0, ErrNotEnoughBytes
	}
	n += int(txtPtrLen)

	field.timeStamp, err = ch.Bytes(8)
	if err != nil {
		return 0, ErrNotEnoughBytes
	}
	n += 8

	dataLen, err := ch.Uint32()
	if err != nil {
		return 0, ErrNotEnoughBytes
	}
	n += 4

	field.value, err = ch.Bytes(int(dataLen))
	if err != nil {
		return 0, ErrNotEnoughBytes
	}
	n += int(dataLen)

	return n, nil
}

// WriteTo implements the tds.FieldData interface.
func (field fieldDataTxtPtr) WriteTo(ch BytesChannel) (int, error) {
	n, err := field.writeToStatus(ch)
	if err != nil {
		return n, err
	}

	if err := ch.WriteUint8(uint8(len(field.txtPtr))); err != nil {
		return n, fmt.Errorf("failed to write TxtPtr length: %w", err)
	}
	n++

	if err := ch.WriteBytes(field.txtPtr); err != nil {
		return n, fmt.Errorf("failed to write TxtPtr: %w", err)
	}
	n += len(field.txtPtr)

	if err := ch.WriteBytes(field.timeStamp); err != nil {
		return n, fmt.Errorf("failed to write TimeStamp: %w", err)
	}
	n += len(field.timeStamp)

	var data []byte
	switch t := field.value.(type) {
	case string:
		data = []byte(t)
	case []byte:
		data = t
	default:
		return n, fmt.Errorf("field value is of type %T instead of string or byte slice", field.value)
	}

	if err := ch.WriteUint32(uint32(len(data))); err != nil {
		return n, fmt.Errorf("failed to write Data length: %w", err)
	}
	n += 4

	if err := ch.WriteBytes(data); err != nil {
		return n, fmt.Errorf("failed to write Data: %w", err)
	}
	n += len(data)

	return n, nil
}

type ImageFieldData struct{ fieldDataTxtPtr }
type TextFieldData struct{ fieldDataTxtPtr }
type UniTextFieldData struct{ fieldDataTxtPtr }
type XMLFieldData struct{ fieldDataTxtPtr }

// Utility-functions

func readLengthBytes(ch BytesChannel, n int) (int, error) {
	var length int
	var err error
	switch n {
	case 4:
		var tmp uint32
		tmp, err = ch.Uint32()
		length = int(tmp)
	case 2:
		var tmp uint16
		tmp, err = ch.Uint16()
		length = int(tmp)
	default:
		var tmp uint8
		tmp, err = ch.Uint8()
		length = int(tmp)
	}

	if err != nil {
		return 0, ErrNotEnoughBytes
	}

	return length, nil
}

func writeLengthBytes(ch BytesChannel, byteCount int, n int64) error {
	var err error
	switch byteCount {
	case 4:
		err = ch.WriteUint32(uint32(n))
	case 2:
		err = ch.WriteUint16(uint16(n))
	default:
		err = ch.WriteUint8(uint8(n))
	}

	if err != nil {
		return fmt.Errorf("failed to write length: %w", err)
	}

	return nil
}

// LookupFieldFmt returns the FieldFmt for a given data type and sets
// required values in it.
func LookupFieldFmt(dataType asetypes.DataType) (FieldFmt, error) {
	var f FieldFmt
	switch dataType {
	case asetypes.BIGDATETIMEN:
		f = &BigDateTimeNFieldFmt{}
	case asetypes.BIGTIMEN:
		f = &BigTimeNFieldFmt{}
	case asetypes.BIT:
		f = &BitFieldFmt{}
	case asetypes.DATETIME:
		f = &DateTimeFieldFmt{}
	case asetypes.DATE:
		f = &DateFieldFmt{}
	case asetypes.SHORTDATE:
		f = &ShortDateFieldFmt{}
	case asetypes.FLT4:
		f = &Flt4FieldFmt{}
	case asetypes.FLT8:
		f = &Flt8FieldFmt{}
	case asetypes.INT1:
		f = &Int1FieldFmt{}
	case asetypes.INT2:
		f = &Int2FieldFmt{}
	case asetypes.INT4:
		f = &Int4FieldFmt{}
	case asetypes.INT8:
		f = &Int8FieldFmt{}
	case asetypes.INTERVAL:
		f = &IntervalFieldFmt{}
	case asetypes.SINT1:
		f = &Sint1FieldFmt{}
	case asetypes.UINT2:
		f = &Uint2FieldFmt{}
	case asetypes.UINT4:
		f = &Uint4FieldFmt{}
	case asetypes.UINT8:
		f = &Uint8FieldFmt{}
	case asetypes.MONEY:
		f = &MoneyFieldFmt{}
	case asetypes.SHORTMONEY:
		f = &ShortMoneyFieldFmt{}
	case asetypes.TIME:
		f = &TimeFieldFmt{}
	case asetypes.BINARY:
		f = &BinaryFieldFmt{}
	case asetypes.BOUNDARY:
		f = &BoundaryFieldFmt{}
	case asetypes.CHAR:
		f = &CharFieldFmt{}
	case asetypes.DATEN:
		f = &DateNFieldFmt{}
	case asetypes.DATETIMEN:
		f = &DateTimeNFieldFmt{}
	case asetypes.FLTN:
		f = &FltNFieldFmt{}
	case asetypes.INTN:
		f = &IntNFieldFmt{}
	case asetypes.UINTN:
		f = &UintNFieldFmt{}
	case asetypes.LONGBINARY:
		f = &LongBinaryFieldFmt{}
		f.setMaxLength(2147483647)
	case asetypes.LONGCHAR:
		f = &LongCharFieldFmt{}
	case asetypes.MONEYN:
		f = &MoneyNFieldFmt{}
	case asetypes.SENSITIVITY:
		f = &SensitivityFieldFmt{}
	case asetypes.TIMEN:
		f = &TimeNFieldFmt{}
	case asetypes.VARBINARY:
		f = &VarBinaryFieldFmt{}
	case asetypes.VARCHAR:
		f = &VarCharFieldFmt{}
		f.setMaxLength(255)
	case asetypes.DECN:
		f = &DecNFieldFmt{}
	case asetypes.NUMN:
		f = &NumNFieldFmt{}
	case asetypes.BLOB:
		f = &BlobFieldFmt{}
	case asetypes.IMAGE:
		f = &ImageFieldFmt{}
	case asetypes.TEXT:
		f = &TextFieldFmt{}
	case asetypes.UNITEXT:
		f = &UniTextFieldFmt{}
	case asetypes.XML:
		f = &XMLFieldFmt{}
	default:
		return nil, fmt.Errorf("unhandled datatype '%s'", dataType)
	}

	f.setDataType(dataType)
	return f, nil
}

/// LookupFieldData returns the FieldData for a given field format.
func LookupFieldData(fieldFmt FieldFmt) (FieldData, error) {
	var d FieldData

	switch fieldFmt.DataType() {
	case asetypes.BIGDATETIMEN:
		d = &BigDateTimeNFieldData{}
	case asetypes.BIGTIMEN:
		d = &BigTimeNFieldData{}
	case asetypes.BIT:
		d = &BitFieldData{}
	case asetypes.DATETIME:
		d = &DateTimeFieldData{}
	case asetypes.DATE:
		d = &DateFieldData{}
	case asetypes.SHORTDATE:
		d = &ShortDateFieldData{}
	case asetypes.FLT4:
		d = &Flt4FieldData{}
	case asetypes.FLT8:
		d = &Flt8FieldData{}
	case asetypes.INT1:
		d = &Int1FieldData{}
	case asetypes.INT2:
		d = &Int2FieldData{}
	case asetypes.INT4:
		d = &Int4FieldData{}
	case asetypes.INT8:
		d = &Int8FieldData{}
	case asetypes.INTERVAL:
		d = &IntervalFieldData{}
	case asetypes.SINT1:
		d = &Sint1FieldData{}
	case asetypes.UINT2:
		d = &Uint2FieldData{}
	case asetypes.UINT4:
		d = &Uint4FieldData{}
	case asetypes.UINT8:
		d = &Uint8FieldData{}
	case asetypes.MONEY:
		d = &MoneyFieldData{}
	case asetypes.SHORTMONEY:
		d = &ShortMoneyFieldData{}
	case asetypes.TIME:
		d = &TimeFieldData{}
	case asetypes.BINARY:
		d = &BinaryFieldData{}
	case asetypes.BOUNDARY:
		d = &BoundaryFieldData{}
	case asetypes.CHAR:
		d = &CharFieldData{}
	case asetypes.DATEN:
		d = &DateNFieldData{}
	case asetypes.DATETIMEN:
		d = &DateTimeNFieldData{}
	case asetypes.FLTN:
		d = &FltNFieldData{}
	case asetypes.INTN:
		d = &IntNFieldData{}
	case asetypes.UINTN:
		d = &UintNFieldData{}
	case asetypes.LONGBINARY:
		d = &LongBinaryFieldData{}
	case asetypes.LONGCHAR:
		d = &LongCharFieldData{}
	case asetypes.MONEYN:
		d = &MoneyNFieldData{}
	case asetypes.SENSITIVITY:
		d = &SensitivityFieldData{}
	case asetypes.TIMEN:
		d = &TimeNFieldData{}
	case asetypes.VARBINARY:
		d = &VarBinaryFieldData{}
	case asetypes.VARCHAR:
		d = &VarCharFieldData{}
	case asetypes.DECN:
		d = &DecNFieldData{}
	case asetypes.NUMN:
		d = &NumNFieldData{}
	case asetypes.BLOB:
		d = &BlobFieldData{}
	case asetypes.IMAGE:
		d = &ImageFieldData{}
	case asetypes.TEXT:
		d = &TextFieldData{}
	case asetypes.UNITEXT:
		d = &UniTextFieldData{}
	case asetypes.XML:
		d = &XMLFieldData{}
	default:
		return nil, fmt.Errorf("unhandled datatype: '%s'", fieldFmt.DataType())
	}

	d.setFormat(fieldFmt)
	return d, nil
}

// LookupFieldFmtData returns both Fieldfmt and FieldData for a given
// data type.
func LookupFieldFmtData(dataType asetypes.DataType) (FieldFmt, FieldData, error) {
	fieldFmt, err := LookupFieldFmt(dataType)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find field format: %w", err)
	}

	data, err := LookupFieldData(fieldFmt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find field data: %w", err)
	}

	return fieldFmt, data, nil
}
