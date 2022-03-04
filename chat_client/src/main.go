package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("dial err: ", err)
		return
	}
	defer conn.Close()

	//for {
	s := "what are you 弄啥？"
	_, err = conn.Write([]byte(s))
	if err != nil {
		fmt.Println("writ err:", err)
		return
	}
	time.Sleep(time.Second)
	//}

}
