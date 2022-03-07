package server

import (
	"encoding/binary"
	"exercises/ch8/ftp"
	"io"
	"log"
	"net"
)

// Run 服务端运行
func Run(port string) error {
	var listener net.Listener
	var con net.Conn
	var err error
	if listener, err = net.Listen("tcp", "localhost:"+port); err != nil {
		return err
	}

	for {
		// 阻塞等待连接
		if con, err = listener.Accept(); err != nil {
			log.Println(err)
			continue
		}

		server := initServer(con)
		go blockWaitCommand(server)
	}
}

// blockWaitCommand 阻塞等待命令
func blockWaitCommand(server *FtpServer) {
	var err error

	// 阻塞监听
	for !server.Exit {
		var cmdID uint8
		var args []byte

		// 获取命令
		if cmdID, args, err = getCommand(server); err != nil {
			if err == io.EOF {
				log.Println("exit")
				break
			}
			log.Println(err)
			continue
		}

		// 处理命令
		if err = handleCommand(server, cmdID, args); err != nil {
			log.Println(err)
			continue
		}
	}
}

// initServer 初始化服务端
func initServer(con net.Conn) *FtpServer {
	return &FtpServer{
		ftp.FtpConn{
			Con:        con,
			ServerHome: ftp.ServerHomeDir,
			ServerCwd:  ftp.ServerHomeDir,
		},
	}
}

// getCommand 获取命令
func getCommand(server *FtpServer) (uint8, []byte, error) {
	var length uint64
	var err error

	// 读取长度
	if err = binary.Read(server.Con, binary.LittleEndian, &length); err != nil {
		return 0, nil, err
	}

	// 读取命令
	var cmdID uint8
	if err = binary.Read(server.Con, binary.LittleEndian, &cmdID); err != nil {
		return 0, nil, err
	}
	// 读取指令
	args := make([]byte, length-uint64(binary.Size(length))-uint64(binary.Size(cmdID)))
	if err = binary.Read(server.Con, binary.LittleEndian, args); err != nil {
		return 0, nil, err
	}
	return cmdID, args, nil
}

// handleCommand 处理命令
func handleCommand(server *FtpServer, cmdID uint8, args []byte) error {
	var err error

	switch cmdID {
	case ftp.Commands["mkdir"]:
		err = server.handleMkdir(args)
	case ftp.Commands["ls"]:
		err = server.handleLs(args)
	}
	return err
}
