package policylru

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type simpleStruct struct {
	int
	string
}

type complexStruct struct {
	int
	simpleStruct
}

func TestZeroValue(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		var lru Cache[int, float64]

		lru.Add(1, 2.0)
		lru.Add(2, 3.0)

		assert.Equal(t, 2, lru.Len())
	})

	t.Run("clear", func(t *testing.T) {
		var lru Cache[bool, struct{}]

		lru.Clear()

		assert.Equal(t, 0, lru.Len())
	})
}

func TestAddAndGet(t *testing.T) {
	t.Run("string_hit", func(t *testing.T) {
		lru := New[string, int](nil)

		lru.Add("foo", 1234)
		value, ok := lru.Get("foo")

		assert.Equal(t, 1, lru.Len())
		assert.True(t, ok)
		assert.Equal(t, 1234, value)
	})

	t.Run("string_miss", func(t *testing.T) {
		lru := New[string, string](nil)

		val, ok := lru.Get("foo")

		assert.Equal(t, 0, lru.Len())
		assert.False(t, ok)
		assert.Equal(t, "", val)
	})

	t.Run("simple_struct_hit", func(t *testing.T) {
		lru := New[simpleStruct, time.Time](nil)

		key := simpleStruct{1, "two"}
		now := time.Now()
		lru.Add(key, now)
		value, ok := lru.Get(key)

		assert.Equal(t, 1, lru.Len())
		assert.True(t, ok)
		assert.Equal(t, now, value)
	})

	t.Run("simple_struct_miss", func(t *testing.T) {
		lru := New[simpleStruct, float64](nil)

		lru.Add(simpleStruct{1, "one"}, 1.0)
		lru.Add(simpleStruct{2, "two"}, 2.0)
		value, ok := lru.Get(simpleStruct{3, "three"})

		assert.Equal(t, 2, lru.Len())
		assert.False(t, ok)
		assert.Equal(t, 0.0, value)
	})

	t.Run("complex_struct_hit", func(t *testing.T) {
		lru := New[complexStruct, int](nil)

		key := complexStruct{1, simpleStruct{2, "three"}}
		lru.Add(key, 4)
		value, ok := lru.Get(key)

		assert.Equal(t, 1, lru.Len())
		assert.True(t, ok)
		assert.Equal(t, 4, value)
	})

	t.Run("with_policy", func(t *testing.T) {
		lru := New[byte, bool](MaxCount[byte, bool](1))

		lru.Add('x', true)
		lru.Add('y', false)
		value1, ok1 := lru.Get('x')
		value2, ok2 := lru.Get('y')

		assert.Equal(t, 1, lru.Len())
		assert.False(t, ok1)
		assert.Equal(t, false, value1)
		assert.True(t, ok2)
		assert.Equal(t, false, value2)
	})

	t.Run("with_added_handler", func(t *testing.T) {
		var olds, news []string
		var updateds []bool
		lru := NewWithHandler[string, string](MaxCount[string, string](2), AddedFunc[string, string](func(k string, old, new string, updated bool) {
			olds = append(olds, k, old)
			news = append(news, k, new)
			updateds = append(updateds, updated)
		}))

		lru.Add("foo", "bar")
		lru.Add("foo", "baz")
		lru.Add("hello", "world")
		lru.Add("foo", "qux")
		lru.Add("bing", "bong")
		lru.Add("foo", "pew")
		lru.Add("hello", "folks-people")
		value, ok := lru.Get("foo")

		assert.Equal(t, 2, lru.Len())
		assert.True(t, ok)
		assert.Equal(t, "pew", value)
		assert.Equal(t, []string{"foo", "", "foo", "bar", "hello", "", "foo", "baz", "bing", "", "foo", "qux", "hello", ""}, olds)
		assert.Equal(t, []string{"foo", "bar", "foo", "baz", "hello", "world", "foo", "qux", "bing", "bong", "foo", "pew", "hello", "folks-people"}, news)
		assert.Equal(t, []bool{false, true, false, true, false, true, false}, updateds)
	})

	t.Run("with_removed_handler", func(t *testing.T) {
		var removedCount int
		var removedKey string
		var removedValue string
		var removedTime time.Time
		lru := NewWithHandler[string, string](MaxCount[string, string](2), RemovedFunc[string, string](func(k string, v string) {
			removedCount++
			removedKey = k
			removedValue = v
			removedTime = time.Now()
		}))

		lru.Add("foo", "bar")
		lru.Add("baz", "qux")
		before := time.Now()
		lru.Add("razzle", "dazzle")
		after := time.Now()
		value1, ok1 := lru.Get("foo")
		value2, ok2 := lru.Get("razzle")

		assert.Equal(t, 2, lru.Len())
		assert.Equal(t, 1, removedCount)
		assert.False(t, removedTime.Before(before))
		assert.False(t, after.Before(removedTime))
		assert.False(t, ok1)
		assert.Equal(t, "foo", removedKey)
		assert.Equal(t, "bar", removedValue)
		assert.Equal(t, "", value1)
		assert.True(t, ok2)
		assert.Equal(t, "dazzle", value2)
	})
}

