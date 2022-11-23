// Copyright 2022 The policy-lru Authors. All rights reserved.
//
// Use of this source code is governed by the Apache License, Version
// 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may find a copy of the license in the file
// LICENSE or at  http://www.apache.org/licenses/LICENSE-2.0.

package policylru

type maxCountPolicy[Key, Value any] int

func (p maxCountPolicy[Key, Value]) Evict(_ Key, _ Value, n int) bool {
	return n > int(p)
}

// MaxCount returns a Policy that evicts the oldest key from the Cache
// whenever the total number of keys in the cash exceeds the given
// maximum count.
func MaxCount[Key, Value any](n int) Policy[Key, Value] {
	return maxCountPolicy[Key, Value](n)
}
