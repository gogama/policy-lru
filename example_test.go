// Copyright 2022 The policy-lru Authors. All rights reserved.
//
// Use of this source code is governed by the Apache License, Version
// 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may find a copy of the license in the file
// LICENSE or at  http://www.apache.org/licenses/LICENSE-2.0.

package policylru_test

import (
	"fmt"

	policylru "github.com/gogama/policy-lru"
)

// This example illustrates the most basic LRU cache where the eviction
// policy is based on limiting the number of keys in the cache.
func ExampleCache_withMaxCount() {
	lru := policylru.New[string, string](policylru.MaxCount[string, string](10))
	lru.Add("foo", "bar")
	value, ok := lru.Get("foo")
	fmt.Printf("In cache? %t. Value: %q.\n", ok, value)
	// Output: In cache? true. Value: "bar".
}

const maxSize = 100

type myValue struct {
	valueSize uint64
}

type myPolicy struct {
	totalSize uint64
}

func (p *myPolicy) Evict(_ string, _ myValue, _ int) bool {
	return p.totalSize > maxSize
}

func (p *myPolicy) Added(_ string, old, new myValue, _ bool) {
	p.totalSize -= old.valueSize
	p.totalSize += new.valueSize
}

func (p *myPolicy) Removed(k string, v myValue) {
	p.totalSize -= v.valueSize
	fmt.Printf("Removed %q with size %d, total size is now %d.\n", k, v.valueSize, p.totalSize)
}

// This example illustrates a more complex LRU cache where the eviction
// policy is based on the total aggregate size of the values in the
// cache.
func ExampleCache_withMaxSizePolicy() {
	policy := &myPolicy{}
	lru := policylru.NewWithHandler[string, myValue](policy, policy)
	lru.Add("foo", myValue{10})
	lru.Add("bar", myValue{90})
	lru.Add("baz", myValue{1})
	lru.Add("qux", myValue{9})
	// Output: Removed "foo" with size 10, total size is now 91.
}
