package net

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kordar/ws/iface"
	"github.com/kordar/ws/utils"
	"sync"
)

type Connection struct {
	// 当前conn属于的房间
	room iface.IRoom
	// 当前连接的 socket TCP 套接字
	Conn *websocket.Conn
	// 当前连接的 ID 也可以称作为 SessionID，ID 全局唯一
	ConnID uint32
	// 当前连接的关闭状态
	isClosed bool
	// 告知该链接已经退出/停止的channel
	ExitBuffChan chan bool
	// 无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan chan []byte
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte

	//消息管理MsgId和对应处理方法的消息管理模块
	MsgHandler iface.IMsgHandle

	//链接属性
	property map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
}

// 创建连接的方法
func NewConnection(room iface.IRoom, conn *websocket.Conn, connID uint32, msgHandler iface.IMsgHandle) *Connection {
	// 初始化Conn属性
	c := &Connection{
		room:         room,
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property:     make(map[string]interface{}),
	}
	// 将新创建的Conn添加到链接管理中
	c.room.GetConnMgr().Add(c)
	return c
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

/*
	读消息Goroutine，用于从客户端中读取数据
*/
func (c *Connection) StartReader() {

	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println("[conn Reader exit!]")
	defer c.Stop()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println("read msg head error ", err)
			break
		}

		fmt.Println(string(message))

		var m Message
		if err := json.Unmarshal(message, &m); err != nil {
			fmt.Println("read msg format bad!!")
			break
		}

		req := Request{
			conn: c,
			msg:  &m,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经启动工作池机制，将消息交给Worker处理
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//从绑定好的消息和对应的处理方法中执行对应的Handle方法
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println("[conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			// 有数据传给客户端
			if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}
		case <-c.ExitBuffChan:
			_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
	}
}

// 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	//1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.room.CallOnConnStart(c)
}

//停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()...ConnID = ", c.ConnID)
	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.room.CallOnConnStop(c)

	// 关闭socket链接
	_ = c.Conn.Close()
	//关闭Writer
	c.ExitBuffChan <- true

	//将链接从连接管理器中删除
	c.room.GetConnMgr().Remove(c)

	//关闭该链接全部管道
	close(c.ExitBuffChan)
	close(c.msgBuffChan)
}

//获取链接属性
func (c *Connection) GetRoom() iface.IRoom {
	return c.room
}

//直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgId uint32, msgMessage string, msgData interface{}) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	// 将data封包，并且发送
	m := Message{Code: msgId, Data: msgData, Message: msgMessage}

	if data, err := json.Marshal(m); err != nil {
		return err
	} else {
		// 写回客户端
		c.msgChan <- data

		return nil
	}
}

func (c *Connection) SendBuffMsg(msgId uint32, msgMessage string, msgData interface{}) error {
	if c.isClosed == true {
		return errors.New("connection closed when send buff msg")
	}

	m := Message{Code: msgId, Data: msgData, Message: msgMessage}

	if data, err := json.Marshal(m); err != nil {
		return err
	} else {
		// 写回客户端
		c.msgBuffChan <- data

		return nil
	}

}


//设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

//获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

//移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
