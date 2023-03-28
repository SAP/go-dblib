// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
// SPDX-FileCopyrightText: 2023 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package dblib

import (
	"database/sql"
	"fmt"
)

// ASEIsolationLevel reflects the ASE isolation levels.
type ASEIsolationLevel int

// Valid ASE isolation levels.
const (
	ASELevelInvalid         ASEIsolationLevel = -1
	ASELevelReadUncommitted ASEIsolationLevel = iota
	ASELevelReadCommitted
	ASELevelRepeatableRead
	ASELevelSerializableRead
)

var (
	// sql2ase maps sql.IsolationLevel to dblib.ASEIsolationLevel.
	sql2ase = map[sql.IsolationLevel]ASEIsolationLevel{
		sql.LevelDefault:         ASELevelReadCommitted,
		sql.LevelReadUncommitted: ASELevelReadUncommitted,
		sql.LevelReadCommitted:   ASELevelReadCommitted,
		sql.LevelWriteCommitted:  ASELevelInvalid,
		sql.LevelRepeatableRead:  ASELevelRepeatableRead,
		sql.LevelSerializable:    ASELevelSerializableRead,
		sql.LevelLinearizable:    ASELevelInvalid,
	}
)

// ASEIsolationLevelFromGo take a database/sql.IsolationLevel and returns
// the relevant isolation level for ASE.
func ASEIsolationLevelFromGo(lvl sql.IsolationLevel) (ASEIsolationLevel, error) {
	aseLvl, ok := sql2ase[lvl]
	if !ok {
		return ASELevelInvalid, fmt.Errorf("Unknown database/sql.IsolationLevel: %v", lvl)
	}

	if aseLvl == ASELevelInvalid {
		return ASELevelInvalid, fmt.Errorf("Isolation level %v is not supported by ASE", lvl)
	}

	return aseLvl, nil
}

// ToGo returns the database/sql.IsolationLevel equivalent of the ASE
// isolation level.
func (lvl ASEIsolationLevel) ToGo() sql.IsolationLevel {
	for sqlLvl, aseLvl := range sql2ase {
		if aseLvl == lvl {
			return sqlLvl
		}
	}

	return sql.LevelDefault
}

// String implements the Stringer interface.
func (lvl ASEIsolationLevel) String() string {
	return lvl.ToGo().String()
}
