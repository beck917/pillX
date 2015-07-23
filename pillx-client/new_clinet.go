//粘包问题演示客户端
package main
 
import (
    "fmt"
    "net"
    "os"
    "time"
	"encoding/binary"
	"bytes"
	"sync"
)

type RequestHeader struct {
	mark		uint8
	cmd 		uint16
	error		uint16
	size		uint16
}

type Request struct {
	Header		*RequestHeader
	Content		[]byte
}

 
func sender(conn net.Conn) {
    for i := 0; i < 1000; i++ {
		reqHeader := &RequestHeader{
			mark:	0xA8,
			size:	5,
			cmd:	0x0DDC,
			error:	0,
		}

		buf := new(bytes.Buffer)
		binary.Write(buf, binary.BigEndian, reqHeader)
		errBin := binary.Write(buf, binary.BigEndian, []byte("Hello"))
		
		if (errBin != nil) {
			fmt.Println("binary.Write failed:", errBin)
			return
		}
		//fmt.Printf("%x", buf.Bytes())
		
        conn.Write(buf.Bytes())
    }
}

type Counter struct {
      mu sync.Mutex
      x int64
}

func (c *Counter) Inc(){
      c.mu.Lock()
      defer c.mu.Unlock()
      c.x++
}
 
func main() {
    server := "127.0.0.1:8080"
	start := time.Now()
	fmt.Print(start)
	c := Counter{}
	for i := 0; i < 1000; i++ {
	    tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	    if err != nil {
	        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
	        os.Exit(1)
	    }
	 
	    conn, err := net.DialTCP("tcp", nil, tcpAddr)
	    if err != nil {
	        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
	        os.Exit(1)
	    }
	 
	    defer conn.Close()
	 
	    //fmt.Println("connect success")
		
		go func() {
			for {
				b := make([]byte, 5)
				conn.Read(b)
				//fmt.Printf("%s", b)
				c.mu.Lock()
				c.x++
				//fmt.Print(t)
				if c.x >= 990000 {
				end := time.Now()
				fmt.Print(end)
				}
				c.mu.Unlock()
			}
		}()
	 
	    go sender(conn)
 	}
    for {
        time.Sleep(1 * 1e9)
    }
}