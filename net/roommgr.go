package net

import (
	"errors"
	"fmt"
	"github.com/kordar/ws/iface"
	"sync"
)

type RoomManager struct {
	rooms    map[uint32]iface.IRoom // 管理的房间信息
	roomLock sync.RWMutex            // 读写连接的读写锁
	// 当前room的消息管理模块，用来绑定MsgId和对应的处理方法
	msgHandler iface.IMsgHandle
}

/*
	创建一个房间管理器
*/
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[uint32]iface.IRoom),
		msgHandler: NewMsgHandle(),
	}
}

func (roomMgr *RoomManager) JoinRoom(room iface.IRoom, limit int) error {
	// 保护共享资源Map 加写锁
	roomMgr.roomLock.Lock()
	defer roomMgr.roomLock.Unlock()

	if room.Members() >= limit {
		return errors.New("Maximum limit exceeded")
	}

	//将room添加到 RoomManager 中
	roomMgr.rooms[room.GetRoomID()] = room

	fmt.Println("room add to RoomManager successfully: room number = ", roomMgr.Len())

	return nil
}

// 添加房间
func (roomMgr *RoomManager) AddRoom(room iface.IRoom) {
	//保护共享资源Map 加写锁
	roomMgr.roomLock.Lock()
	defer roomMgr.roomLock.Unlock()

	//将room添加到 RoomManager 中
	roomMgr.rooms[room.GetRoomID()] = room

	fmt.Println("room add to RoomManager successfully: room number = ", roomMgr.Len())
}

// 删除房间
func (roomMgr *RoomManager) RemoveRoom(room iface.IRoom) {
	// 保护共享资源Map 加写锁
	roomMgr.roomLock.Lock()
	defer roomMgr.roomLock.Unlock()

	// 清空房间
	room.Clear()
	// 删除房间信息
	delete(roomMgr.rooms, room.GetRoomID())

	fmt.Println("rooms Remove RoomID=", room.GetRoomID(), " successfully: room num = ", roomMgr.Len())
}

// 利用RoomID获取房间
func (roomMgr *RoomManager) GetRoom(roomID uint32) (iface.IRoom, error) {
	// 保护共享资源Map 加读锁
	roomMgr.roomLock.RLock()
	defer roomMgr.roomLock.RUnlock()

	if room, ok := roomMgr.rooms[roomID]; ok {
		return room, nil
	} else {
		return nil, errors.New("room not found")
	}
}

// 获取当前房间数目
func (roomMgr *RoomManager) Len() int {
	return len(roomMgr.rooms)
}

// 清除并停止所有房间
func (roomMgr *RoomManager) ClearRooms() {
	// 保护共享资源Map 加写锁
	roomMgr.roomLock.Lock()
	defer roomMgr.roomLock.Unlock()

	// 停止并删除全部的连接信息
	for roomID, room := range roomMgr.rooms {
		// 停止
		room.Clear()
		// 删除
		delete(roomMgr.rooms, roomID)
	}

	fmt.Println("Clear All rooms successfully: room num = ", roomMgr.Len())
}

//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (roomMgr *RoomManager) AddRouter(msgId uint32, router iface.IRouter) {
	roomMgr.msgHandler.AddRouter(msgId, router)
}

func (roomMgr *RoomManager) GetMsgHandler() iface.IMsgHandle {
	return roomMgr.msgHandler
}
