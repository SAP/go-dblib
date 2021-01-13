// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import "fmt"

var _ Package = (*ReturnStatusPackage)(nil)

// ReturnStatusPackage communicates the return status of stored
// procedures.
type ReturnStatusPackage struct {
	ReturnValue int32
}

// ReadFrom implements the tds.Package interface.
func (pkg *ReturnStatusPackage) ReadFrom(ch BytesChannel) error {
	var err error

	if pkg.ReturnValue, err = ch.Int32(); err != nil {
		return ErrNotEnoughBytes
	}

	return nil
}

// WriteTo implements the tds.Package interface.
func (pkg ReturnStatusPackage) WriteTo(ch BytesChannel) error {
	return ch.WriteInt32(pkg.ReturnValue)
}

func (pkg ReturnStatusPackage) String() string {
	return fmt.Sprintf("%T(%d)", pkg, pkg.ReturnValue)
}
