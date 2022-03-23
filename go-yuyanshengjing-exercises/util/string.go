package util

import (
	"bytes"
	"io"
)

// ReadString 读取字符串
func ReadString(r io.Reader) string {
	buf := bytes.Buffer{}
	buf.ReadFrom(r) // notice：忽略错误
	return buf.String()
}
