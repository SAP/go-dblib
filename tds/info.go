// SPDX-FileCopyrightText: 2021 SAP SE
// SPDX-FileCopyrightText: 2022 SAP SE
// SPDX-FileCopyrightText: 2023 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import (
	"fmt"
	"os"

	"github.com/SAP/go-dblib/dsn"
)

// Info is a github.com/SAP/go-dblib/dsn compatible struct with required
// information for the TDS driver.
type Info struct {
	dsn.Info

	Network        string `json:"network" doc:"Network to use, either 'tcp' or 'udp'"`
	ClientHostname string `json:"client-hostname" doc:"Hostname to send to server"`

	TLSEnable         bool   `json:"tls-enable" doc:"Enforce TLS use"`
	TLSHostname       string `json:"tls-hostname" doc:"Remote hostname to validate against SANs"`
	TLSSkipValidation bool   `json:"tls-skip-validation" doc:"Skip TLS validation - accepts any TLS certificate"`
	TLSCAFile         string `json:"tls-ca-file" doc:"Path to CA file to validate server certificate against"`

	PacketReadTimeout       int `json:"packet-read-timeout" doc:"Time in seconds to wait before aborting a connection when no response is received from the server"`
	ChannelPackageQueueSize int `json:"channel-package-queue-size" doc:"How many TDS packages can be queued in a TDS channel"`

	DebugLogPackages bool `json:"debug-log-packages" doc:"Log packages as they are transmitted/received"`
}

func SetInfo(info *Info) error {
	info.Network = "tcp"

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("dsn: error getting hostname: %w", err)
	}
	info.ClientHostname = hostname

	info.PacketReadTimeout = 50
	info.ChannelPackageQueueSize = 100

	return nil
}
