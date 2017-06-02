package pillx

import (
	"sync"
)

//handler接口,ServeRouter类实现了此serve方法
type Handler interface {
	serve(*Response, IProtocol)
}

// 这里将HandlerFunc定义为一个函数类型，因此以后当调用a = HandlerFunc(f)之后, 调用a的serve实际上就是调用f的对应方法, 拥有相同参数和相同返回值的函数属于同一种类型。
type HandlerFunc func(*Response, IProtocol)

// Serve calls f(w, r).
func (f HandlerFunc) serve(w *Response, r IProtocol) {
	f(w, r)
}

type ServeRouter struct {
	mu          sync.RWMutex
	opcode_list map[uint16]OpcodeHandler
}

//将router对应的opcode,方法存储
func (router *ServeRouter) Handle(name uint16, handler Handler) {
	router.mu.Lock()
	defer router.mu.Unlock()

	router.opcode_list[name] = OpcodeHandler{handler: handler, name: name}
}

// HandleFunc registers the handler function for the given opcode.
func (rounter *ServeRouter) handleFunc(name uint16, handler func(*Response, IProtocol)) {
	rounter.Handle(name, HandlerFunc(handler))
}

// 取出opcode对应的操作方法,然后回调
func (router *ServeRouter) serve(w *Response, r IProtocol) {
	router.mu.RLock()
	defer router.mu.RUnlock()

	var handler Handler
	if router.opcode_list[r.GetCmd()].handler != nil {
		handler = router.opcode_list[r.GetCmd()].handler
		handler.serve(w, r)
	}
}

// 专门用于onmessege,close,connect等
func (router *ServeRouter) serveOnfunc(w *Response, r IProtocol, cmd uint16) {
	router.mu.RLock()
	defer router.mu.RUnlock()
	//TODO 正式上线打开
	/**
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	*/

	var handler Handler
	if router.opcode_list[cmd].handler != nil {
		handler = router.opcode_list[cmd].handler

		handler.serve(w, r)
	}
}

type OpcodeHandler struct {
	name    uint16
	handler Handler
}

// NewServeRouter allocates and returns a new ServeRouter.
func NewServeRouter() *ServeRouter { return &ServeRouter{opcode_list: make(map[uint16]OpcodeHandler)} }
