package conn

import (
	pb "chat_server/src/pb"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"net"
)

type TCPConn struct {
	Conns map[string]net.Conn
}

func Listen(listen net.Listener) {
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept err", err)
			return
		}

		go HandlerConnect(conn)
	}
}

func HandlerConnect(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr().String()
	fmt.Println("新连接来啦~", addr)
	for {
		buf := make([]byte, 8192)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("client addr exit:", err, addr)
			return
		}
		fmt.Println(string(buf))
		req := &pb.Request{}
		proto.Unmarshal(buf[:n], req)

		continue

		handle := handlers[req.Api]
		if handle == nil {
			fmt.Println("ERROR: client invalid API!")
			return
		}
		handle(conn, req.Data)
	}
}
