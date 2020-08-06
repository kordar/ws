package iface

/*
	房间管理抽象
*/
type IRoomMgr interface {
	AddRoom(room IRoom)                   // 添加房间
	RemoveRoom(room IRoom)                // 删除房间
	GetRoom(roomID uint32) (IRoom, error) // 利用roomID获取房间
	Len() int                             // 获取当前房间
	ClearRooms()                          // 删除并停止所有房间
	// 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	AddRouter(msgId uint32, router IRouter)
	GetMsgHandler() IMsgHandle
}
