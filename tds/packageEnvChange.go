// SPDX-FileCopyrightText: 2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import "fmt"

//go:generate stringer -type=EnvChangeType

// EnvChangeType signals which part of the environment was updated.
type EnvChangeType uint8

const (
	TDS_ENV_DB EnvChangeType = iota + 1
	TDS_ENV_LANG
	TDS_ENV_CHARSET
	TDS_ENV_PACKSIZE
)

// EnvChangePackage is used to communicate multiple environment changes.
type EnvChangePackage struct {
	members []EnvChangePackageField
}

// ReadFrom implements the tds.Package interface.
func (pkg *EnvChangePackage) ReadFrom(ch BytesChannel) error {
	length, err := ch.Uint16()
	if err != nil {
		return ErrNotEnoughBytes
	}

	var n uint16 = 0
	for n < length {
		member := EnvChangePackageField{}
		i, err := member.ReadFrom(ch)
		if err != nil {
			return fmt.Errorf("error reading EnvChangePackage member: %w", err)
		}
		n += uint16(i)

		pkg.members = append(pkg.members, member)
	}

	if n > length {
		return fmt.Errorf("read too many bytes, %d instead of expected %d", n, length)
	}

	return nil
}

// WriteTo implements the tds.Package interface.
func (pkg EnvChangePackage) WriteTo(ch BytesChannel) error {
	if err := ch.WriteUint8(byte(TDS_ENVCHANGE)); err != nil {
		return fmt.Errorf("error writing TDS token %s: %w", TDS_ENVCHANGE, err)
	}

	totalLength := 0
	for _, member := range pkg.members {
		totalLength += member.ByteLength()
	}

	if err := ch.WriteUint16(uint16(totalLength)); err != nil {
		return fmt.Errorf("error writing length: %w", err)
	}

	length := 0
	for _, member := range pkg.members {
		n, err := member.WriteTo(ch)
		if err != nil {
			return fmt.Errorf("error writing EnvChangePackage member: %w", err)
		}
		length += n
	}

	if length != totalLength {
		return fmt.Errorf("wrote %d bytes instead of expected %d bytes", length, totalLength)
	}

	return nil
}

func (pkg EnvChangePackage) String() string {
	s := fmt.Sprintf("%T(", pkg)

	for _, member := range pkg.members {
		s += fmt.Sprintf("%s(%s -> %s)", member.Type, member.OldValue, member.NewValue)
	}

	return s + ")"
}

// EnvChangePackageField is a single environment change.
type EnvChangePackageField struct {
	Type               EnvChangeType
	NewValue, OldValue string
}

// ReadFrom implements the tds.Package interface.
func (field *EnvChangePackageField) ReadFrom(ch BytesChannel) (int, error) {
	// n is the amount of bytes read from channel
	n := 0

	typ, err := ch.Uint8()
	if err != nil {
		return n, ErrNotEnoughBytes
	}
	field.Type = EnvChangeType(typ)
	n++

	length, err := ch.Uint8()
	if err != nil {
		return n, ErrNotEnoughBytes
	}
	n++

	if length > 0 {
		field.NewValue, err = ch.String(int(length))
		if err != nil {
			return n, ErrNotEnoughBytes
		}
		n += int(length)
	}

	length, err = ch.Uint8()
	if err != nil {
		return n, ErrNotEnoughBytes
	}
	n++

	if length > 0 {
		field.OldValue, err = ch.String(int(length))
		if err != nil {
			return n, ErrNotEnoughBytes
		}
		n += int(length)
	}

	return n, nil
}

// WriteTo implements the tds.Package interface.
func (field EnvChangePackageField) WriteTo(ch BytesChannel) (int, error) {
	if err := ch.WriteUint8(uint8(field.Type)); err != nil {
		return 0, fmt.Errorf("error writing type: %w", err)
	}
	n := 1

	if err := ch.WriteUint8(uint8(len(field.NewValue))); err != nil {
		return n, fmt.Errorf("error writing new value length: %w", err)
	}
	n++

	if err := ch.WriteString(field.NewValue); err != nil {
		return n, fmt.Errorf("error writing new value: %w", err)
	}
	n += len(field.NewValue)

	if err := ch.WriteUint8(uint8(len(field.OldValue))); err != nil {
		return n, fmt.Errorf("error writing old value length: %w", err)
	}
	n++

	if err := ch.WriteString(field.OldValue); err != nil {
		return n, fmt.Errorf("error writing old value: %w", err)
	}
	n += len(field.OldValue)

	return n, nil
}

// ByteLength returns the length in bytes.
func (field EnvChangePackageField) ByteLength() int {
	// type byte
	// + new value length byte + new value length
	// + old value length byte + old value length
	return 3 + len(field.NewValue) + len(field.OldValue)
}
