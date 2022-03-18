package echo

import (
	"io"
	"log"
	"net"
	"os"
)

// Netcat 请求
func Netcat(port string) {
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}
	remoteDone := make(chan struct{})
	localDone := make(chan struct{})
	go func() {
		mustCopy(os.Stdout, conn)
		log.Println("remoteDone")
		remoteDone <- struct{}{} // signal the main goroutine
	}()

	go func() {
		mustCopy(conn, os.Stdin)
		log.Println("输入断掉")
		localDone <- struct{}{}
	}()

	for {
		select {
		case <-remoteDone:
			conn.Close()
			return
		case <-localDone:
			conn.(*net.TCPConn).CloseWrite()
		}
	}
}

// mustCopy 拿到连接数据
func mustCopy(dst io.Writer, src io.Reader) {
	nr, err := io.Copy(dst, src)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(nr)
}
