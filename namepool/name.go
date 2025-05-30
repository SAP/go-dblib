// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package namepool

// Name is a member of a Pool. It contains the formatted string of the
// format and the ID as well as the ID itself.
type Name struct {
	name string
	id   *uint64

	pool *pool
}

// Name returns the name of name.
func (name Name) Name() string {
	return name.name
}

func (name Name) String() string {
	return name.Name()
}

// ID returns the id of name.
func (name Name) ID() uint64 {
	return *name.id
}

// Release calls to the Names' Pool to release itself. The
// restrictions and affects of Pool.Release apply.
func (name *Name) Release() {
	name.pool.Release(name)
}
