// SPDX-FileCopyrightText: 2021 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package dsn

import (
	"reflect"
	"testing"
)

func TestTagToField(t *testing.T) {
	type Outer struct {
		A string `json:"A" multiref:"anotherA"`
		B int    `json:"B" multiref:"anotherB"`
	}

	type Embed struct {
		EmbedA string `json:"EmbedA" multiref:"anotherEmbedA"`
	}

	type OuterEmbed struct {
		A string `json:"A" multiref:"anotherA"`
		B int    `json:"B" multiref:"anotherB"`
		Embed
	}

	simple := new(Outer)
	simple.A = "a string"
	simple.B = 5

	t.Run("simple", func(t *testing.T) {
		checker(t, false,
			TagToField(*simple, OnlyJSON),
			map[string]reflect.Value{
				"A": reflect.ValueOf("a string"),
				"B": reflect.ValueOf(5),
			},
		)
	})

	t.Run("simple ptr", func(t *testing.T) {
		checker(t, true,
			TagToField(simple, OnlyJSON),
			map[string]reflect.Value{
				"A": reflect.ValueOf("a string"),
				"B": reflect.ValueOf(5),
			},
		)
	})

	t.Run("simple multiref", func(t *testing.T) {
		checker(t, false,
			TagToField(*simple, Multiref),
			map[string]reflect.Value{
				"A":        reflect.ValueOf("a string"),
				"anotherA": reflect.ValueOf("a string"),
				"B":        reflect.ValueOf(5),
				"anotherB": reflect.ValueOf(5),
			},
		)
	})

	t.Run("simple multiref ptr", func(t *testing.T) {
		checker(t, true,
			TagToField(simple, Multiref),
			map[string]reflect.Value{
				"A":        reflect.ValueOf("a string"),
				"anotherA": reflect.ValueOf("a string"),
				"B":        reflect.ValueOf(5),
				"anotherB": reflect.ValueOf(5),
			},
		)
	})

	embed := new(OuterEmbed)
	embed.A = "a string"
	embed.B = 5
	embed.EmbedA = "an embedded string"

	t.Run("embed", func(t *testing.T) {
		checker(t, false,
			TagToField(*embed, OnlyJSON),
			map[string]reflect.Value{
				"A":      reflect.ValueOf("a string"),
				"B":      reflect.ValueOf(5),
				"EmbedA": reflect.ValueOf("an embedded string"),
			},
		)
	})

	t.Run("embed ptr", func(t *testing.T) {
		checker(t, true,
			TagToField(embed, OnlyJSON),
			map[string]reflect.Value{
				"A":      reflect.ValueOf("a string"),
				"B":      reflect.ValueOf(5),
				"EmbedA": reflect.ValueOf("an embedded string"),
			},
		)
	})
}

func checker(t *testing.T, isPtr bool, recv, expect map[string]reflect.Value) {
	if len(recv) != len(expect) {
		t.Errorf("Expected %d members in returned map, got %d", len(expect), len(recv))
	}

	for key, expect := range expect {
		received, ok := recv[key]
		if !ok {
			t.Errorf("Expected key %q not in returned map", key)
			continue
		}

		if received.Type() != expect.Type() {
			t.Errorf("Expected key %q to hold value of type %s, got %s",
				key, expect.Type(), received.Type())
			continue
		}

		if received.Interface() != expect.Interface() {
			t.Errorf("Expected key %q to hold value %q, got %q",
				key, expect.Interface(), received.Interface())
		}

		if isPtr {
			if !received.CanSet() {
				t.Errorf("Received reflect.Value for key %q is not settable", key)
				continue
			}
		}
	}
}
