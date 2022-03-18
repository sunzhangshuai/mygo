package chat

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	TimeOut = 10 * time.Second
)

type client struct {
	name string
	ch   chan<- string
}

var (
	entering = make(chan *client)
	leaving  = make(chan *client)
	message  = make(chan string, 2)
)

// StartServer 开启服务端
func StartServer(port string) {
	var listener net.Listener
	var err error
	if listener, err = net.Listen("tcp", "localhost:"+port); err != nil {
		log.Print("服务故障，启动失败")
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go func() {
			handlerConn(conn)
			defer conn.Close()
		}()
	}
}

// handlerConn 处理链接
func handlerConn(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	var name string

	// 获取名称
	for scanner.Scan() {
		name = scanner.Text()
		break
	}

	ch := make(chan string, 10)
	cli := &client{
		name: name,
		ch:   ch,
	}

	// 将收到的消息发送客户端
	go func() {
		for msg := range ch {
			_, err := fmt.Fprintln(conn, msg)
			if err != nil {
				close(ch)
				return
			}
		}
	}()

	// 进入
	ch <- "You are " + cli.name
	entering <- cli

	active := make(chan struct{})

	// 遍历器
	go func() {
		for scanner.Scan() {
			active <- struct{}{}
			message <- cli.name + "：" + scanner.Text()
		}
	}()

	timer := time.NewTimer(TimeOut)

loop:
	for {
		select {
		case <-timer.C:
			conn.(*net.TCPConn).CloseWrite()
			break loop
		case <-active:
			// 时间重置了
			timer.Reset(TimeOut)
		}
	}

	// 离开
	leaving <- cli
}

// broadcaster 广播员
func broadcaster() {
	clients := make(map[*client]bool)
	for {
		select {
		// 进入
		case cli := <-entering:
			clients[cli] = true
			message <- cli.name + " has arrived"
		// 离开
		case cli := <-leaving:
			delete(clients, cli)
			message <- cli.name + " has left"
			close(cli.ch)
		case msg := <-message:
			for cli := range clients {
				cli.ch <- msg
			}
		}
	}
}
