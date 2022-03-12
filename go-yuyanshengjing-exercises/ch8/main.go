package ch8

import (
	"exercises/ch8/chat"
	"exercises/ch8/clock"
	"exercises/ch8/echo"
	"exercises/ch8/ftp/client"
	"exercises/ch8/ftp/server"
	"exercises/ch8/mirror"
	"fmt"
	"gopl.io/ch5/links"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Exercises 练习
type Exercises struct {
}

// crawl 爬url
func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

// Task1 修改clock2来支持传入参数作为端口号，
// 然后写一个clockwall的程序，这个程序可以同时与多个clock服务器通信，
// 从多个服务器中读取时间，并且在一个表格中一次显示所有服务器传回的结果，类似于你在某些办公室里看到的时钟墙。
// 如果你有地理学上分布式的服务器可以用的话，让这些服务器跑在不同的机器上面；
// 或者在同一台机器上跑多个不同的实例，这些实例监听不同的端口，假装自己在不同的时区。
// go run main.go -ch 8 -task 1 clock 8010 US/Eastern
// go run main.go -ch 8 -task 1 clock 8020 Asia/Tokyo
// go run main.go -ch 8 -task 1 clock 8030 Europe/London
// go run main.go -ch 8 -task 1 clockWall NewYork=8010 Tokyo=8020 London=8030
func (e *Exercises) Task1() {
	opType := os.Args[5]
	switch opType {
	case "clock":
		task1.Clock(os.Args[6], os.Args[7])
	case "clockWall":
		task1.ClockWall()
	}
}

// Task2 实现一个并发FTP服务器。服务器应该解析客户端发来的一些命令，
// 比如cd命令来切换目录，ls来列出目录内文件，get和send来传输文件，close来关闭连接。
// 你可以用标准的ftp命令来作为客户端，或者也可以自己实现一个。
// go run main.go -ch 8 -task 2 client 8010
// go run main.go -ch 8 -task 2 server 8010
func (e *Exercises) Task2() {
	var err error
	opType := os.Args[5]
	switch opType {
	case "server":
		err = server.Run(os.Args[6])
	case "client":
		err = client.Run(os.Args[6])
	}
	if err != nil {
		log.Println(err)
	}
}

// Task3
// 1. 在netcat3例子中，conn虽然是一个interface类型的值，但是其底层真实类型是*net.TCPConn，代表一个TCP连接。 一个TCP连接有读和写两个部分，可以使用CloseRead和CloseWrite方法分别关闭它们。 修改netcat3的主goroutine代码，只关闭网络连接中写的部分，这样的话后台goroutine可以在标准输入被关闭后继续打印从reverb1服务器传回的数据。 （要在reverb2服务器也完成同样的功能是比较困难的）
// 2. 修改reverb2服务器，在每一个连接中使用sync.WaitGroup来计数活跃的echo goroutine。 当计数减为零时，关闭TCP连接的写入，像练习8.3中一样。验证一下你的修改版netcat3客户端会一直等待所有的并发“喊叫”完成， 即使是在标准输入流已经关闭的情况下。
// 3. 使用select来改造8.3节中的echo服务器，为其增加超时，这样服务器可以在客户端10秒中没有任何喊话时自动断开连接。
// go run main.go -ch 8 -task 3 netcat 8010
// go run main.go -ch 8 -task 3 reverb 8010
func (e *Exercises) Task3() {
	opType := os.Args[5]
	switch opType {
	case "netcat":
		echo.Netcat(os.Args[6])
	case "reverb":
		echo.Reverb(os.Args[6])
	}
}

// Task4
// go run main.go -ch 8 -task 4 reverb 8010
// go run main.go -ch 8 -task 4 netcat 8010
func (e *Exercises) Task4() {
	e.Task3()
}

// Task5 使用一个已有的CPU绑定的顺序程序，比如在3.3节中我们写的Mandelbrot程序或者3.2节中的3-D surface计算程序，
// 并将他们的主循环改为并发形式，使用channel来进行通信。
// 在多核计算机上这个程序得到了多少速度上的改进？使用多少个goroutine是最合适的呢？
// go run main.go -ch 8 -task 5 zhangshuai
func (e *Exercises) Task5() {
	_, mainName, _, _ := runtime.Caller(0)
	filePath := path.Dir(mainName)
	fileName := path.Join(filePath, os.Args[5]+".svg")
	fmt.Println(fileName)
	file, err := os.Create(fileName)
	if err != nil {
		if os.IsExist(err) {
			file.Truncate(0)
		}
	}
	defer file.Close()
	surface(file)
}

// Task6WorkUrl ...
type Task6WorkUrl struct {
	depth int
	data  string
}

// Task6 为并发爬虫增加深度限制。
// 也就是说，如果用户设置了depth=1，那么只有从首页跳转三次以内能够跳到的页面才能被抓取到。
// go run main.go -ch 8 -task 6 2 http://gopl.io/ https://golang.org/help/ https://golang.org/doc/ https://golang.org/blog/
func (e *Exercises) Task6() {
	depth, _ := strconv.Atoi(os.Args[5])

	urls := make([]Task6WorkUrl, len(os.Args[6:]))
	// 初始化url
	for i, v := range os.Args[6:] {
		urls[i] = Task6WorkUrl{0, v}
	}
	workList := make(chan []Task6WorkUrl)  // URL列表，可能有重复项
	unseenLinks := make(chan Task6WorkUrl) // 消除重复的URL

	// 将命令行参数添加到工作列表
	go func() { workList <- urls }()

	// 创建20个爬虫goroutines来获取没有遍历到的链接
	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link.data)
				urls := make([]Task6WorkUrl, len(foundLinks))
				for i, v := range foundLinks {
					urls[i] = Task6WorkUrl{link.depth + 1, v}
				}
				// 专有的goroutine中发送到worklist中来避免死锁
				go func() { workList <- urls }()
			}
		}()
	}

	// 主goroutine消除重复工作列表项
	// 然后把没爬到的传送给爬虫。
	seen := make(map[string]bool)
	for list := range workList {
		for _, link := range list {
			if !seen[link.data] && link.depth < depth {
				seen[link.data] = true
				unseenLinks <- link
			}
		}
	}
}

