//粘包问题演示客户端
package main
 
import (
    "fmt"
    "net"
    "os"
    "time"
)
 
func sender(conn net.Conn) {
    for i := 0; i < 100; i++ {
        words := "a8"
        conn.Write([]byte(words))
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
			b := make([]byte, 7)
			conn.Read(b)
			fmt.Printf("%x", b)
		}
	}()
 
    go sender(conn)
 
    for {
        time.Sleep(1 * 1e9)
    }
}