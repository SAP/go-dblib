// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import (
	"encoding/binary"
	"io"
)

// endian is a driver-level configuration.
var endian binary.ByteOrder = binary.LittleEndian

// writeBasedOnEndian writes either little or big based on the set
// endianness.
func writeBasedOnEndian(stream io.Writer, little byte, big byte) (int, error) {
	if endian == binary.LittleEndian {
		return stream.Write([]byte{little})
	} else {
		return stream.Write([]byte{big})
	}
}
