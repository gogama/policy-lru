// Copyright 2022 The policy-lru Authors. All rights reserved.
//
// Use of this source code is governed by the Apache License, Version
// 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may find a copy of the license in the file
// LICENSE or at  http://www.apache.org/licenses/LICENSE-2.0.

package policylru

// The PolicyFunc type is an adapter that allows you to use an ordinary
// function as a Policy without implementing the Policy interface.
//
// If f is a function a signature matching that of PolicyFunc, then
// PolicyFunc[k, v](f) is a Policy that calls f.
type PolicyFunc[Key, Value any] func(k Key, v Value, n int) bool

func (f PolicyFunc[Key, Value]) Evict(k Key, v Value, n int) bool {
	return f(k, v, n)
}

// The AddedFunc type is an adapter that allows you to use a single
// ordinary add-handling function as a Handler without implementing
// the whole Handler interface.
//
// If f is a function whose signature matches the Added method of a
// Handler[k, v], then AddedFunc[k, v](f) is a Handler[k, v] with a
// no-op Removed method and an Added method that calls f.
type AddedFunc[Key, Value any] func(k Key, old, new Value, updated bool)

func (f AddedFunc[Key, Value]) Added(k Key, old, new Value, updated bool) {
	f(k, old, new, updated)
}

func (f AddedFunc[Key, Value]) Removed(k Key, v Value) {
}

// The RemovedFunc type is an adapter that allows you to use a single
// ordinary remove-handling function as a Handler without implementing
// the whole Handler interface.
//
// If f is a function whose signature matches the Removed method of a
// Handler[k, v], then RemovedFunc[k, v](f) is a Handler[k, v] with a
// no-op Added method and a Removed method that calls f.
type RemovedFunc[Key, Value any] func(k Key, v Value)

func (f RemovedFunc[Key, Value]) Added(_ Key, _, _ Value, _ bool) {
}

func (f RemovedFunc[Key, Value]) Removed(k Key, v Value) {
	f(k, v)
}
