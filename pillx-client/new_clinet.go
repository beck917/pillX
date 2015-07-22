//粘包问题演示客户端
package main
 
import (
    "fmt"
    "net"
    "os"
    "time"
	"encoding/binary"
	"bytes"
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
    for i := 0; i < 100; i++ {
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
		fmt.Printf("%x", buf.Bytes())
		
        conn.Write(buf.Bytes())
    }
}
 
func main() {
    server := "127.0.0.1:8080"
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
 
    fmt.Println("connect success")
	
	go func() {
		for {
			b := make([]byte, 5)
			conn.Read(b)
			fmt.Printf("%s", b)
		}
	}()
 
    go sender(conn)
 
    for {
        time.Sleep(1 * 1e9)
    }
}