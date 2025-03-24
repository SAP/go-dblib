// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import "fmt"

// HeaderOnlyPackage is used to communicate header-only packets using
// the same communication channels as regular token-based packages in
// go-ase.
type HeaderOnlyPackage struct {
	Header PacketHeader
}

// ReadFrom implements the tds.Package interface.
func (pkg HeaderOnlyPackage) ReadFrom(ch BytesChannel) error {
	return fmt.Errorf("HeaderOnlyPackages cannot be read from a ByteChannel")
}

// WriteTo implements the tds.Package interface.
func (pkg HeaderOnlyPackage) WriteTo(ch BytesChannel) error {
	return fmt.Errorf("HeaderOnlyPackages cannot be written to a ByteChannel")
}

func (pkg HeaderOnlyPackage) String() string {
	return fmt.Sprintf("Header: %s", pkg.Header)
}
