package net

import (
	"errors"
	"fmt"
	"github.com/kordar/ws/iface"
)

type Room struct {
	// 房间ID
	RoomID uint32
	//
	ConnManager iface.IConnManager

	//该Room的连接创建时Hook函数
	OnConnStart func(conn iface.IConnection)
	//该Room的连接断开时的Hook函数
	OnConnStop func(conn iface.IConnection)

	//链接属性
	property map[string]interface{}
}

// 获取房间号
func (room *Room) GetRoomID() uint32 {
	return room.RoomID
}

// 停止房间服务
func (room *Room) Clear() {
	room.ConnManager.ClearConn()
}

// 得到链接管理器
func (room *Room) GetConnMgr() iface.IConnManager {
	return room.ConnManager
}

//调用连接OnConnStart Hook函数
func (room *Room) CallOnConnStart(conn iface.IConnection) {
	if room.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		room.OnConnStart(conn)
	}
}

//调用连接OnConnStop Hook函数
func (room *Room) CallOnConnStop(conn iface.IConnection) {
	if room.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		room.OnConnStop(conn)
	}
}

//设置该Server的连接创建时Hook函数
func (room *Room) SetOnConnStart(hookFunc func(iface.IConnection)) {
	room.OnConnStart = hookFunc
}

//设置该Server的连接断开时的Hook函数
func (room *Room) SetOnConnStop(hookFunc func(iface.IConnection)) {
	room.OnConnStop = hookFunc
}

// 房间人数
func (room *Room) Members() int {
	return room.ConnManager.Len()
}

// 获取房间属性
func (room *Room) GetProperty(key string) (interface{}, error) {
	if value, ok := room.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (room *Room) Broadcast(msgId uint32, message string, data interface{}, ignore uint32)  {
	for _,conn := range room.ConnManager.GetAll() {
		if (conn.GetConnID() != ignore) {
			_ = conn.SendBuffMsg(msgId, message, data)
		}
	}
}

func NewRoom(roomID uint32) iface.IRoom {
	return &Room{
		ConnManager: NewConnManager(),
		RoomID:      roomID,
	}
}


