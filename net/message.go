package net

type Message struct {
	Code    uint32      `json:"code,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"msg,omitempty"`
}

// 获取消息数据段长度
func (msg *Message) GetMessage() string {
	return msg.Message
}

// 获取消息ID
func (msg *Message) GetCode() uint32 {
	return msg.Code
}

// 获取消息内容
func (msg *Message) GetData() interface{} {
	return msg.Data
}

// 设置消息数据段长度
func (msg *Message) SetMessage(message string) {
	msg.Message = message
}

// 设计消息ID
func (msg *Message) SetCode(code uint32) {
	msg.Code = code
}

// 设计消息内容
func (msg *Message) SetData(data interface{}) {
	msg.Data = data
}
