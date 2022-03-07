package client

import (
	"exercises/ch8/ftp"
	"fmt"
	"os"
	"path"
)

type FtpClient struct {
	ftp.FtpConn
}

// sendGetCommand 发送get命令
func (f *FtpClient) sendGetCommand(args []string) error {
	return f.sendCommandCommon(ftp.Commands["get"], args)
}

// receiveGetData 接收get数据
func (f *FtpClient) receiveGetData(args []string) error {
	if _, err := f.receiveDataCommon(); err != nil {
		return err
	}
	return f.saveFile(args[0])
}

// sendPutCommand 发送put命令
func (f *FtpClient) sendPutCommand(args []string) error {
	var file *os.File
	var err error
	var filePath string

	// 获取本地文件句柄
	if file, err = f.getFile(args[0]); err != nil {
		return err
	}
	defer file.Close()

	// 发送命令
	if err = f.sendCommandCommon(ftp.Commands["put"], []string{path.Base(filePath)}); err != nil {
		return err
	}

	// 发送文件
	if err = f.sendFile(file); err != nil {
		return err
	}
	return nil
}

// receivePutData 接收put数据
func (f *FtpClient) receivePutData() error {
	if _, err := f.receiveDataCommon(); err != nil {
		return err
	}
	return nil
}

// sendMkdirCommand 发送mkdir命令
func (f *FtpClient) sendMkdirCommand(args []string) error {
	return f.sendCommandCommon(ftp.Commands["mkdir"], args)
}

// receiveMkdirData 接收mkdir数据
func (f *FtpClient) receiveMkdirData() error {
	if _, err := f.receiveDataCommon(); err != nil {
		return err
	}
	return nil
}

// sendCdCommand 发送Cd命令
func (f *FtpClient) sendCdCommand(args []string) error {
	return f.sendCommandCommon(ftp.Commands["cd"], args)
}

// receiveCdData 接收cd数据
func (f *FtpClient) receiveCdData() error {
	var res []byte
	var err error

	if res, err = f.receiveDataCommon(); err != nil {
		return err
	}
	f.ServerCwd = ftp.Sbyte2Str(res)
	fmt.Printf("cd %s：\n", f.ServerCwd)
	return nil
}

// sendLsCommand 发送ls命令
func (f *FtpClient) sendLsCommand(args []string) error {
	return f.sendCommandCommon(ftp.Commands["ls"], args)
}

// receiveLsData 接收Ls数据
func (f *FtpClient) receiveLsData() error {
	if _, err := f.receiveDataCommon(); err != nil {
		return err
	}
	return nil
}

// sendExitCommand 发送exit命令
func (f *FtpClient) sendExitCommand(args []string) error {
	return f.sendCommandCommon(ftp.Commands["exit"], args)
}

// receiveExitData 接收exit数据
func (f *FtpClient) receiveExitData() error {
	if _, err := f.receiveDataCommon(); err != nil {
		return err
	}
	f.Exit = true
	return nil
}
