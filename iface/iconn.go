package iface

type IConnection interface {
	// 启动连接，让当前连接开始工作
	Start()
	// 停止连接，结束当前连接状态M
	Stop()

	// 获取当前房间
	GetRoom() IRoom

	// 获取当前连接ID
	GetConnID() uint32
	// 直接将Message数据发送数据给远程的TCP客户端(无缓冲)
	SendMsg(msgId uint32, message string, data interface{}) error
	SendBuffMsg(msgId uint32, message string, data interface{}) error

	// 设置链接属性
	SetProperty(key string, value interface{})
	// 获取链接属性
	GetProperty(key string) (interface{}, error)
	// 移除链接属性
	RemoveProperty(key string)
}
