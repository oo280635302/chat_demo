package core

import (
	pb "chat_server/src/pb"
	"fmt"
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
		u.reply("不要重复登录")
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
			u.reply("该用户正在登陆")
			return
		}

		u.RoomId = 0
		u.LastLoginTime = tm
		u.IsOnline = true
	}
	u.UserName = req.UserName

	// 保存数据库
	u.save()

	u.reply("登陆成功")
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
	u.save()
}

// 加入房间
func (u *User) JoinRoom(b []byte) {
	req := new(pb.JoinRoomReq)
	proto.Unmarshal(b, req)

	if req.RoomId <= 0 {
		u.reply("无效的房间号")
		return
	}

	u.mutex.Lock()
	defer u.mutex.Unlock()

	// 退出之前的房间
	if u.RoomId != 0 {
		lastRoom := Rooms.GetRoom(u.RoomId)
		if lastRoom != nil {
			lastRoom.ExitRoom(u)
		}
	}

	u.RoomId = req.RoomId

	room := Rooms.GetOrCreateRoom(u.RoomId)
	room.JoinRoom(u)

	// 保存数据
	u.save()
}

// 发送消息
func (u *User) SendMessage(b []byte) {
	req := new(pb.SendMessageReq)
	proto.Unmarshal(b, req)

	if len(u.UserName) == 0 {
		u.reply("需要登录")
		return
	}

	if u.RoomId == 0 {
		u.reply("需要在房间内才能聊天")
		return
	}

	room := Rooms.GetRoom(u.RoomId)
	if room == nil {
		fmt.Println("ERROR:User SendMessage Get Room err:", u.RoomId)
		return
	}

	// 广播
	room.BroadMessage(u.UserName, req.Msg)
}

// 保存用户数据
func (u *User) save() {
	UserDB[u.UserName] = User{
		UserName:      u.UserName,
		RoomId:        u.RoomId,
		LastLoginTime: u.LastLoginTime,
		OnlineTime:    u.OnlineTime,
		IsOnline:      u.IsOnline,
	}
}

// 回复消息
func (u *User) reply(str string) {
	u.conn.Write([]byte(str))
}
