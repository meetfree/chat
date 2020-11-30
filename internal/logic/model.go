package logic

import (
	"time"

	"im-service/pkg/util"
)

const TopicOnline = "online"
const TopicChat = "chat"
const Separator = "☠"
const MsgTypeIsMessage = "1" // 聊天消息
const MsgTypeIsImage = "2"   // 聊天图片

type MessageModel struct {
	MsgId     int64  `json:"req,omitempty"`                                        // 消息ID
	MsgBody   string `form:"content" validate:"required" json:"content,omitempty"` // 消息的内容
	MsgType   string `form:"type" validate:"required" json:"type,omitempty"`       // 消息类别，文字或图片
	Client    string `form:"client" validate:"required" json:"client,omitempty"`   // web或app
	RoomId    string `json:"room,omitempty"`                                       // 房间号
	Nickname  string `form:"nick" json:"nick,omitempty"`                           // 用户名
	Username  string `form:"name" json:"name,omitempty"`						   // 用户昵称 new add
	Level     int64  `form:"level" json:"level,omitempty"`                         // 用户等级 new add
	Avatar    string `form:"avatar" json:"avatar,omitempty"`                       // 头像
	Color     string `form:"color" json:"color,omitempty"`                       // 颜色
	Timestamp int64  `json:"timestamp,omitempty"`                                  // 发送时间(时间戳)
	Ip        string `json:"ip"`                                                   // ip
}

func NewMessageModel(roomId string, nickname string, msgBody string, msgType string, client string) *MessageModel {
	return &MessageModel{
		MsgId:     util.GenId(),
		MsgBody:   msgBody,
		RoomId:    roomId,
		MsgType:   msgType,
		Client:    client,
		Nickname:  nickname,
		Timestamp: time.Now().Unix(),
	}
}

type message struct {
	Topic string      `json:"topic,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

func NewMessage(topic string, data interface{}) *message {
	return &message{
		Topic: topic,
		Data:  data,
	}
}
