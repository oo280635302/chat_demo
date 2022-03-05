package main

import (
	"bufio"
	"chat_server/src/common"
	"chat_server/src/core"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	s.Start()
	defer s.End()

	// GM
	for {
		reader := bufio.NewReader(os.Stdin)
		line, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println("ERROR: os stdin err:", err)
			break
		}
		if len(line) == 0 {
			continue
		}
		if line[0] != '/' {
			continue
		}
		strs := strings.Split(string(line), " ")
		if len(strs) != 2 {
			continue
		}
		switch strs[0] {
		case "/stats":
			core.GetUserInfo(strs[1])
		case "/popular":
			id, _ := strconv.Atoi(strs[1])
			if id <= 0 {
				fmt.Println("错误的房间号: ", id)
				continue
			}
			core.GetPopular(int64(id))
		default:
			continue
		}
	}

}
