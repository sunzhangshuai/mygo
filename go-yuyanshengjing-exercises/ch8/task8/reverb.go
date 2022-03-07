package task8

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

	scan := make(chan struct{})

	// 等待结束
	wg := sync.WaitGroup{}

	go func() {
		defer wg.Done()
		for {
			select {
			case <-scan:
				// 继续循环
			case <-time.After(10 * time.Second):
				fmt.Println("停掉了")
				c.Close()
				return
			}
		}
	}()

	for input.Scan() {
		scan <- struct{}{}
		wg.Add(1)
		go echo(c, input.Text(), 1*time.Second)
	}
	// NOTE: ignoring potential errors from input.Err()
	c.(*net.TCPConn).CloseRead()

	wg.Wait()
}

// echo 输出
func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}
