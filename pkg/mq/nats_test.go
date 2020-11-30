package mq

import (
	"log"
	"testing"

	"github.com/nats-io/nats.go"
)

func TestConnect(t *testing.T) {
	// create server connection
	conn, _ := nats.Connect(nats.DefaultURL)
	log.Println("connected to " + nats.DefaultURL)

	// subscribe to subject
	log.Printf("subscribing to subject 'foo' \n")
	conn.Subscribe("foo", func(msg *nats.Msg) {
		//handle the message
		log.Printf("received message '%s\n", string(msg.Data)+"'")
	})
}