func TestRemove(t *testing.T) {
	t.Run("removed", func(t *testing.T) {
		lru := New[string, int](nil)

		lru.Add("foo", 1001)
		removed := lru.Remove("foo")

		assert.True(t, removed)
		assert.Equal(t, 0, lru.Len())
	})

	t.Run("not_removed", func(t *testing.T) {
		lru := New[int, int](nil)

		removed := lru.Remove(0)

		assert.False(t, removed)
		assert.Equal(t, 0, lru.Len())
	})

	t.Run("with_removed_handler", func(t *testing.T) {
		var removedKey int
		var removedValue string
		var removedTime time.Time
		lru := NewWithHandler[int, string](nil, RemovedFunc[int, string](func(k int, v string) {
			removedKey = k
			removedValue = v
			removedTime = time.Now()
		}))

		lru.Add(10, "lorem")
		lru.Add(15, "ipsum")
		before := time.Now()
		shouldBeTrue := lru.Remove(15)
		after := time.Now()
		shouldBeFalse := lru.Remove(15)

		assert.True(t, shouldBeTrue)
		assert.Equal(t, 15, removedKey)
		assert.Equal(t, "ipsum", removedValue)
		assert.False(t, removedTime.Before(before))
		assert.False(t, after.Before(removedTime))
		assert.Equal(t, 1, lru.Len())
		assert.False(t, shouldBeFalse)
	})
}

func TestEvict(t *testing.T) {
	t.Run("implicit_during_add", func(t *testing.T) {
		lru := New[int, int](MaxCount[int, int](2))

		lru.Add(1, 11)
		lru.Add(2, 22)
		lru.Add(1, 11)
		lru.Add(3, 33)
		value, ok := lru.Get(2)

		assert.Equal(t, 2, lru.Len())
		assert.Equal(t, 0, value)
		assert.False(t, ok)
	})

	t.Run("explicit", func(t *testing.T) {
		maxSize := 10
		policy := PolicyFunc[string, string](func(_, _ string, n int) bool {
			return n > maxSize
		})
		lru := New[string, string](policy)

		lru.Add("doomed", "to eviction")
		lru.Add("ill-fated", "due to being an eviction target")
		lru.Add("lucky", "to survive")
		lru.Add("blessed", "to avoid the evict-pocalypse")

		assert.Equal(t, 4, lru.Len())

		maxSize = 2
		lru.Evict()
		_, ok1 := lru.Get("ill-fated")
		value2, ok2 := lru.Get("lucky")
		value3, ok3 := lru.Get("blessed")

		assert.Equal(t, 2, lru.Len())
		assert.False(t, ok1)
		assert.True(t, ok2)
		assert.Equal(t, "to survive", value2)
		assert.True(t, ok3)
		assert.Equal(t, "to avoid the evict-pocalypse", value3)
	})
}

func TestClear(t *testing.T) {
	var removed []int
	lru := NewWithHandler[int, int](nil, RemovedFunc[int, int](func(k, v int) {
		removed = append(removed, k, v)
	}))

	lru.Add(1, 2)
	lru.Add(3, 4)
	lru.Add(5, 6)
	lru.Clear()

	assert.Equal(t, 0, lru.Len())
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, removed)
}

/*
func TestEvict(t *testing.T) {
	evictedKeys := make([]Key, 0)
	onEvictedFun := func(key Key, value interface{}) {
		evictedKeys = append(evictedKeys, key)
	}

	lru := New(20)
	lru.OnEvicted = onEvictedFun
	for i := 0; i < 22; i++ {
		lru.Add(fmt.Sprintf("myKey%d", i), 1234)
	}

	if len(evictedKeys) != 2 {
		t.Fatalf("got %d evicted keys; want 2", len(evictedKeys))
	}
	if evictedKeys[0] != Key("myKey0") {
		t.Fatalf("got %v in first evicted key; want %s", evictedKeys[0], "myKey0")
	}
	if evictedKeys[1] != Key("myKey1") {
		t.Fatalf("got %v in second evicted key; want %s", evictedKeys[1], "myKey1")
	}
}
*/
