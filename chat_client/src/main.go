package main

import (
	"bufio"
	pb "chat_client/src/pb"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("dial err: ", err)
		return
	}
	defer conn.Close()

	go func() {
		for {
			buf := make([]byte, 4096)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("ERROR: read err:", err)
				conn.Close()
				os.Exit(0)
			}
			if n == 0 {
				continue
			}
			fmt.Println(string(buf[:n]))
		}
	}()

	fmt.Println("请输入~ 以:为分隔符")
	for {

		reader := bufio.NewReader(os.Stdin)
		res, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println("ERROR: reader down err:", err)
			return
		}

		// 正常流程
		strs := strings.Split(string(res), ":")
		if len(strs) != 2 {
			fmt.Println("WARN:无效的输入!!!", len(strs), strs)
			continue
		}

		api := pb.API(pb.API_value[strs[0]])

		switch api {
		case pb.API_Login:
			req := &pb.LoginReq{UserName: strs[1]}
			b, _ := proto.Marshal(req)
			request := &pb.Request{
				Api:  pb.API_Login,
				Data: b,
			}
			b, _ = proto.Marshal(request)
			_, err := conn.Write(b)
			if err != nil {
				fmt.Println("write err:", err)
				return
			}
		case pb.API_JoinRoom:
			roomId, _ := strconv.Atoi(strs[1])

			req := &pb.JoinRoomReq{RoomId: int64(roomId)}
			b, _ := proto.Marshal(req)
			request := &pb.Request{
				Api:  pb.API_JoinRoom,
				Data: b,
			}
			b, _ = proto.Marshal(request)
			conn.Write(b)
		case pb.API_SendMessage:
			req := &pb.SendMessageReq{Msg: strs[1]}
			b, _ := proto.Marshal(req)
			request := &pb.Request{
				Api:  pb.API_SendMessage,
				Data: b,
			}
			b, _ = proto.Marshal(request)
			conn.Write(b)
		default:
			fmt.Println("Invalid Stdio")
		}
	}
}
