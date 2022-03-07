package server

import (
	"encoding/binary"
	"path"
	"strings"
)

// SendContent 发送数据
func (f *FtpServer) SendContent(content []byte) error {
	var length uint64
	length = uint64(len(content))
	if length == 0 {
		return binary.Write(f.Con, binary.LittleEndian, &length)
	}
	length = length + uint64(binary.Size(length))
	if err := binary.Write(f.Con, binary.LittleEndian, &length); err != nil {
		return err
	}
	if err := binary.Write(f.Con, binary.LittleEndian, content); err != nil {
		return err
	}
	return nil
}

// getPath 获取路径
func (f *FtpServer) getPath(dir string) string {
	if strings.HasPrefix(dir, "/") {
		dir = path.Join(f.ServerHome, dir)
	} else {
		dir = path.Join(f.ServerCwd, dir)
	}
	return dir
}
