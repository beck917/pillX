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

// A conn represents the server side of connection.
type conn struct {
	remote_addr string
	server 		*Server
	rc 			net.Conn
	w 			io.Writer
	sr         	liveSwitchReader     // where the LimitReader reads from; usually the rwc
	lr         	*io.LimitedReader    // io.LimitReader(sr)
	buf 		*bufio.ReadWriter
	
	mu 			sync.Mutex
}

func (srv *Server) newConn(rc net.Conn) (c *conn, err error) {
	c = new(conn);
	c.remote_addr = rc.RemoteAddr().String()
	c.server = srv
	c.rc = rc
	c.w = w
	
	c.sr = liveSwitchReader{r: c.rc}
	c.lr = io.LimitReader(&c.sr, noLimit).(*io.LimitedReader)
	br := newBufioReader(c.lr)
	bw := newBufioWriterSize(checkConnErrorWriter{c}, 4<<10)
	c.buf = bufio.NewReadWriter(br, bw)
	return c, nil
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
	
	// ErrorLog specifies an optional logger for errors accepting
	// connections and unexpected behavior from handlers.
	// If nil, logging goes to os.Stderr via the log package's
	// standard logger.
	ErrorLog *log.Logger
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
		c.setState(c.rc, StateNew) // before Serve can return
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