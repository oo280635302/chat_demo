package main

import (
	"chat_server/src/common"
	"chat_server/src/core"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// 工具
	common.Init()

	// 核心
	core.Init()

	fmt.Println("init success")

	// 连接
	var s core.TCPConn
	err := s.Listen("127.0.0.1:8888")
	if err != nil {
		fmt.Println("ERROR: Listen tcp 8888 err:", err)
		return
	}
	defer s.End()
	s.Start()

	for {
		time.Sleep(time.Second * 5)
		fmt.Println("当前连接：", len(s.Users), s.Users)
		fmt.Println("数据库数据：", core.UserDB)
	}

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
