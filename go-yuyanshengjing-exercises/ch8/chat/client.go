package chat

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// StartClient 开启客户端
func StartClient(port string, name string) {
	client, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		log.Println(err)
	}

	remoteDone := make(chan struct{})
	localDone := make(chan struct{})

	// 发送数据
	go func() {
		// 先发送名称
		fmt.Fprintln(client, name)

		if _, err := io.Copy(client, os.Stdin); err != nil {
			log.Println(err)
		}
		log.Println("local done")
		localDone <- struct{}{}
	}()

	// 接收数据
	go func() {
		if _, err := io.Copy(os.Stdout, client); err != nil {
			log.Println(err)
		}
		log.Println("remote done")
		remoteDone <- struct{}{}
	}()

	for {
		select {
		case <-remoteDone:
			client.Close()
			return
		case <-localDone:
			client.(*net.TCPConn).CloseWrite()
		}
	}
}
