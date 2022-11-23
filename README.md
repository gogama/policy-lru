# policy-lru

A very simple LRU cache in GoLang supporting Go generics and defining
your own custom eviction policies.

Use It If You Want
==================

1. A lightweight LRU with no external dependencies.
2. Go 1.18 generics.
3. To define your own eviction policy, for example only evict from the
   cache when consumed memory or disk space exceeds some threshold.

Examples
========

Below is an example of the most basic LRU cache with a simple eviction
policy based on a maximum number of keys. More examples are available in
the [API reference documentation](https://pkg.go.dev/github.com/gogama/policy-lru.

```go
package main

import (
	"fmt"
	"github.com/gogama/policy-lru"
)

func main() {
	lru := policylru.New[string, string](policylru.MaxCount[string, string](10))
	lru.Add("foo", "bar")
	value, ok := lru.Get("foo")
	fmt.Printf("In cache? %t. Value: %q.\n", ok, value)
}
```

Status
======

[![Go Report Card](https://goreportcard.com/badge/github.com/gogama/policy-lru)](https://goreportcard.com/report/github.com/gogama/policy-lru) [![PkgGoDev](https://pkg.go.dev/badge/github.com/gogama/policy-lru)](https://pkg.go.dev/github.com/gogama/policy-lru)

License
=======

This project is licensed under the terms of the Apache License 2.0.

Acknowledgements
================

The code for this project is *heavily* based on the code from package
[`lru`](https://pkg.go.dev/github.com/golang/groupcache/lru) from
[github.com/golang/groupcache](https://github.com/golang/groupcache).
The borrowed code is licensed under the
[Apache License 2.0](https://github.com/golang/groupcache/blob/master/LICENSE)

Developer happiness on this project was boosted by JetBrains' generous donation
of an [open source license](https://www.jetbrains.com/opensource/) for their
lovely GoLand IDE. ❤