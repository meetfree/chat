package logic

import (
	"bytes"
	"encoding/json"
	"sync"
	"time"

	"im-service/pkg/db/mgo"
	"im-service/pkg/filter"
	"im-service/pkg/mq"
	"im-service/pkg/pb"
	"im-service/pkg/util"
)

type Hub struct {
	Rooms      map[string][]*Client
	Clients    map[*Client]string
	Broadcast  chan []byte   // 广播内容
	Temp       chan []byte
	Register   chan *Client
	Unregister chan *Client
	Counter    map[string]int
	Locker     *sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte, 1024),
		Temp:       make(chan []byte, 1024),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]string),
		Rooms:      make(map[string][]*Client),
		Counter:    make(map[string]int),
		Locker:     new(sync.RWMutex),
	}
}

// 开始消息读写队列
func (h *Hub) Run() {
	util.RecoverPanic()
	// 订阅消息
	go mq.Sub(TopicChat, h.Temp)
	go mq.Sub(TopicOnline, h.Broadcast)
	for {
		select {
		// socket client chan
		case client := <-h.Register:
			h.Clients[client] = string(client.RoomId)
			h.Rooms[string(client.RoomId)] = append(h.Rooms[string(client.RoomId)], client)
			// 读锁
			h.Locker.Lock()
			if h.Counter[string(client.RoomId)] == 0 {
				h.Counter[string(client.RoomId)] = 1
				//h.Counter[string(client.RoomId)][client.Token] = true
			} else {
				//h.Counter[string(client.RoomId)][client.Token] = true
				h.Counter[string(client.RoomId)] += 1
			}

			// socket close people number -1
			client.Conn.SetCloseHandler(func(code int, text string) error {
				h.Counter[string(client.RoomId)] -= 1
				return nil
			})

			// 释放互斥锁
			h.Locker.Unlock()

			// 进入聊天室欢迎信息
			// go func() {
			// 	if string(client.Nickname) != "" {
			// 		for _, i := range h.Rooms[string(client.RoomId)] {
			// 			if !bytes.Equal(i.Nickname, client.Nickname) {
			// 				h.Broadcast <- []byte(string(client.Nickname) + "进入聊天室")
			// 			}
			// 		}
			// 	}
			// }()
		// 读消息
		//case client := <-h.Unregister:
		//	if _, ok := h.Clients[client]; ok {
		//		delete(h.Clients, client)
		//		close(client.Send)
		//	}

		// 消息存取转二进制
		case t := <-h.Temp:
			var m = new(MessageModel)
			_ = json.Unmarshal(t, m)
			h.handleMsg(m)

		// 发消息
		case message := <-h.Broadcast:
			for client := range h.Clients {
				// 禁言的时候跳过  新的禁言在接口处理了
				//if client.Mute {
				//	continue
				//}
				//fmt.Println(client.RoomId)
				array := bytes.Split(message, []byte(Separator))

				// 房间 与消息 房间相同时发送消息
				if bytes.Equal(client.RoomId, array[0]) {
					select {
					case client.Send <- array[1]:
						// client.Updated = time.Now()
					default:
						// 关闭发送通道
						close(client.Send)
						// 删除连接
						delete(h.Clients, client)
						// 移出聊天室
						delete(h.Rooms, string(client.RoomId))
					}
				}
			}
		}
	}
}

// 消息句柄操作
func (h *Hub) handleMsg(m *MessageModel) {
	m.MsgId = util.GenId()
	m.MsgBody = filter.Parse(m.MsgBody)
	m.Timestamp = time.Now().Unix()
	// 保存用户发送的消息
	if m.MsgType == MsgTypeIsMessage || m.MsgType == MsgTypeIsImage {
		go mgo.Save(m.RoomId, m)
	}

	// protocol buffer 转换消息内容(二进制)
	b := pb.Encode(pb.NewDownstream(pb.Encode(pb.NewMessageWrapper(m.Client, m.Nickname, m.MsgBody, m.Username, m.Level, m.Avatar, m.Color))))
	_ = pb.DecodeDown(new(pb.Downstream), b)
	h.Broadcast <- util.BytesCombine([]byte(m.RoomId+Separator), b)
}
