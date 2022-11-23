// Package policylru provides a generic LRU cache that lets you decide
// your own eviction policy.
package policylru

import (
	"container/list"
)

// Policy represents a cache eviction policy.
type Policy[Key, Value any] interface {
	// Evict decides whether a given cache entry should be evicted
	// from the cache based on the entry's key and value, and the
	// current number of items in the cache.
	//
	// Immediately after Evict returns true, the specified cache entry
	// will be deleted from the Cache which called Evict.
	Evict(k Key, v Value, n int) bool
}

// Handler can optionally be used to handle cache removal events.
type Handler[Key, Value any] interface {
	// Added is called after an element is added to the cache.
	Added(k Key, old, new Value, update bool)
	// Removed is called after an element is removed from the cache.
	//
	// Removal can happen either by operation of the eviction policy or
	// by a direct call to the Cache's Remove method.
	Removed(k Key, v Value)
}

// Cache is a Policy-driven LRU cache. It is not safe for concurrent
// access.
type Cache[Key comparable, Value any] struct {
	// Policy is the cache eviction policy. If Policy is nil, no element
	// will ever be evicted from the cache.
	Policy Policy[Key, Value]
	// Handler is the optional cache eviction handler.
	Handler Handler[Key, Value]

	ll    *list.List
	cache map[Key]*list.Element
}

type entry[Key, Value any] struct {
	key   Key
	value Value
}

// New creates a new policy-driven Cache.
//
// If policy is nil, the cache has no limit, and it is assumed that
// eviction is handled by the caller.
func New[Key comparable, Value any](policy Policy[Key, Value]) *Cache[Key, Value] {
	return NewWithHandler(policy, nil)
}

// NewWithHandler creates a new policy-driven Cache with a removal
// event handler.
//
// If policy is nil, the cache has no limit, and it is assumed that
// eviction is handled by the caller. If handler is nil, removal events
// will not be generated.
func NewWithHandler[Key comparable, Value any](policy Policy[Key, Value], handler Handler[Key, Value]) *Cache[Key, Value] {
	return &Cache[Key, Value]{
		Policy:  policy,
		Handler: handler,
		ll:      list.New(),
		cache:   make(map[Key]*list.Element),
	}
}

// Add adds a value to the cache.
func (c *Cache[Key, Value]) Add(k Key, v Value) {
	if c.cache == nil {
		c.ll = list.New()
		c.cache = make(map[Key]*list.Element)
	}
	h := c.Handler
	if ele, ok := c.cache[k]; ok {
		c.ll.MoveToFront(ele)
		e := ele.Value.(*entry[Key, Value])
		old := e.value
		e.value = v
		if h != nil {
			h.Added(k, old, v, true)
		}
		return
	}
	ele := c.ll.PushFront(&entry[Key, Value]{k, v})
	c.cache[k] = ele
	if h != nil {
		var old Value
		h.Added(k, old, v, false)
	}
	c.Evict()
}

// Get looks up a key's value from the cache.
func (c *Cache[Key, Value]) Get(k Key) (v Value, hit bool) {
	var ele *list.Element
	if ele, hit = c.cache[k]; hit {
		c.ll.MoveToFront(ele)
		v = ele.Value.(*entry[Key, Value]).value
	}
	return
}

// Remove removes the provided key from the cache.
func (c *Cache[Key, Value]) Remove(k Key) (removed bool) {
	if ele, hit := c.cache[k]; hit {
		c.removeElement(ele, k)
		return true
	}
	return false
}

// Evict continuously removes the oldest item from cache as long as the
// eviction policy returns true for that item. This process ends when
// the policy returns false for the oldest item or the cache is empty.
//
// The value returned is the number of items removed.
func (c *Cache[Key, Value]) Evict() (n int) {
	p := c.Policy
	if p == nil {
		return
	}
	ele := c.ll.Back()
	for ele != nil {
		e := ele.Value.(*entry[Key, Value])
		if p.Evict(e.key, e.value, c.ll.Len()) {
			c.removeElement(ele, e.key)
			n++
			ele = c.ll.Back()
		} else {
			break
		}
	}
	return
}

func (c *Cache[Key, Value]) removeElement(ele *list.Element, k Key) {
	c.ll.Remove(ele)
	delete(c.cache, k)
	h := c.Handler
	if h != nil {
		h.Removed(k, ele.Value.(*entry[Key, Value]).value)
	}
}

// Len returns the number of items in the cache.
func (c *Cache[Key, Value]) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

// Clear purges all stored items from the cache.
func (c *Cache[Key, Value]) Clear() {
	cache := c.cache
	c.ll = nil
	c.cache = nil
	h := c.Handler
	if h != nil {
		for _, ele := range cache {
			e := ele.Value.(*entry[Key, Value])
			h.Removed(e.key, e.value)
		}
	}
}
