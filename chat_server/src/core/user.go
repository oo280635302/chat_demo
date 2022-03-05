package core

import (
	"chat_server/src/common"
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
		u.Reply("不要重复登录")
		return
	}
	u.mutex.Lock()
	defer u.mutex.Unlock()

	// 查
	user, ok := UserDB[req.UserName]

	tm := time.Now().Unix()

	// 新增用户
	if !ok {
		u.IsOnline = true
		u.LastLoginTime = tm

		// 修改用户
	} else {
		if user.IsOnline {
			u.Reply("该用户正在登陆")
			return
		}

		u.RoomId = 0
		u.LastLoginTime = tm
		u.IsOnline = true
	}
	u.UserName = req.UserName

	// 保存数据库
	u.Save()

	u.Reply("登陆成功")
}

// 登出
func (u *User) Logout() {
	u.mutex.Lock()
	u.mutex.Unlock()

	if len(u.UserName) == 0 {
		return
	}

	u.IsOnline = false
	// 在线时长
	onlineTime := time.Now().Unix() - u.LastLoginTime
	if onlineTime > 0 {
		u.OnlineTime += onlineTime
	}
	u.RoomId = 0

	// 保存
	u.Save()
}

// 加入房间
func (u *User) JoinRoom(b []byte) {
	req := new(pb.JoinRoomReq)
	proto.Unmarshal(b, req)

	if req.RoomId <= 0 {
		u.Reply("无效的房间号")
		return
	}

	u.mutex.Lock()
	defer u.mutex.Unlock()

	// 退出之前的房间
	if u.RoomId != 0 {
		RoomDB[u.RoomId].ExitRoom(u)
	}

	u.RoomId = req.RoomId

	roomData, ok := RoomDB[req.RoomId]
	// 没有就新建
	if !ok {
		roomData = CreateRoom(req.RoomId)
	}
	roomData.Users[u.UserName] = u

	// 加入房间
	roomData.JoinRoom(u)

	// 保存数据
	u.Save()
}

// 发送消息
func (u *User) SendMessage(b []byte) {
	req := new(pb.SendMessageReq)
	proto.Unmarshal(b, req)

	if len(u.UserName) == 0 {
		u.Reply("需要登录")
		return
	}

	if u.RoomId == 0 {
		u.Reply("需要在房间内才能聊天")
		return
	}

	room := RoomDB[u.RoomId]

	// 脏词过滤
	msg := common.Trie.Replace(req.Msg, '*')

	// 广播
	room.BroadMessage(u.UserName, msg)
}

// 保存用户数据
func (u *User) Save() {
	UserDB[u.UserName] = User{
		UserName:      u.UserName,
		RoomId:        u.RoomId,
		LastLoginTime: u.LastLoginTime,
		OnlineTime:    u.OnlineTime,
		IsOnline:      u.IsOnline,
	}
}

// 回复消息
func (u *User) Reply(str string) {
	u.conn.Write([]byte(str))
}
