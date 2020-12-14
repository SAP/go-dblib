<!--
SPDX-FileCopyrightText: 2020 SAP SE

SPDX-License-Identifier: Apache-2.0
-->

# go-dblib

[![PkgGoDev](https://pkg.go.dev/badge/github.com/SAP/go-dblib)](https://pkg.go.dev/github.com/SAP/go-dblib)
[![Go Report Card](https://goreportcard.com/badge/github.com/SAP/go-dblib)](https://goreportcard.com/report/github.com/SAP/go-dblib)
[![REUSE
status](https://api.reuse.software/badge/github.com/SAP/go-dblib)](https://api.reuse.software/info/github.com/SAP/go-dblib)
![Actions: CI](https://github.com/SAP/go-dblib/workflows/CI/badge.svg)

## Description

`go-dblib` is a shared library for [`go-ase`][purego] and
[`cgo-ase`][cgo]. These are driver-implementations for the
[`database/sql`][pkg-database-sql] package of [Go (golang)][go] to
provide access to SAP ASE instances.

SAP ASE is the shorthand for [SAP Adaptive Server Enterprise][sap-ase],
a relational model database server originally known as Sybase SQL
Server.

## Requirements

The package `go-dblib` is a shared library for the
driver-implementations of [`go-ase`][purego] and [`cgo-ase`][cgo]. Thus, one of
these implementations is required.

## Download

The packages in this repo can be `go get` and imported as usual.

```sh
go get github.com/SAP/go-dblib/<package>
```

## Usage (Example)

By importing `go-dblib` there are several use-cases. For example, the
`dsn`-package can be used to set up DSN information that is required to
connect to the ASE-database by using a connector.

```go
package main

import (
    "database/sql"

    "github.com/SAP/go-dblib/dsn"
    "github.com/SAP/go-ase"
)

func main() {
    d := dsn.NewInfo()
    d.Host = "hostname"
    d.Port = "4901"
    d.Username = "user"
    d.Password = "pass"

    connector, err := ase.NewConnector(*d)
    if err != nil {
        log.Printf("Failed to create connector: %v", err)
        return
    }

    db, err := sql.OpenDB(connector)
    if err != nil {
        log.Printf("Failed to open database: %v", err)
        return
    }
    defer db.Close()

    if err := db.Ping(); if err != nil {
        log.Printf("Failed to ping ASE: %v", err)
    }
}
```

## Unit tests

Unit tests for the packages are included in their respective directories
and can be run using `go test`.

## Known Issues

The list of known issues is available [here][issues].

## How to obtain support

Feel free to open issues for feature requests, bugs, or general feedback [here][issues].

## Contributing

Any help to improve this library is highly appreciated.

For details on how to contribute please see the [contribution](CONTRIBUTING.md) file.

## License

Copyright (c) 2019-2020 SAP SE or an SAP affiliate company. All rights reserved.
This file is licensed under the Apache License 2.0 except as noted otherwise in the [LICENSE file](LICENSES).

[cgo]: https://github.com/SAP/cgo-ase
[go]: https://golang.org/
[issues]: https://github.com/SAP/go-dblib/issues
[pkg-database-sql]: https://golang.org/pkg/database/sql
[purego]: https://github.com/SAP/go-ase
[sap-ase]: https://www.sap.com/products/sybase-ase.html
