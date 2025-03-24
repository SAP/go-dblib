// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"database/sql"
	"database/sql/driver"

	"github.com/SAP/go-dblib/dsn"
)

// genSQLDBFn is the signature of functions stored in the genSQLDBMap.
type genSQLDBFn func() (*sql.DB, error)

// genSQLDBMap maps abstract names to functions, which are expected to
// return unique sql.DBs.
type genSQLDBMap map[string]genSQLDBFn

var sqlDBMap = make(genSQLDBMap)

// ConnectorCreator is the interface for function expected by InitDBs to
// initialize driver.Connectors.
type ConnectorCreator func(interface{}) (driver.Connector, error)

// RegisterDSN registers at least one new genSQLDBFn in genSQLDBMap
// based on sql.Open.
// If connectorFn is non-nil a second genSQLDBFn is stored with the
// suffix `connector`.
func RegisterDSN(name string, info interface{}, connectorFn ConnectorCreator) error {
	sqlDBMap[name] = func() (*sql.DB, error) {
		db, err := sql.Open("ase", dsn.FormatSimple(info))
		if err != nil {
			return nil, err
		}
		return db, nil
	}

	if connectorFn != nil {
		sqlDBMap[name+" connector"] = func() (*sql.DB, error) {
			connector, err := connectorFn(info)
			if err != nil {
				return nil, err
			}

			return sql.OpenDB(connector), nil
		}
	}

	return nil
}
