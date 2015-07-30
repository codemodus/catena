// Package catena aids the composition of http.Handler wrapper catenas.
//
// Review the test file for examples covering catena manipulation. Benchmarks
// are available showing a negligible increase in processing time and memory
// consumption, and no increase in memory allocations compared to nesting
// functions without an aid.
package catena

import "net/http"

// Catena holds the basic components used to order http.Handler wrapper catenas.
type Catena struct {
	hws []func(http.Handler) http.Handler
}

// New takes one or more http.Handler wrappers, and returns a new Catena.
func New(hws ...func(http.Handler) http.Handler) Catena {
	return Catena{hws: hws}
}

// Append takes one or more http.Handler wrappers, and appends the value to the
// returned Catena.
func (c Catena) Append(hws ...func(http.Handler) http.Handler) Catena {
	c.hws = append(c.hws, hws...)
	return c
}

// Merge takes one or more Catena objects, and appends the values' http.Handler
// wrappers to the returned Catena.
func (c Catena) Merge(cs ...Catena) Catena {
	for k := range cs {
		c.hws = append(c.hws, cs[k].hws...)
	}
	return c
}

// End takes a http.Handler and returns an http.Handler.
func (c Catena) End(h http.Handler) http.Handler {
	if h == nil {
		h = http.HandlerFunc(nilHandler)
	}

	for i := len(c.hws) - 1; i >= 0; i-- {
		h = c.hws[i](h)
	}

	return h
}

// EndFn takes a func that matches the http.HandlerFunc type, then passes it to
// End.
func (c Catena) EndFn(h http.HandlerFunc) http.Handler {
	if h == nil {
		h = http.HandlerFunc(nilHandler)
	}
	return c.End(h)
}

func nilHandler(w http.ResponseWriter, r *http.Request) {
	return
}
