package task1

import (
	"io"
	"log"
	"net"
	"time"
)

// Clock 时钟
func Clock(port string, tz string) {
	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn, tz)
	}
}

// handleConn 链接
func handleConn(c net.Conn, tz string) {
	var cstSh, _ = time.LoadLocation(tz)
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().In(cstSh).Format("15:04:05"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}
