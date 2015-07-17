package pillx

import (
	"sync"
)

//handler接口,ServeRouter类实现了此serve方法
type Handler interface {
	serve(*Response, *Request)
}

// 这里将HandlerFunc定义为一个函数类型，因此以后当调用a = HandlerFunc(f)之后, 调用a的serve实际上就是调用f的对应方法, 拥有相同参数和相同返回值的函数属于同一种类型。
type HandlerFunc func(*Response, *Request)

// Serve calls f(w, r).
func (f HandlerFunc) serve(w *Response, r *Request) {
	f(w, r)
}

type ServeRouter struct {
	mu    sync.RWMutex
	opcode_list     map[uint16]OpcodeHandler
}

//将router对应的opcode,方法存储
func (router *ServeRouter) Handle(name uint16, handler Handler) {
	router.mu.Lock()
	defer router.mu.Unlock()

	router.opcode_list[name] = OpcodeHandler{handler: handler, name: name}
}

// HandleFunc registers the handler function for the given opcode.
func (rounter *ServeRouter) handleFunc(name uint16, handler func(*Response, *Request)) {
	rounter.Handle(name, HandlerFunc(handler))
}

// 取出opcode对应的操作方法,然后回调
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

// NewServeRouter allocates and returns a new ServeRouter.
func NewServeRouter() *ServeRouter { return &ServeRouter{opcode_list: make(map[uint16]OpcodeHandler)} }

// defaultServeRouter is the default ServeRouter used by Serve.
var defaultServeRouter = NewServeRouter()

// HandleFunc registers the handler function for the given opcode
// in the defaultServeRouter.
func HandleFunc(name uint16, handler func(*Response, *Request)) {
	defaultServeRouter.handleFunc(name, handler)
}