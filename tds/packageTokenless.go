// SPDX-FileCopyrightText: 2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import (
	"bytes"
	"fmt"
)

// TokenlessPackage contains data.
type TokenlessPackage struct {
	Data *bytes.Buffer
}

// NewTokenLessPackage creates a tokenless-package.
func NewTokenlessPackage() *TokenlessPackage {
	return &TokenlessPackage{
		Data: &bytes.Buffer{},
	}
}

// ReadFrom implements the tds.Package interface.
func (pkg *TokenlessPackage) ReadFrom(ch BytesChannel) error {
	_, err := pkg.Data.ReadFrom(ch)
	return err
}

// WriteTo implements the tds.Package interface.
func (pkg TokenlessPackage) WriteTo(ch BytesChannel) error {
	return ch.WriteBytes(pkg.Data.Bytes())
}

func (pkg TokenlessPackage) String() string {
	return fmt.Sprintf("%T(possibleToken=%x) %#v", pkg, pkg.Data.Bytes()[0], pkg)
}
