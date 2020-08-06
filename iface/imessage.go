package iface

type IMessage interface {
	GetMessage() string   //获取消息提示内容
	GetCode() uint32      //获取消息ID
	GetData() interface{} //获取消息内容

	SetCode(uint32)      //设置消息ID
	SetData(interface{}) //设置消息内容
	SetMessage(string)   //设置消息提示内容
}
