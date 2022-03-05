package core

import (
	pb "chat_server/src/pb"
	"github.com/gogo/protobuf/proto"
	"net"
	"sync"
	"time"
)

type User struct {
	conn  net.Conn
	mutex sync.Mutex

	UserName      string `json:"user_name"`
	RoomId        int64  `json:"room_id"`
	LastLoginTime int64  `json:"last_login_time"` // 时间戳,单位s
	OnlineTime    int64  `json:"online_time"`     // 单位s
	IsOnline      bool   `json:"is_online"`       // 是否在线
}

// 登录
func (u *User) Login(b []byte) {
	req := new(pb.LoginReq)
	proto.Unmarshal(b, req)

	if len(u.UserName) != 0 {
		u.Send("不要重复登录")
		return
	}
	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.UserName = req.UserName

	// 查
	userData, ok := UserDB[req.UserName]

	// 新增用户
	if !ok {
		userData = User{
			UserName:      req.UserName,
			LastLoginTime: time.Now().Unix(),
			RoomId:        0,
			OnlineTime:    0,
			IsOnline:      true,
		}

		// 修改用户
	} else {
		u.RoomId = 0
		u.LastLoginTime = time.Now().Unix()
		u.IsOnline = true

		userData.RoomId = 0
		userData.IsOnline = true
		userData.LastLoginTime = time.Now().Unix()
	}

	// 保存数据库
	UserDB[req.UserName] = userData
	u.Send("登陆成功")
}

// 登出
func (u *User) Logout() {
	u.mutex.Lock()
	u.mutex.Unlock()

	if len(u.UserName) == 0 {
		return
	}

	useData := UserDB[u.UserName]
	useData.IsOnline = false

	// 在线时长
	onlineTime := time.Now().Unix() - useData.LastLoginTime
	if onlineTime > 0 {
		useData.OnlineTime += onlineTime
	}

	// 踢出房间

	// 保存
	UserDB[u.UserName] = useData
}

// 加入房间
func (u *User) JoinRoom(b []byte) {

}

// 发送消息
func (u *User) SendMessage() {

}

func (u *User) Send(str string) {
	u.conn.Write([]byte(str))
}
