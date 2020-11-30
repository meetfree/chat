package mq

import (
	"im-service/conf"
	"time"

	"github.com/nats-io/nats.go"
)

//const addr string = "nats://localhost:4222"

var (
	Conn *nats.Conn
	err  error
)

func Open() {
	if Conn, err = nats.Connect(conf.TomlConfig.NatConfig(), nats.Name(""), nats.Timeout(10*time.Second)); err != nil {
		panic("消息队列服务未启用或异常")
	}
}

func Pub(topic string, data []byte) {
	if err = Conn.Publish(topic, data); err != nil {
		panic("消息发布异常")
	}
}

// func SubSync() {
func Sub(topic string, ch chan []byte) {
	if _, err := Conn.Subscribe(topic, func(msg *nats.Msg) {
		ch <- msg.Data
	}); err != nil {
		panic("消息订阅异常")
	}
}
