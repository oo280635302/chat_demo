package core

import (
	"sync"
)

type Room struct {
	Lock sync.Locker

	Id        int64
	Users     map[string]*User // 用户群
	AllRecord []ChatRecord     // 所有聊天记录 逆序
	HotRecord []string         // 热词榜
}

// 聊天记录
type ChatRecord struct {
	SendUserName string // 发送人
	Content      string // 内容
	SendTime     int64  // 发送时间
}

// 获取最近的聊天记录 50条
func (r *Room) GetRoomRecord() {

}

// 广播消息
func (r *Room) BroadMessage() {

}

// 获取房间热词榜
func (r *Room) GetRoomPopular() {

}
