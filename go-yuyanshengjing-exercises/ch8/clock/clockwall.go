package task1

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

// ClockWall 时间墙
func ClockWall(args []string) {
	// 先输出时区
	title := "时区\t"
	for _, val := range args {
		values := strings.Split(val, "=")
		title += values[0] + "\t"
		//go connTcp(values[1])
	}
	fmt.Println(title)

	// 开启异步收取
	valueArr := make([]string, len(args))
	for i, val := range args {
		values := strings.Split(val, "=")
		go connTcp(values[1], &valueArr[i])
	}

	// 刷新时间
	for {
		value := "时间\t"
		for i := range valueArr {
			for {
				if len(valueArr[i]) > 0 {
					value += valueArr[i] + "\t"
					break
				}
			}
		}
		fmt.Println(value)
		time.Sleep(1 * time.Second)
	}
}

// 链接tcp
func connTcp(port string, out *string) {
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(out, conn)
}

// 拿到连接数据
func mustCopy(out *string, src io.Reader) {
	buf := make([]byte, 100)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			*out = string(buf[:nr])
		}
		if er != nil {
			log.Fatal(er)
		}
	}
}
