// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import "fmt"

//go:generate stringer -type=LanguageStatus

// LanguageStatus is a bitmask for options on a language package.
type LanguageStatus int

const (
	TDS_LANGUAGE_NOARGS   LanguageStatus = 0x0
	TDS_LANGUAGE_HASARGS  LanguageStatus = 0x1
	TDS_LANG_BATCH_PARAMS LanguageStatus = 0x4
)

// LanguagePackage is used to execute an SQL statement without
// arguments.
type LanguagePackage struct {
	Status LanguageStatus
	Cmd    string
}

// ReadFrom implements the tds.Package interface.
func (pkg *LanguagePackage) ReadFrom(ch BytesChannel) error {
	totalLength, err := ch.Uint32()
	if err != nil {
		return ErrNotEnoughBytes
	}

	status, err := ch.Byte()
	if err != nil {
		return ErrNotEnoughBytes
	}
	pkg.Status = LanguageStatus(status)

	pkg.Cmd, err = ch.String(int(totalLength) - 1)
	if err != nil {
		return ErrNotEnoughBytes
	}

	return nil
}

// WriteTo implements the tds.Package interface.
func (pkg *LanguagePackage) WriteTo(ch BytesChannel) error {
	if err := ch.WriteByte(byte(TDS_LANGUAGE)); err != nil {
		return fmt.Errorf("failed to write TDS token %s: %w", TDS_LANGUAGE, err)
	}

	length := 1 + len(pkg.Cmd)
	if err := ch.WriteUint32(uint32(length)); err != nil {
		return fmt.Errorf("failed to write length: %w", err)
	}

	if err := ch.WriteByte(byte(pkg.Status)); err != nil {
		return fmt.Errorf("failed to write status: %w", err)
	}

	if err := ch.WriteString(pkg.Cmd); err != nil {
		return fmt.Errorf("failed to write language command: %w", err)
	}

	return nil
}

func (pkg LanguagePackage) String() string {
	return fmt.Sprintf("%T(%s): %s", pkg, pkg.Status, pkg.Cmd)
}
