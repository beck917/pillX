// server
package pillx

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// serverHandler delegates to either the server's Handler or
// DefaultServeMux and also handles "OPTIONS *" requests.
type ServerHandler struct {
	server *Server
}

func (sh ServerHandler) serve(rw ResponseWriter, req *Request) {
	handler := sh.server.Handler
	if handler == nil {
		handler = defaultServeRouter
	}

	handler.serve(rw, req)
}

// A liveSwitchReader can have its Reader changed at runtime. It's
// safe for concurrent reads and switches, if its mutex is held.
type liveSwitchReader struct {
	sync.Mutex
	r io.Reader
}

func (sr *liveSwitchReader) Read(p []byte) (n int, err error) {
	sr.Lock()
	r := sr.r
	sr.Unlock()
	return r.Read(p)
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

type Server struct {
	addr			string        // TCP address to listen on, ":http" if empty
	Handler        Handler       // handler to invoke, http.DefaultServeMux if nil
	
	// ErrorLog specifies an optional logger for errors accepting
	// connections and unexpected behavior from handlers.
	// If nil, logging goes to os.Stderr via the log package's
	// standard logger.
	ErrorLog *log.Logger
}

func (srv *Server) newConn(rc net.Conn) (c *conn, err error) {
	c = new(conn);
	c.remote_addr = rc.RemoteAddr().String()
	c.server = srv
	c.remonte_conn = rc
	c.io_writer = rc
	
	c.sr = liveSwitchReader{r: c.rc}
	c.lr = io.LimitReader(&c.sr, noLimit).(*io.LimitedReader)
	br := newBufioReader(c.lr)
	bw := newBufioWriterSize(checkConnErrorWriter{c}, 4<<10)
	c.buf = bufio.NewReadWriter(br, bw)
	return c, nil
}

// ListenAndServe listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.  If
// srv.Addr is blank, ":5917" is used.
func (srv *Server) ListenAndServe() error {
	addr := srv.addr
	if addr == "" {
		addr = ":5917"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each.  The service goroutines read requests and
// then call srv.Handler to reply to them.
func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		//remote connnetion
		rc, e := l.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				srv.logf("http: Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return e
		}
		tempDelay = 0
		c, err := srv.newConn(rc)
		if err != nil {
			continue
		}
		//c.setState(c.rc, StateNew) // before Serve can return
		go c.serve()
	}
}

func (s *Server) logf(format string, args ...interface{}) {
	if s.ErrorLog != nil {
		s.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

// checkConnErrorWriter writes to c.rwc and records any write errors to c.werr.
// It only contains one field (and a pointer field at that), so it
// fits in an interface value without an extra allocation.
type checkConnErrorWriter struct {
	c *conn
}

func (w checkConnErrorWriter) Write(p []byte) (n int, err error) {
	n, err = w.c.w.Write(p) // c.w == c.rwc, except after a hijack, when rwc is nil.
	if err != nil && w.c.werr == nil {
		w.c.werr = err
	}
	return
}