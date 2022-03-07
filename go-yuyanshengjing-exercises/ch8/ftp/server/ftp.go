package server

import (
	"exercises/ch8/ftp"
	"fmt"
	"os"
)

// FtpServer ftp服务端
type FtpServer struct {
	ftp.FtpConn
}

// handleMkdir 处理mkdir
func (f *FtpServer) handleMkdir(args []byte) error {
	dir := f.getPath(ftp.Sbyte2Str(args))
	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		return err
	}
	return f.SendContent([]byte("mkdir ok"))
}

// handleLs 处理ls
func (f *FtpServer) handleLs(args []byte) error {
	var dir *os.File
	var dirInfo os.FileInfo
	var dirList []os.FileInfo
	var err error
	var res string

	pathStr := f.getPath(ftp.Sbyte2Str(args))

	if dir, err = os.Open(pathStr); err != nil {
		if err = f.SendContent(ftp.Str2sbyte(err.Error())); err != nil {
			return err
		}
		return nil
	}
	if dirInfo, err = dir.Stat(); err != nil {
		if err = f.SendContent(ftp.Str2sbyte(err.Error())); err != nil {
			return err
		}
		return nil
	}
	if !dirInfo.IsDir() {
		res = fmt.Sprintf("%.30s\t%.10d\t%s\n", dirInfo.Name(), dirInfo.Size(), dirInfo.ModTime())
		if err = f.SendContent(ftp.Str2sbyte(res)); err != nil {
			return err
		}
		return nil
	}
	if dirList, err = dir.Readdir(0); err != nil {
		if err = f.SendContent(ftp.Str2sbyte(res)); err != nil {
			return err
		}
		return nil
	}
	res = fmt.Sprintf("Total:%d\n", len(dirList))
	for _, entry := range dirList {
		res += fmt.Sprintf("%.30s\t%.10d\t%s\n", entry.Name(), entry.Size(), entry.ModTime())
	}
	if err = f.SendContent(ftp.Str2sbyte(res)); err != nil {
		return err
	}
	return nil
}
