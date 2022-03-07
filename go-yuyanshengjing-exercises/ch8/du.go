package ch8

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

// baseDir 基础dir
var baseDir string

// runDu 编写一个du工具，每隔一段时间将root目录下的目录大小计算并显示出来。
func runDu() {
	// 获取到要计算的文件夹
	_, fileName, _, _ := runtime.Caller(0)
	baseDir = path.Dir(fileName)
	tick := time.Tick(10 * time.Second)
	for {
		select {
		case <-tick:
			go calculate()
		}
	}
}

// calculate 计算大小
func calculate() {
	fileSize := make(chan int64)
	totalSize := int64(0)
	totalCount := 0
	// 计算大小
	go func() {
		walkDir(baseDir, fileSize)
		close(fileSize)
	}()

	// 每隔10毫秒打印一次
	tick := time.Tick(1 * time.Millisecond)
loop:
	for {
		select {
		case <-tick:
			fmt.Printf("%d files  %.1f KB\n", totalCount, float64(totalSize)/1e3)
		case size, ok := <-fileSize:
			if !ok {
				fmt.Printf("%d files  %.1f KB\n", totalCount, float64(totalSize)/1e3)
				fmt.Println("========================================================")
				break loop
			}
			totalCount++
			totalSize += size
		}
	}
}

// walkDir 处理dir
func walkDir(dirName string, fileSize chan int64) {
	var list []os.DirEntry
	var err error
	if list, err = os.ReadDir(dirName); err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return
	}
	for _, entry := range list {
		info, _ := entry.Info()
		if entry.IsDir() {
			walkDir(path.Join(dirName, info.Name()), fileSize)
		} else {
			fileSize <- info.Size()
		}
	}
}
