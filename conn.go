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

// A conn represents the server side of connection.
type conn struct {
	remote_addr string
	server 		*Server
	remonte_conn 			net.Conn
	io_writer			io.Writer
	sr         	liveSwitchReader     // where the LimitReader reads from; usually the rwc
	lr         	*io.LimitedReader    // io.LimitReader(sr)
	buf 		*bufio.ReadWriter
	
	mu 			sync.Mutex
}

func (c *conn) readRequest() (response *response, err error) {
	
}

// Serve a new connection.
func (c *conn) serve() {
	origConn := c.rwc // copy it before it's set nil on Close or Hijack
	for {
		w, err := c.readRequest()
	}
	ServerHandler{c.server}.serve(w, w.req)
}