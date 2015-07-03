package pillx

import (
	"sync"
)

type Handler interface {
	serve(*Response, *Request)
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler object that calls f.
type HandlerFunc func(*Response, *Request)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) serve(w *Response, r *Request) {
	f(w, r)
}

type ServeRouter struct {
	mu    sync.RWMutex
	opcode_list     map[uint16]OpcodeHandler
}

// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
func (router *ServeRouter) Handle(name uint16, handler Handler) {
	router.mu.Lock()
	defer router.mu.Unlock()

	router.opcode_list[name] = OpcodeHandler{handler: handler, name: name}
}

// HandleFunc registers the handler function for the given pattern.
func (rounter *ServeRouter) handleFunc(name uint16, handler func(*Response, *Request)) {
	rounter.Handle(name, HandlerFunc(handler))
}

// Serve dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (router *ServeRouter) serve(w *Response, r *Request) {
	router.mu.RLock()
	defer router.mu.RUnlock()
	
	var handler Handler
	handler = router.opcode_list[r.name].handler
	handler.serve(w, r)
}

type OpcodeHandler struct {
	name     	uint16
	handler     Handler
}

// NewServeMux allocates and returns a new ServeMux.
func NewServeRouter() *ServeRouter { return &ServeRouter{opcode_list: make(map[uint16]OpcodeHandler)} }

// defaultServeRouter is the default ServeMux used by Serve.
var defaultServeRouter = NewServeRouter()

// HandleFunc registers the handler function for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.
func HandleFunc(name uint16, handler func(*Response, *Request)) {
	defaultServeRouter.handleFunc(name, handler)
}