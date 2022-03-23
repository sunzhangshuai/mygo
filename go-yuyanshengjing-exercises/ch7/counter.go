package ch7

import (
	"bufio"
	"io"
)

// MyCounter 数量统计
type MyCounter struct {
	words, lines int
}

// CountWriter 计数器
type CountWriter struct {
	cw io.Writer
	C  int64
}

// Write 写入
func (c *MyCounter) Write(p []byte) (int, int) {
	var words, lines int

	for i := 0; i < len(p); {
		advance, _, _ := bufio.ScanWords(p, true)
		words++
		i += advance
	}

	for i := 0; i < len(p); {
		advance, _, _ := bufio.ScanLines(p, true)
		lines++
		i += advance
	}
	c.words += words
	c.lines += lines
	return words, lines
}

// Write 计数器写方法
func (c *CountWriter) Write(p []byte) (int, error) {
	c.C = int64(len(p))
	n, err := c.cw.Write(p)
	return n, err
}

// CountingWriter io中间件
func CountingWriter(w io.Writer) (io.Writer, *int64) {
	c := &CountWriter{cw: w}
	return c, &c.C
}
