package client

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"exercises/ch8/ftp"
)

// Run 启动客户端
func Run(port string) error {
	var err error
	var client *FtpClient

	// 初始化客户端
	if client, err = initClient(port); err != nil {
		return err
	}

	// 阻塞等待命令
	blockWaitCommand(client)
	return nil
}

// initClient 初始化客户端
func initClient(port string) (*FtpClient, error) {
	var conn net.Conn
	var err error

	// 建立连接
	if conn, err = net.Dial("tcp", "localhost:"+port); err != nil {
		return nil, err
	}

	return &FtpClient{
		ftp.FtpConn{
			Con:        conn,
			ServerCwd:  ftp.ServerHomeDir,
			ServerHome: ftp.ServerHomeDir,
			ClientCwd:  ftp.ClientHomeDir,
			ClientHome: ftp.ClientHomeDir,
		},
	}, nil
}

// blockWaitCommand 等待命令
func blockWaitCommand(client *FtpClient) {
	// 初始化输入
	input := bufio.NewScanner(os.Stdin)

	// 等待输入
	for input.Scan() && !client.Exit {
		// 处理输入内容
		argStr := input.Text()

		arg := strings.Split(strings.TrimSpace(argStr), " ")
		if len(arg) == 0 {
			printHelp()
		}

		// 获取命令&参数
		command := arg[0]
		if len(arg) > 1 {
			arg = arg[1:]
		} else {
			arg = nil
		}

		// 执行命令
		if err := handleCommand(client, command, arg); err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err)
		}
	}
}

// printHelp 打印帮助信息
func printHelp() {
	log.Println("Help:\t[command][args]\ncp [path]")
}

// handleCommand 处理命令
func handleCommand(client *FtpClient, command string, args []string) error {
	var cmdID uint8
	var err error
	var ok bool

	if cmdID, ok = ftp.Commands[command]; !ok {
		return errors.New("unsupported command")
	}

	// 发送命令
	if err = sendCommand(client, cmdID, args); err != nil {
		return err
	}

	// 接收返回值
	if err = receiveData(client, cmdID, args); err != nil {
		return err
	}
	return nil
}

// sendCommand 发送指令
func sendCommand(client *FtpClient, cmdID uint8, args []string) error {
	var err error
	switch cmdID {
	case ftp.Commands["get"]:
		err = client.sendGetCommand(args)
	case ftp.Commands["put"]:
		err = client.sendPutCommand(args)
	case ftp.Commands["mkdir"]:
		err = client.sendMkdirCommand(args)
	case ftp.Commands["ls"]:
		err = client.sendLsCommand(args)
	case ftp.Commands["cd"]:
		err = client.sendMkdirCommand(args)
	}
	return err
}

// receiveData 接收数据
func receiveData(client *FtpClient, cmdID uint8, args []string) error {
	var err error

	switch cmdID {
	case ftp.Commands["get"]:
		err = client.receiveGetData(args)
	case ftp.Commands["put"]:
		err = client.receivePutData()
	case ftp.Commands["mkdir"]:
		err = client.receiveMkdirData()
	case ftp.Commands["ls"]:
		err = client.receiveLsData()
	case ftp.Commands["exit"]:
		err = client.receiveMkdirData()
	case ftp.Commands["cd"]:
		err = client.sendCdCommand(args)
	}
	return err
}
