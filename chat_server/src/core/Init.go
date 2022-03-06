package core

import (
	"fmt"
	"sync"
	"time"
)

var Rooms *RoomGroup       // 房间组
var UserDB map[string]User // 用户库

const (
	CHAT_TMP = "【%s】: %s" // 消息模板
)

func Init() {
	Rooms = &RoomGroup{Rooms: sync.Map{}}
	UserDB = make(map[string]User)
}

// 获取用户的在线时长
func GetUserInfo(str string) {
	user := UserDB[str]
	if user.UserName == "" {
		fmt.Println("没有该用户~")
		return
	}

	if user.IsOnline {
		tm := time.Now().Unix() - user.LastLoginTime + user.OnlineTime
		fmt.Printf("用户: %s ,在线时长：%d秒, 房间号：%d", user.UserName, tm, user.RoomId)
	} else {
		fmt.Printf("用户: %s ,在线时长：%d秒, 房间号：%d", user.UserName, user.OnlineTime, user.RoomId)
	}
}

// 获取当前房间10分钟
func GetPopular(id int64) {
	room := Rooms.GetRoom(id)
	if room == nil {
		fmt.Println("房间不存在")
		return
	}
	if len(room.AllRecord) == 0 {
		fmt.Println("当前房间最近10分钟没有聊天记录~")
		return
	}
	room.GetRoomPopular()
}
