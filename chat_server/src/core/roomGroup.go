package core

import "sync"

type RoomGroup struct {
	Rooms sync.Map // key int64,value *Room
}

// 获取/新增房间信息
func (g *RoomGroup) GetOrCreateRoom(roomId int64) *Room {
	roomI, ok := g.Rooms.Load(roomId)
	if ok {
		return roomI.(*Room)
	}

	tryRoom := &Room{}
	tryRoom.Id = roomId
	tryRoom.Users = make(map[string]*User, 100)
	tryRoom.HotRecord = make(map[string]int64, 100)

	roomI, _ = g.Rooms.LoadOrStore(roomId, tryRoom)
	return roomI.(*Room)
}

// 获取房间信息
func (g *RoomGroup) GetRoom(roomId int64) *Room {
	roomI, _ := g.Rooms.Load(roomId)
	return roomI.(*Room)
}
