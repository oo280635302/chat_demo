package conn

import (
	"chat_server/src/core"
	pb "chat_server/src/pb"
	"net"
)

var handlers = map[pb.API]func(net.Conn, []byte){
	pb.API_Login: core.Login,
}
