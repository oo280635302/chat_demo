package core

var onlineMap map[string]*Client

func Init() {
	onlineMap = make(map[string]*Client)
}
