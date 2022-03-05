package core

import (
	"fmt"
	"time"
)

var UserDB map[string]User // 用户库
var RoomDB map[int64]*Room // 房间库

const (
	CHAT_TMP = "【%s】: %s" // 消息模板
)

func Init() {
	UserDB = make(map[string]User)
	RoomDB = make(map[int64]*Room)
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
	room := RoomDB[id]
	if room == nil || len(room.AllRecord) == 0 {
		fmt.Println("当前房间最近10分钟没有聊天记录~")
		return
	}
	room.GetRoomPopular()
}
