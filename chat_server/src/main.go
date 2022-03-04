package main

import (
	"chat_server/src/common"
	"chat_server/src/conn"
	"chat_server/src/core"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// 工具
	common.Init()

	// 核心
	core.Init()

	fmt.Println("init success")
	// 连接
	l, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("ERROR: Listen tcp 8888 err:", err)
		return
	}
	defer l.Close()
	go conn.Listen(l)
	fmt.Println("listen 8888...")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGINT,
		syscall.SIGILL,
		syscall.SIGFPE,
		syscall.SIGSEGV,
		syscall.SIGTERM,
		syscall.SIGABRT)
	<-signalChan
}