// Task7 完成一个并发程序来创建一个线上网站的本地镜像，把该站点的所有可达的页面都抓取到本地硬盘。
// 为了省事，我们这里可以只取出现在该域下的所有页面（比如golang.org开头，译注：外链的应该就不算了。）
// 当然了，出现在页面里的链接你也需要进行一些处理，使其能够在你的镜像站点上进行跳转，而不是指向原始的链接。
// go run main.go -ch 8 -task 7
func (e *Exercises) Task7() {
	mirror.Run(2, "http://gopl.io/")
}

// Task8
// go run main.go -ch 8 -task 8 reverb 8010
// go run main.go -ch 8 -task 8 netcat 8010
func (e *Exercises) Task8() {
	e.Task3()
}

// Task9 编写一个du工具，每隔一段时间将root目录下的目录大小计算并显示出来。
// go run main.go -ch 8 -task 9
func (e *Exercises) Task9() {
	runDu()
}

// Task10 HTTP请求可能会因http.Request结构体中Cancel channel的关闭而取消。
// 修改8.6节中的web crawler来支持取消http请求。
// （提示：http.Get并没有提供方便地定制一个请求的方法。你可以用http.NewRequest来取而代之，设置它的Cancel字段，然后用http.DefaultClient.Do(req)来进行这个http请求。）
// go run main.go -ch 8 -task 10
func (e *Exercises) Task10() {
	var resp *http.Response
	var req *http.Request
	var err error

	if req, err = http.NewRequest(http.MethodGet, "http://gopl.io/", nil); err != nil {
		fmt.Println(err)
		return
	}

	cancelChan := make(chan struct{})
	doneChan := make(chan struct{}, 1)
	req.Cancel = cancelChan

	go func() {
		resp, err = http.DefaultClient.Do(req)
		doneChan <- struct{}{}
	}()

	tick := time.Tick(10 * time.Second)

	select {
	case <-doneChan:
		close(cancelChan)
		close(doneChan)
		if err != nil {
			io.Copy(os.Stdout, resp.Body)
			defer resp.Body.Close()
		}
	case <-tick:
		close(cancelChan)
		fmt.Println("req timeout")
		return
	}

	if err != nil {
		defer resp.Body.Close()
	}
}

// Task11 紧接着8.4.4中的mirroredQuery流程，实现一个并发请求url的fetch的变种。当第一个请求返回时，直接取消其它的请求。
// go run main.go -ch 8 -task 11
func (e *Exercises) Task11() {
	urls := []string{
		"http://gopl.io/",
		"http://golang.org/help/",
		"http://golang.org/doc/",
		"http://golang.org/blog/",
	}

	wg := sync.WaitGroup{}

	cancelChan := make(chan struct{})
	doneChan := make(chan struct{}, 1)

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			var resp *http.Response
			var req *http.Request
			var err error
			if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
				fmt.Println("fail：", err)
				return
			}
			req.Cancel = cancelChan
			if resp, err = http.DefaultClient.Do(req); err != nil {
				fmt.Println("fail：", err)
				return
			}
			resp.Body.Close()
			fmt.Println("success")
			doneChan <- struct{}{}
		}(url)
	}

	// 拿到done
	<-doneChan

	close(cancelChan)

	wg.Wait()
}

// Task12
// 1. 使broadcaster能够将arrival事件通知当前所有的客户端。这需要你在clients集合中，以及entering和leaving的channel中记录客户端的名字。
// 2. 使聊天服务器能够断开空闲的客户端连接，比如最近五分钟之后没有发送任何消息的那些客户端。提示：可以在其它goroutine中调用conn.Close()来解除Read调用，就像input.Scanner()所做的那样。
// 3. 修改聊天服务器的网络协议，这样每一个客户端就可以在entering时提供他们的名字。将消息前缀由之前的网络地址改为这个名字。
// 4. 如果一个客户端没有及时地读取数据可能会导致所有的客户端被阻塞。修改broadcaster来跳过一条消息，而不是等待这个客户端一直到其准备好读写。 或者为每一个客户端的消息发送channel建立缓冲区，这样大部分的消息便不会被丢掉；broadcaster应该用一个非阻塞的send向这个channel中发消息。
// go run main.go -ch 8 -task 12 server 8080
// go run main.go -ch 8 -task 12 client 8080 zs
func (e *Exercises) Task12() {
	opType := os.Args[5]
	switch opType {
	case "server":
		chat.StartServer(os.Args[6])
	case "client":
		chat.StartClient(os.Args[6], os.Args[7])
	}
}
