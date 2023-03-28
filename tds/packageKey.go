// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
// SPDX-FileCopyrightText: 2023 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import (
	"fmt"

	"github.com/SAP/go-dblib/asetypes"
)

// KeyPackage is used to communicate the row position of a cursor.
type KeyPackage struct {
	DataType asetypes.DataType
	Value    interface{}
}

// ReadFrom implements the tds.Package interface.
func (pkg *KeyPackage) ReadFrom(ch BytesChannel) error {
	var length int

	if pkg.DataType.ByteSize() > 0 {
		// fixed-length data type
		length = pkg.DataType.ByteSize()
	} else {
		// variable-length data type
		length8, err := ch.Uint8()
		if err != nil {
			return ErrNotEnoughBytes
		}
		length = int(length8)
	}

	bs, err := ch.Bytes(length)
	if err != nil {
		return ErrNotEnoughBytes
	}

	pkg.Value, err = pkg.DataType.GoValue(endian, bs)
	if err != nil {
		return fmt.Errorf("tds: error converting bytes to %s: %w", pkg.DataType, err)
	}

	return nil
}

// WriteTo implements the tds.Package interface.
func (pkg KeyPackage) WriteTo(ch BytesChannel) error {
	if err := ch.WriteByte(byte(TDS_KEY)); err != nil {
		return err
	}

	bs, err := pkg.DataType.Bytes(endian, pkg.Value, int64(pkg.DataType.LengthBytes()))
	if err != nil {
		return fmt.Errorf("tds: error converting value %s to bytes: %w", pkg.DataType, err)
	}

	if pkg.DataType.ByteSize() == -1 {
		// variable-length data type
		if err := ch.WriteUint8(uint8(len(bs))); err != nil {
			return err
		}
	}

	if err := ch.WriteBytes(bs); err != nil {
		return err
	}

	return nil
}

func (pkg KeyPackage) String() string {
	return fmt.Sprintf("Key(%s: %v)", pkg.DataType, pkg.Value)
}
