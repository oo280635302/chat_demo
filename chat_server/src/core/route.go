package core

import (
	pb "chat_server/src/pb"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"io"
	"net"
	"sync"
)

type TCPConn struct {
	Lock     sync.Mutex
	Listener net.Listener

	Users map[string]*User
}

// TCP 处理
func (t *TCPConn) Listen(address string) error {
	t.Users = make(map[string]*User)
	l, err := net.Listen("tcp", address)
	if err == nil {
		t.Listener = l
	}
	fmt.Println("监听：", address)
	return err
}

func (t *TCPConn) End() {
	t.Listener.Close()
}

func (t *TCPConn) Start() {
	go func() {
		for {
			conn, err := t.Listener.Accept()
			if err != nil {
				fmt.Println("ERROR: accept err:", err)
			} else {
				client := t.accept(conn)
				go t.serve(client)
			}
		}
	}()
}

func (t *TCPConn) serve(user *User) {
	defer t.quiet(user)
	for {
		buf := make([]byte, 4096)
		n, err := user.conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println("ERROR: route read err:", err)
			return
		}

		if n == 0 {
			continue
		}

		fmt.Println(string(buf[:n]))

		req := &pb.Request{}
		err = proto.Unmarshal(buf[:n], req)
		if err != nil {
			fmt.Println("ERROR: route unmarshal err:", err)

		}

		HandleAPI(req.Api, user, req.Data)
	}
}

func (t *TCPConn) accept(conn net.Conn) *User {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	fmt.Println("新连接来了~", conn.RemoteAddr().String(), len(t.Users)+1)
	user := &User{
		conn: conn,
	}
	t.Users[conn.RemoteAddr().String()] = user
	return user
}

func (t *TCPConn) quiet(user *User) {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	// 退出保存数据库
	user.Logout()

	// 清理连接
	delete(t.Users, user.conn.RemoteAddr().String())

	fmt.Println("溜了~ ", user.conn.RemoteAddr().String())
	user.conn.Close()
}

func HandleAPI(api pb.API, user *User, b []byte) {
	switch api {
	case pb.API_Login:
		user.Login(b)

		// todo
	}
}
