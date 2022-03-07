package task4

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// Reverb 执行
func Reverb(port string) {
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
		go handleConn(conn)
	}
}

// handleConn 处理连接
func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	wg := sync.WaitGroup{}
	for input.Scan() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			echo(c, input.Text(), 3*time.Second)
		}()
	}
	// NOTE: ignoring potential errors from input.Err()
	c.(*net.TCPConn).CloseRead()
	go func() {
		wg.Wait()
		fmt.Println("wg wait done")
		c.(*net.TCPConn).CloseWrite()
	}()
}

// echo 输出
func echo(c net.Conn, shout string, delay time.Duration) {
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}
