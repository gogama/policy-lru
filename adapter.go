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

// The HandlerFunc type is an adapter that allows you to use an ordinary
// function as a Handler without implementing the Handler interface.
//
// If f is a function a signature matching that of HandlerFunc, then
// HandlerFunc[k, v](f) is a Handler that calls f.
type HandlerFunc[Key, Value any] func(k Key, v Value)

func (f HandlerFunc[Key, Value]) Removed(k Key, v Value) {
	f(k, v)
}
