// server
package pillx

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// 用来获得当前的router并回调相应的处理方法
type ServerHandler struct {
	server *Server
}

func (sh ServerHandler) serve(res *Response, req IProtocol) {
	router := sh.server.Handler
	router.Serve(res, req)
}

var (
	bufioReaderPool   sync.Pool
	bufioWriter2kPool sync.Pool
	bufioWriter4kPool sync.Pool
	DefaultHammerTime time.Duration = 60 * time.Second
)

func bufioWriterPool(size int) *sync.Pool {
	switch size {
	case 2 << 10:
		return &bufioWriter2kPool
	case 4 << 10:
		return &bufioWriter4kPool
	}
	return nil
}

func newBufioReader(r io.Reader) *bufio.Reader {
	if v := bufioReaderPool.Get(); v != nil {
		br := v.(*bufio.Reader)
		br.Reset(r)
		return br
	}
	return bufio.NewReader(r)
}

func newBufioWriterSize(w io.Writer, size int) *bufio.Writer {
	pool := bufioWriterPool(size)
	if pool != nil {
		if v := pool.Get(); v != nil {
			bw := v.(*bufio.Writer)
			bw.Reset(w)
			return bw
		}
	}
	return bufio.NewWriterSize(w, size)
}

// 继承自mutex锁,可以安全的读取数据
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

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// noLimit is an effective infinite upper bound for io.LimitedReader
const noLimit int64 = (1 << 63) - 1

type Server struct {
	Addr     string  // TCP address to listen on, ":http" if empty
	Handler  Handler // handler to invoke, http.DefaultServeMux if nil
	Protocol IProtocol

	// ErrorLog specifies an optional logger for errors accepting
	// connections and unexpected behavior from handlers.
	// If nil, logging goes to os.Stderr via the log package's
	// standard logger.
	ErrorLog *log.Logger
}

//注册路由相应处理方法
func (srv *Server) HandleFunc(name uint16, handler func(*Response, IProtocol)) {
	if srv.Handler == nil {
		srv.Handler = NewServeRouter()
	}

	srv.Handler.(*ServeRouter).handleFunc(name, handler)
}

func (srv *Server) newConn(rc net.Conn) (c *Conn, err error) {
	c = new(Conn)
	c.remote_addr = rc.RemoteAddr().String()
	c.server = srv
	c.remonte_conn = rc
	c.io_writer = rc

	c.sr = liveSwitchReader{r: c.remonte_conn}
	c.lr = io.LimitReader(&c.sr, noLimit).(*io.LimitedReader)
	br := newBufioReader(c.lr) //bufio pool
	bw := newBufioWriterSize(checkConnErrorWriter{c}, 4<<10)
	c.buf = bufio.NewReadWriter(br, bw)
	return c, nil
}

// ListenAndServe listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.  If
// srv.Addr is blank, ":5917" is used.
func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":5917"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}

func (srv *Server) ListenAndServeUdp() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":5917"
	}
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	for {
		c, err := srv.newConn(conn)
		if err != nil {
			continue
		}
		go c.serve()
	}
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

/*
handleSignals listens for os Signals and calls any hooked in function that the
user had registered with the signal.
*/
func (srv *Server) handleSignals() {
	var sig os.Signal
	hookableSignals := []os.Signal{
		syscall.SIGHUP,
		//syscall.SIGUSR1,
		//syscall.SIGUSR2,
		syscall.SIGINT,
		syscall.SIGTERM,
		//syscall.SIGTSTP,
		syscall.SIGQUIT,
	}
	sigChan := make(chan os.Signal)

	signal.Notify(
		sigChan,
		hookableSignals...,
	)

	pid := syscall.Getpid()
	for {
		sig = <-sigChan
		switch sig {
		case syscall.SIGHUP:
			log.Println(pid, "Received SIGHUP. forking.")
		case syscall.SIGINT:
			log.Println(pid, "Received SIGINT.")
			srv.shutdown()
		case syscall.SIGTERM:
			log.Println(pid, "Received SIGTERM.")
			srv.shutdown()
		default:
			log.Printf("Received %v: nothing i care about...\n", sig)
		}
		break
	}
}

func (srv *Server) shutdown() {
	time.Sleep(DefaultHammerTime)
	return
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
	c *Conn
}

func (w checkConnErrorWriter) Write(p []byte) (n int, err error) {
	n, err = w.c.io_writer.Write(p) // c.w == c.rwc, except after a hijack, when rwc is nil.
	if err != nil && w.c.io_writer_err == nil {
		w.c.io_writer_err = err
	}
	return
}
