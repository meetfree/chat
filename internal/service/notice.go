package service

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	"im-service/pkg/mq"
	"im-service/pkg/util"
)

type NoticeService struct {
	Content string `json:"content,omitempty"` // 消息类别
}

func (ns *NoticeService) Run(ctx *gin.Context) {
	defer util.RecoverPanic()
	var (
		data []byte
		err  error
	)
	data = []byte(ctx.PostForm("data"))
	if data, err = json.Marshal(ns); err != nil {
		panic(err)
	}
	go mq.Pub("notice", data)
}
