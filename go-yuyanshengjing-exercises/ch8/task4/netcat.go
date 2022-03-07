package task4

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
	done := make(chan struct{})
	go func() {
		mustCopy(os.Stdout, conn)
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	mustCopy(conn, os.Stdin)
	conn.(*net.TCPConn).CloseWrite()
	<-done // wait for background goroutine to finish
}

// mustCopy 拿到连接数据
func mustCopy(dst io.Writer, src io.Reader) {
	nr, err := io.Copy(dst, src)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(nr)
}
