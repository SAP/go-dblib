// SPDX-FileCopyrightText: 2020 SAP SE
// SPDX-FileCopyrightText: 2021 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package dsn

// Info serves as both an example and an embeddable default to use in
// DSN structs.
type Info struct {
	Host     string `json:"host" multiref:"hostname" doc:"Hostname to connect to"`
	Port     string `json:"port" doc:"Port (Example: '443' or 'tls') to connect to"`
	Username string `json:"username" multiref:"user" doc:"Username"`
	Password string `json:"password" multiref:"passwd,pass" doc:"Password"`
	Database string `json:"database" multiref:"db" doc:"Database"`
}
