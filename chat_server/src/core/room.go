package core

import (
	"chat_server/src/common"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Room struct {
	mutex sync.Mutex

	Id        int64
	Users     map[string]*User // 用户群
	AllRecord []ChatRecord     // 所有聊天记录 正序

	Left      int              // 用于聊天记录10分钟的左指针 当 广播消息/查询热词 的时候移动
	HotRecord map[string]int64 // 近10分钟词出现的次数
}

// 聊天记录
type ChatRecord struct {
	SendUserName string   // 发送人
	Content      string   // 内容
	SendTime     int64    // 发送时间戳 单位s
	Words        []string // 词
}

// 加入房间
func (r *Room) JoinRoom(user *User) {
	r.mutex.Lock()
	r.mutex.Unlock()

	r.Users[user.UserName] = user
	user.reply(fmt.Sprintf("欢迎进入房间号:%d\n", r.Id))

	// 发送近50条聊天记录
	s := len(r.AllRecord) - 50
	if s < 0 {
		s = 0
	}

	sb := strings.Builder{}
	defer sb.Reset()
	for _, v := range r.AllRecord[s:] {
		sb.WriteString(fmt.Sprintf("【%s】: %s\n", v.SendUserName, v.Content))
	}
	user.reply(sb.String())
}

// 退出房间
func (r *Room) ExitRoom(user *User) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.Users, user.UserName)

	user.reply(fmt.Sprintf("你退出了房间号:%d", r.Id))
}

// 广播消息
func (r *Room) BroadMessage(userName, msg string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 脏词过滤
	msg = common.Trie.Replace(msg, '*')

	// 分词
	words := strings.Split(msg, " ")
	tm := time.Now().Unix()

	record := ChatRecord{
		SendUserName: userName,
		Content:      msg,
		SendTime:     tm,
		Words:        words,
	}
	r.AllRecord = append(r.AllRecord, record)

	// 加入词语出现频率榜
	for _, word := range words {
		r.HotRecord[word]++
	}

	// 移动左指针  两个条件：1.时间小于10分钟前 2.左指针不能超过记录长度
	for r.AllRecord[r.Left].SendTime < tm-10*60*60 && r.Left < len(r.AllRecord) {
		for _, word := range r.AllRecord[r.Left].Words {
			r.HotRecord[word]--
		}
		r.Left++
	}

	for _, user := range r.Users {
		fmt.Println(user, user.IsOnline)
		if user != nil && user.IsOnline {
			fmt.Println("广播消息", fmt.Sprintf(CHAT_TMP, record.SendUserName, record.Content))
			user.reply(fmt.Sprintf(CHAT_TMP, record.SendUserName, record.Content))
		}
	}
}

// 获取房间热词榜
func (r *Room) GetRoomPopular() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	tm := time.Now().Unix()

	// 移动左指针  两个条件：1.时间小于10分钟前 2.左指针不能超过记录长度
	for r.AllRecord[r.Left].SendTime < tm-10*60*60 && r.Left < len(r.AllRecord) {
		for _, word := range r.AllRecord[r.Left].Words {
			r.HotRecord[word]--
		}
		r.Left++
	}

	// 找出热词
	var hotWord string
	var hotWordNum int64

	for word, num := range r.HotRecord {
		if num > hotWordNum {
			hotWord = word
			hotWordNum = num
		}
	}

	if hotWord == "" {
		fmt.Println("当前房间最近10分钟没有聊天记录~")
		return
	}

	fmt.Printf("当前房间最近10分钟频率最高的词是:%s 出现的次数为:%d \n", hotWord, hotWordNum)
}
