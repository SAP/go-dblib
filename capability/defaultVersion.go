// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package capability

// DefaultVersion implements the Version interface.
type DefaultVersion struct {
	// Spec is the string representation of the version.
	// Depending on the used versioning scheme these may vary wildly.
	spec         string
	capabilities map[*Capability]bool
}

// NewDefaultVersion returns an initialized DefaultVersion.
func NewDefaultVersion(spec string) Version {
	return &DefaultVersion{spec, map[*Capability]bool{}}
}

// VersionString implements the capability.Version interface.
func (v DefaultVersion) VersionString() string {
	return v.spec
}

// SetCapability implements the capability.Version interface.
func (v *DefaultVersion) SetCapability(cap *Capability, b bool) {
	v.capabilities[cap] = b
}

// Has implements the capability.Version interface.
func (v DefaultVersion) Has(cap *Capability) bool {
	canCap, ok := v.capabilities[cap]
	if !ok {
		return false
	}
	return canCap
}
