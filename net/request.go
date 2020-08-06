package net

import "github.com/kordar/ws/iface"

type Request struct {
	conn iface.IConnection //已经和客户端建立好的 链接
	msg iface.IMessage    //客户端请求的数据
}

// 获取请求连接信息
func (r *Request) GetConnection() iface.IConnection {
	return r.conn
}

// 获取请求消息的数据
func (r *Request) GetData() interface{} {
	return r.msg.GetData()
}

// 获取请求的消息的ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetCode()
}

// 获取消息提示信息
func (r *Request) GetMessage() string {
	return r.msg.GetMessage()
}
