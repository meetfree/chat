package pb

import (
	"strconv"
	"time"

	"github.com/gogo/protobuf/proto"

	"im-service/pkg/filter"
	"im-service/pkg/logger"
	"im-service/pkg/util"
)

func Encode(pm proto.Message) []byte {
	util.RecoverPanic()
	var (
		data []byte
		err  error
	)
	if data, err = proto.Marshal(pm); err != nil {
		logger.Error("protobuf编码失败，错误原因：" + err.Error())
	}
	return data
}

func DecodeUp(us *Upstream, data []byte) error {
	util.RecoverPanic()
	if err := proto.Unmarshal(data, us); err != nil {
		logger.Error("protobuf解码失败，错误原因：" + err.Error())
		return err
	}
	return nil
}

func DecodeDown(ds *Downstream, data []byte) error {
	util.RecoverPanic()
	if err := proto.Unmarshal(data, ds); err != nil {
		logger.ZapLogger.Error("protobuf解码失败，错误原因：" + err.Error())
		return err
	}
	return nil
}

func NewDownstream(data []byte) *Downstream {
	return &Downstream{
		Type: PackageType_PT_MESSAGE,
		Data: data,
	}
}

func NewMessageWrapper(client string, nickname string, msgBody string, username string, level int64, avatar string, color string) *MessageWrapper {
	return &MessageWrapper{
		Topic:     "chat",
		SeqId:     util.GenId(),
		Client:    client,
		Username:  nickname,
		Nickname:  username, // new add username
		Level:     level,    // new add level
		Avatar:    avatar,   // new add avatar
		Color:     color,    // new add color
		MsgBody:   &MessageBody{MessageType: MessageType_MT_TEXT, MessageContent: filter.Parse(msgBody)},
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
	}
}

//

// NewMessageBody 创建一个消息体类型
// func NewMessageBody(msgType int, msgContent string) *Upstream {
// 	up := Upstream{}
// 	switch up.Type {
// 	case PackageType_PT_HEARTBEAT:
// 		content.Content = &MessageContent_Text{Text: &text}
// 	case MessageType_MT_FACE:
// 		var face Face
// 		err := jsoniter.Unmarshal(util.Str2bytes(msgContent), &face)
// 		if err != nil {
// 			logger.Sugar.Error(err)
// 			return nil
// 		}
// 		content.Content = &MessageContent_Face{Face: &face}
// 	case MessageType_MT_COMMAND:
// 		var command Command
// 		err := jsoniter.Unmarshal(util.Str2bytes(msgContent), &command)
// 		if err != nil {
// 			logger.Sugar.Error(err)
// 			return nil
// 		}
// 		content.Content = &MessageContent_Command{Command: &command}
// 	case MessageType_MT_CUSTOM:
// 		var custom Custom
// 		err := jsoniter.Unmarshal(util.Str2bytes(msgContent), &custom)
// 		if err != nil {
// 			logger.Sugar.Error(err)
// 			return nil
// 		}
// 		content.Content = &MessageContent_Custom{Custom: &custom}
// 	}
//
// 	return &MessageBody{
// 		MessageType:    MessageType(msgType),
// 		MessageContent: &content,
// 	}
// }
