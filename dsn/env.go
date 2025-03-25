// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package dsn

import (
	"fmt"
	"os"
	"strings"
)

// FromEnv reads values from environment variables into the given input.
//
// Environment variables are matched to members of the input based on
// the json and multiref metadata tags.
//
// Example:
//   os.Setenv("MY_MEMBER", "an example")
//   os.Setenv("MY_SECOND_MEMBER", "5")
//
//   type Example struct {
//       Member string `json:"member"`
//       AnotherMember int `json:"another-member" multiref:"second-member"`
//   }
//
//   ex := new(Example)
//   if err := FromEnv("MY", ex); err != nil {
//       return err
//   }
//
// ex.Member will be "an example" and ex.AnotherMember will be 5.
//
// Also see dsn/examples/from_env.
func FromEnv(prefix string, input interface{}) error {
	// prefix = "MY_"
	prefix = strings.ToUpper(prefix) + "_"
	ttf := TagToField(input, Multiref)

	for _, env := range os.Environ() {
		// MY_EXAMPLE=VALUE=value2
		// -> []string{"MY_EXAMPLE", "VALUE=value2"}
		envSplit := strings.SplitN(env, "=", 2)
		key, value := envSplit[0], envSplit[1]

		if !strings.HasPrefix(key, prefix) {
			continue
		}

		// MY_EXAMPLE
		// -> example
		// MY_OTHER_EXAMPLE
		// -> other-example
		key = strings.ReplaceAll(
			strings.ToLower(
				strings.TrimPrefix(key, prefix),
			),
			"_", "-",
		)

		field, ok := ttf[key]
		if !ok {
			continue
		}

		if err := setValue(field, value); err != nil {
			return fmt.Errorf("dsn: error setting field %q to value %q: %w",
				key, value, err)
		}
	}

	return nil
}
