package policylru

type maxCountPolicy[Key, Value any] int

func (p maxCountPolicy[Key, Value]) Evict(_ Key, _ Value, n int) bool {
	return n > int(p)
}

func MaxCount[Key, Value any](n int) Policy[Key, Value] {
	return maxCountPolicy[Key, Value](n)
}
