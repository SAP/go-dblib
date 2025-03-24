// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import (
	"errors"
	"fmt"
)

// EEDError contains the extended error data packages and the wrapped
// error.
type EEDError struct {
	EEDPackages  []*EEDPackage
	WrappedError error
}

// Add adds an EED-package to an err.EEDPackages.
func (err *EEDError) Add(eed *EEDPackage) {
	err.EEDPackages = append(err.EEDPackages, eed)
}

// Is reports whether any wrapped EEDError in errs chain matches other.
func (err EEDError) Is(other error) bool {
	if err.WrappedError == nil {
		return false
	}
	return errors.Is(err.WrappedError, other)
}

// Error returns the string representation of EEDPackages.
func (err EEDError) Error() string {
	s := fmt.Sprintf("%s: received EED messages: ", err.WrappedError)

	for _, eed := range err.EEDPackages {
		s += fmt.Sprintf("%d: %s; ", eed.MsgNumber, eed.Msg)
	}

	return s
}
