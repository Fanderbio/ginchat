package models

import "gorm.io/gorm"

// 消息
type Message struct {
	gorm.Model
	FromId   uint   // 消息的发送方
	TargetId uint   // 消息接收者
	Type     string // 消息类型 群聊 私聊 广播
	Media    int    // 消息类型 文字 图片 音频
	Content  string // 消息内容
	Pic      string
	Url      string
	Desc     string // 描述相关的
	Amount   int    // 其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}
