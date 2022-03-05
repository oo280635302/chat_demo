package core

var UserDB map[string]User // 用户数据库
var RoomDB map[string]User // 房间数据库

func Init() {
	UserDB = make(map[string]User)
	RoomDB = make(map[string]User)
}
