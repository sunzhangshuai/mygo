package client

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"exercises/ch8/ftp"
)

// sendCommandCommon 发送命令通用方法
func (f *FtpClient) sendCommandCommon(cmdID uint8, args []string) error {
	var err error
	var length uint64

	// 发送格式 长度 + 命令id + 内容
	// 1. 发送总长度
	argStr := strings.Join(args, " ")
	length = uint64(binary.Size(length)+binary.Size(cmdID)) + uint64(len(argStr))
	if err = binary.Write(f.Con, binary.LittleEndian, length); err != nil {
		return err
	}

	// 2. 发送指令
	if err = binary.Write(f.Con, binary.LittleEndian, cmdID); err != nil {
		return err
	}

	// 3. 发送内容
	if err = binary.Write(f.Con, binary.LittleEndian, ftp.Str2sbyte(argStr)); err != nil {
		return err
	}
	return nil
}

// receiveDataCommon 接收数据通用方法
func (f *FtpClient) receiveDataCommon() ([]byte, error) {
	var length uint64
	var err error

	// 1. 获取结果长度，长度为0结束
	if err = binary.Read(f.Con, binary.LittleEndian, &length); err != nil {
		return nil, err
	}
	if length == 0 {
		fmt.Printf("%s：\n", f.ServerCwd)
		return nil, nil
	}

	// 2. 接收数据
	res := make([]byte, length-uint64(binary.Size(length)))
	if err = binary.Read(f.Con, binary.LittleEndian, res); err != nil {
		return nil, err
	}
	fmt.Printf("res：%s：\n", res)
	return res, nil
}

// sendFile 发送文件
func (f *FtpClient) sendFile(file *os.File) error {
	var fileInfo os.FileInfo
	var err error

	// 发送文件长度
	if fileInfo, err = file.Stat(); err != nil {
		return err
	}
	if fileInfo.IsDir() {
		return errors.New("不支持 put 文件夹")
	}
	if err = binary.Write(f.Con, binary.LittleEndian, fileInfo.Size()); err != nil {
		return err
	}

	// 发送文件内容
	buf := make([]byte, 4096)
	bufReader := bufio.NewReader(file)
	for {
		var nr int
		if nr, err = bufReader.Read(buf); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if err = binary.Write(f.Con, binary.LittleEndian, buf[:nr]); err != nil {
			return err
		}
	}
	return nil
}

// getFile 获取文件
func (f *FtpClient) getFile(filePath string) (*os.File, error) {
	var file *os.File
	var err error

	filePath = strings.Replace(filePath, "\\", "/", -1)
	if strings.HasPrefix(filePath, "/") {
		filePath = path.Join(f.ClientHome, filePath)
	} else {
		filePath = path.Join(f.ClientCwd, filePath)
	}
	if file, err = os.Open(filePath); err != nil {
		return nil, err
	}
	return file, nil
}

// createFile 创建文件
func (f *FtpClient) createFile(filePath string) (*os.File, error) {
	var file *os.File
	var err error

	fileName := path.Base(filePath)
	if file, err = os.Create(path.Join(f.ClientCwd, fileName)); err != nil {
		if os.IsExist(err) {
			if err = file.Truncate(0); err != nil {
				return nil, err
			}
		}
		return nil, err
	}
	return file, nil
}

// saveFile 保存文件
func (f *FtpClient) saveFile(filePath string) error {
	var file *os.File
	var err error

	// 创建文件
	if file, err = f.createFile(filePath); err != nil {
		return err
	}
	defer file.Close()

	// 保存文件内容
	var length uint64
	var total, bufSize uint64
	var nr int
	if err = binary.Read(f.Con, binary.LittleEndian, &length); err != nil {
		return err
	}
	bufSize = length
	if length > 4096 {
		bufSize = 4096
	}
	buf := make([]byte, bufSize)
	for total < length {
		if err = binary.Read(f.Con, binary.LittleEndian, buf); err != nil {
			return err
		}
		if nr, err = file.Write(buf); err != nil {
			return err
		}
		total += uint64(nr)
		if length-total < bufSize {
			buf = buf[0 : length-total]
		}
	}
	return err
}
