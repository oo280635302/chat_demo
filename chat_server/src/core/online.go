package core

import (
	pb "chat_server/src/pb"
	"github.com/gogo/protobuf/proto"
	"net"
	"time"
)

type Client struct {
	Conn          net.Conn
	UserName      string `json:"user_name"`
	LastLoginTime int64  `json:"last_login_time"`
	RoomId        string `json:"room_id"`
}

func Login(conn net.Conn, b []byte) {
	req := new(pb.LoginReq)
	proto.Unmarshal(b, req)

	onlineMap[req.UserName] = &Client{
		Conn:          conn,
		UserName:      time.Now().String(),
		LastLoginTime: time.Now().UnixNano(),
	}

}
