// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package namepool

import (
	"sync"
	"testing"
)

func TestNewPool(t *testing.T) {
	cases := map[string]struct {
		format string
		first  string
	}{
		"empty": {
			format: "",
			first:  "%!(EXTRA uint64=1)",
		},
		"only id": {
			format: "%d",
			first:  "1",
		},
		"format with only %d": {
			format: "name %d",
			first:  "name 1",
		},
		"format with multiple formats": {
			format: "name %d %d %s",
			first:  "name 1 %!d(MISSING) %!s(MISSING)",
		},
		"format with string verb": {
			format: "name %s",
			first:  "name %!s(uint64=1)",
		},
	}

	for title, cas := range cases {
		t.Run(title, func(t *testing.T) {
			pool := Pool(cas.format)
			name := pool.Acquire()
			defer pool.Release(name)

			if name.Name() != cas.first {
				t.Errorf("Expected to receive '%s' as first name, received: %s", cas.first, name.Name())
			}
		})
	}
}

func TestPool_AcquireConcurrent(t *testing.T) {
	pool := Pool("")

	fn := func(start, finished *sync.WaitGroup, t *testing.T) {
		// Wait for all goroutines to be started by the main goroutine,
		// then execute the test.
		start.Wait()
		defer finished.Done()

		name := pool.Acquire()
		if name == nil {
			t.Errorf("pool.Acquire returned nil value")
			return
		}
		defer pool.Release(name)

		if name.name == "" {
			t.Errorf("name.name is empty")
		}

		if *name.id == 0 {
			t.Errorf("name.id is 0")
		}
	}

	concurrents := 1000

	start := new(sync.WaitGroup)
	start.Add(1)

	finished := new(sync.WaitGroup)
	finished.Add(concurrents)

	for i := 0; i < concurrents; i++ {
		go fn(start, finished, t)
	}

	start.Done()
	finished.Wait()
}

func TestPool_Release(t *testing.T) {
	pool := Pool("%d")

	name := pool.Acquire()
	if name == nil {
		t.Errorf("Acquired Name is nil")
		return
	}

	pool.Release(name)
	if name.id != nil {
		t.Errorf("Released Name has non-nil ID pointer")
	}
	if (*name).name != "" {
		t.Errorf("Released Name has non-empty name")
	}
}

func TestPool_ReleaseMultiple(t *testing.T) {
	pool := Pool("")

	name := pool.Acquire()

	name.Release()
	if name.id != nil {
		t.Errorf("Release Name has non-nil ID pointer")
	}

	// Release a second time to check that no nil pointer gets into the
	// pool.
	name.Release()

	secondName := pool.Acquire()
	if secondName.id == nil {
		t.Errorf("Received nil ID pointer from pool")
	}
}
