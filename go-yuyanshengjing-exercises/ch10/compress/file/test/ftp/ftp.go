package ftp

import (
	"net"
	"unsafe"
)

const (
	// ServerHomeDir 服务端家目录
	ServerHomeDir = "/Users/zhangshuai24/Documents/notes/语言/编译型语言/go/Go程序设计语言/exercises/test/ftp/file/server"
	// ClientHomeDir 客户端家目录
	ClientHomeDir = "/Users/zhangshuai24/Documents/notes/语言/编译型语言/go/Go程序设计语言/exercises/test/ftp/file/client"
)

// Commands 命令行配置
var Commands = map[string]uint8{
	"cd":    uint8(1),
	"ls":    uint8(2),
	"exit":  uint8(3),
	"mkdir": uint8(4),
	"put":   uint8(5),
	"get":   uint8(6),
}

// FtpConn ftp 连接
type FtpConn struct {
	Con        net.Conn
	ServerCwd  string
	ServerHome string
	ClientCwd  string
	ClientHome string
	Exit       bool
}

// Str2sbyte 字符串 -> byte数组
func Str2sbyte(s string) (b []byte) {
	*(*string)(unsafe.Pointer(&b)) = s
	*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&b)) + 2*unsafe.Sizeof(&b))) = len(s)
	return
}

// Sbyte2Str byte数组 -> 字符串
func Sbyte2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
