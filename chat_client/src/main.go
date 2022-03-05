package main

import (
	pb "chat_client/src/pb"
	"fmt"
	"github.com/gogo/protobuf/proto"
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

	go func() {
		req := &pb.LoginReq{UserName: "小三"}
		b, _ := proto.Marshal(req)
		request := &pb.Request{
			Api:  pb.API_Login,
			Data: b,
		}
		b, _ = proto.Marshal(request)
		_, err = conn.Write(b)
		if err != nil {
			fmt.Println("writ err:", err)
			return
		}
		time.Sleep(time.Second)
	}()

	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read err:", err)
			return
		}
		if n == 0 {
			continue
		}
		fmt.Println(string(buf[:n]))
	}

}
