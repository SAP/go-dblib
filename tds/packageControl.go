// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

var _ Package = (*ControlPackage)(nil)

type ControlPackage struct {
}

func (pkg *ControlPackage) ReadFrom(ch BytesChannel) error {
	return nil
}

func (pkg ControlPackage) WriteTo(ch BytesChannel) error {
	return nil
}

func (pkg ControlPackage) String() string {
	return ""
}
