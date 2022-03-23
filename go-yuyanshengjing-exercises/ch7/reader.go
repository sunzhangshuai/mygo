package ch7

import (
	"golang.org/x/net/html"
	"io"

	"exercises/util"
)

// MyReader my reader
type MyReader struct {
	data   string
	curLen int64
	maxLen int64
}

// NewReader r
func NewReader(s string) *MyReader {
	return &MyReader{
		data:   s,
		maxLen: int64(len(s)),
	}
}

// HtmlParse HTML解析器
func HtmlParse(s string) (*html.Node, error) {
	r := NewReader(s)
	return html.Parse(r)
}

// Read 读取
func (r *MyReader) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	n = copy(b, r.data[r.curLen:r.maxLen])
	if r.curLen += int64(n); r.curLen == r.maxLen {
		err = io.EOF
	}
	return
}

// LimitReader 置顶最多的读取数量
func LimitReader(r io.Reader, n int64) io.Reader {
	return &MyReader{data: util.ReadString(r), maxLen: n}
}
