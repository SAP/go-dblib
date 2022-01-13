// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import (
	"fmt"
)

var _ Package = (*ParamsPackage)(nil)
var _ Package = (*RowPackage)(nil)

// ParamsPackage is used to communicate a parameters or arguments.
type ParamsPackage struct {
	paramFmt   *ParamFmtPackage
	rowFmt     *RowFmtPackage
	DataFields []FieldData
}

// RowPackage is used to communicate a row.
type RowPackage struct {
	ParamsPackage
}

// NewParamsPackage returns an initialized ParamsPkg.
func NewParamsPackage(data ...FieldData) *ParamsPackage {
	return &ParamsPackage{
		DataFields: data,
	}
}

// LastPkg implements the tds.LastPkgAcceptor interface.
func (pkg *ParamsPackage) LastPkg(other Package) error {
	switch otherPkg := other.(type) {
	case *ParamFmtPackage:
		pkg.paramFmt = otherPkg
	case *RowFmtPackage:
		pkg.rowFmt = otherPkg
	case *ParamsPackage:
		pkg.paramFmt = otherPkg.paramFmt
	case *RowPackage:
		pkg.rowFmt = otherPkg.rowFmt
	case *OrderByPackage:
		pkg.rowFmt = otherPkg.rowFmt
	case *OrderBy2Package:
		pkg.rowFmt = otherPkg.rowFmt
	default:
		return fmt.Errorf("TDS_PARAMS or TDS_ROW received without preceding TDS_PARAMFMT/2 or TDS_ROWFMT")
	}

	if pkg.DataFields != nil {
		// pkg.Datafields has already been filled - this package was
		// created by the client and is being added to the message.
		return nil
	}

	var fieldFmts []FieldFmt
	if pkg.paramFmt != nil {
		fieldFmts = pkg.paramFmt.Fmts
	} else if pkg.rowFmt != nil {
		fieldFmts = pkg.rowFmt.Fmts
	} else {
		return fmt.Errorf("both paramFmt and rowFmt are nil")
	}

	pkg.DataFields = make([]FieldData, len(fieldFmts))

	// Make copies of the formats to store data in
	var err error
	for i, field := range fieldFmts {
		pkg.DataFields[i], err = LookupFieldData(field)
		if err != nil {
			return fmt.Errorf("error copying field: %w", err)
		}
	}

	return nil
}

// ReadFrom implements the tds.Package interface.
func (pkg *ParamsPackage) ReadFrom(ch BytesChannel) error {
	for i, field := range pkg.DataFields {
		// TODO can the written byte count be validated?
		if _, err := field.ReadFrom(ch); err != nil {
			return fmt.Errorf("error occurred reading param field %d data (%s): %w",
				i, field.Format().DataType(), err)
		}
	}

	return nil
}

// WriteTo implements the tds.Package interface.
func (pkg ParamsPackage) WriteTo(ch BytesChannel) error {
	var token Token
	if pkg.paramFmt != nil {
		token = TDS_PARAMS
	} else if pkg.rowFmt != nil {
		token = TDS_ROW
	} else {
		return fmt.Errorf("both paramFmt and rowFmt are nil")
	}

	if err := ch.WriteByte(byte(token)); err != nil {
		return fmt.Errorf("error occurred writing TDS token %s: %w", token, err)
	}

	for i, field := range pkg.DataFields {
		if _, err := field.WriteTo(ch); err != nil {
			return fmt.Errorf("error occurred writing param field %d data: %w", i, err)
		}
	}
	return nil
}

func (pkg ParamsPackage) String() string {
	s := make([]string, len(pkg.DataFields))
	for i, field := range pkg.DataFields {
		s[i] = fmt.Sprintf("%v", field.Value())
	}
	return fmt.Sprintf("%T(%d): %s", pkg, len(pkg.DataFields), s)
}
