package iface

type IRoom interface {
	// 获取房间号
	GetRoomID() uint32
	// 得到链接管理器
	GetConnMgr() IConnManager
	// 停止房间服务
	Clear()
	// 房间人数
	Members() int
	// 获取属性
	GetProperty(key string) (interface{}, error)
	//设置该Server的连接创建时Hook函数
	SetOnConnStart(func(IConnection))
	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(IConnection))
	//调用连接OnConnStart Hook函数
	CallOnConnStart(conn IConnection)
	//调用连接OnConnStop Hook函数
	CallOnConnStop(conn IConnection)
	// 广播
	Broadcast(msgId uint32, message string, data interface{}, ignore uint32)
}
