package service

import (
	"net/http"

	"im-service/pkg/filter"
	"im-service/pkg/util"

	"github.com/gin-gonic/gin"
)

var (
	data []byte
	err  error
)

type ChatService struct {
	err error
}

func (cs *ChatService) Run(ctx *gin.Context) {
	defer util.RecoverPanic()
	var (
	// mm         = new(pb.MessageModel)
	// downstream pb.Downstream
	// b          []byte
	)
	// if b := !base.Authorize(ctx); b {
	// 	return
	// }
	// if err = ctx.ShouldBind(&mm); err != nil {
	// 	log.Warnf("Websocket err:%v\n", cs.err.Error())
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"msg": "缺少参数或参数不匹配，请重试", "success": false, "code": http.StatusBadRequest})
	// 	return
	// }

	// mm.RoomId = ctx.PostForm("room")
	// mm.MsgId = util.GenId()
	// timestamp := time.Now().Unix()
	// mm.Timestamp = fmt.Sprintf("%d", timestamp)
	// mm.Ip = ctx.ClientIP()
	// 保存发送消息记录
	// go mgo.Save(mm.RoomId, mm)
	//
	// ds := pb.NewDownstream(mm)
	// b = pb.Encode(ds)
	// if err = pb.DecodeDown(&downstream, b); err != nil {
	// 	logger.Error(err)
	// }
	// go mq.Pub(base.TopicChat, util.BytesCombine([]byte(mm.RoomId+base.Separator), b))
	// ctx.JSON(http.StatusOK, gin.H{"msg": "消息发送成功", "success": true, "code": http.StatusOK})
}

// func (cs *ChatService) Recall(ctx *gin.Context) {
// 	msgId := ctx.PostForm("id")
// 	data, err = json.Marshal(model.NewMessage("recall", msgId))
// 	go mq.Pub(base.TopicChat, data)
// }

// 加载字典
func (cs *ChatService) Reload(ctx *gin.Context) {
	defer util.RecoverPanic()
	filter.Reload()
	ctx.JSON(http.StatusOK, gin.H{"msg": "reload success", "success": true})
}

// 拉取对话记录
// func Pull(room string, conn *websocket.Conn) {
//	var (
//		results []*Response
//		cur     *mongo.Cursor
//		err     error
//	)
//	if room == "" {
//		return
//	}
//	db := mgo.NewMgo(room)
//	if cur, err = db.GetAll(); err != nil {
//		log.Info(err)
//	}
//	// 查找多个文档返回一个游标，遍历游标允许我们一次解码一个文档
//	for cur.Next(context.TODO()) {
//		// 创建一个值，将单个文档解码为该值
//		var elem Response
//		err := cur.Decode(&elem)
//		if err != nil {
//			log.Fatal(err)
//		}
//		results = append(results, &elem)
//	}
//	log.Info(len(results))
//	for _, i := range results {
//		//c.send <- []byte(i.MsgBody)
//		conn.WriteMessage(websocket.TextMessage, []byte(i.MsgBody))
//	}
// }
