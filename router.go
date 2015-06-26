package pillx

import (
	"sync"
)

type Handler interface {
	serve(ResponseWriter, *Request)
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler object that calls f.
type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) serve(w ResponseWriter, r *Request) {
	f(w, r)
}

type ServeRouter struct {
	mu    sync.RWMutex
	m     map[uint8]OpcodeHandler
}

// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
func (router *ServeRouter) Handle(name uint8, handler Handler) {
	router.mu.Lock()
	defer router.mu.Unlock()

	router.m[name] = OpcodeHandler{h: handler, name: name}
}

// HandleFunc registers the handler function for the given pattern.
func (rounter *ServeRouter) handleFunc(name uint8, handler func(ResponseWriter, *Request)) {
	rounter.Handle(name, HandlerFunc(handler))
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (router *ServeRouter) serve(w ResponseWriter, r *Request) {
	var handler Handler
	handler = router.m[r.name].handler
	handler.serve(w, r)
}

type OpcodeHandler struct {
	name     	uint8
	handler     Handler
}

// NewServeMux allocates and returns a new ServeMux.
func NewServeRouter() *ServeRouter { return &ServeMux{m: make(map[uint8]OpcodeHandler)} }

// defaultServeRouter is the default ServeMux used by Serve.
var defaultServeRouter = NewServeRouter()

// HandleFunc registers the handler function for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.
func HandleFunc(name uint8, handler func(ResponseWriter, *Request)) {
	defaultServeRouter.handleFunc(name, handler)
}