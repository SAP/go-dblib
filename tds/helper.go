// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import (
	"fmt"
	"io"
)

// writeString writes s padded to padTo and its length to buf.
func writeString(stream io.Writer, s string, padTo int) error {
	if len(s) > padTo {
		return fmt.Errorf("string '%s' is too large, must be at most %d bytes long",
			s, padTo)
	}

	if _, err := stream.Write([]byte(s)); err != nil {
		return err
	}

	if _, err := stream.Write(make([]byte, padTo-len(s))); err != nil {
		return err
	}

	_, err := stream.Write([]byte{byte(len(s))})
	return err
}

func deBitmask(bitmask int, maxValue int) []int {
	curVal := 1
	ret := []int{}

	for curVal <= maxValue {
		if bitmask&curVal == curVal {
			ret = append(ret, curVal)
		}
		curVal = curVal << 1
	}

	return ret
}

func deBitmaskString(bitmask, maxValue int, toString func(i int) string, defaultValue string) string {
	values := deBitmask(bitmask, maxValue)
	if len(values) == 0 {
		return defaultValue
	}

	ret := ""
	for i, value := range values {
		ret += toString(value)
		if i+1 != len(values) {
			ret += "|"
		}
	}

	return ret
}
