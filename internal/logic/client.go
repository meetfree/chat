package logic

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"im-service/pkg/logger"
	"im-service/pkg/pb"
	"im-service/pkg/util"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	Hub      *Hub            //
	Conn     *websocket.Conn // websocket连接
	Send     chan []byte     // 发送信息的缓冲channel
	Uuid     []byte          // 用户id
	//Token    string          // 令牌（登录获取）
	Nickname []byte          // 用户名登录获取，否则统一为guest
	RoomId   []byte          // 房间id
	Mute     bool            // 是否禁言
	IpAddr   []byte          // ip地址
	Created  time.Time       // 进入时间
	Updated  time.Time       // 活跃时间
}

func NewClient(hub *Hub, conn *websocket.Conn, name []byte, room []byte, ip []byte) *Client {
	now := time.Now()
	return &Client{
		//Token:    token,
		Hub:      hub,
		Conn:     conn,
		Send:     make(chan []byte, maxMessageSize),
		Nickname: name,
		RoomId:   room,
		Mute:     false,
		IpAddr:   ip,
		Created:  now,
		Updated:  now,
	}
}

// 读取从websocket转到hub的消息
func (c *Client) Read() {
	var (
		p   []byte
		err error
	)
	defer util.RecoverPanic()
	defer func() {
		c.Hub.Unregister <- c
		_ = c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(
		func(string) error {
			_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

	for {
		if _, p, err = c.Conn.ReadMessage(); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				if !strings.Contains(err.Error(), "websocket: close 1005") {
					logger.ZapLogger.Error(err.Error())
				}
				logger.ZapLogger.Info(string(c.Nickname) + " conn is closed.")
			}
			c.Conn.Close()
			return
		} else {
			var upstream pb.Upstream
			if err = pb.DecodeUp(&upstream, p); err != nil {
				c.Conn.Close()
				return
			} else {
				c.Updated = time.Now()
			}
		}
	}
}

// 将消息从hub写入websocket连接
func (c *Client) Write() {
	// 捕获异常
	defer util.RecoverPanic()
	var (
		w   io.WriteCloser
		err error
	)
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			// 设置写入的超时时间
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// 如果取值出错关闭连接，设置写入状态和对应的数据
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				panic(err)
			}
			if w, err = c.Conn.NextWriter(websocket.BinaryMessage); err != nil {
				panic(err)
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				panic(err)
			}
		case <-ticker.C:
			// 心跳包，ping出错就会报错退出断开这个连接
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// 处理websocket请求
func (h *Hub) Websocket(ctx *gin.Context) {
	defer util.RecoverPanic()
	// if !Authorize(ctx) {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "认证失败，请重试", "success": false})
	// }
	var (
		client *Client
		conn   *websocket.Conn
		room   = ctx.Param("room")
		//room   = ctx.DefaultQuery("room", "88888888")
		nick   = ctx.Query("nickname")
		//token  = ctx.Query("token")
		err    error

		// 将网络请求变为websocket
		upgrader = &websocket.Upgrader{
			// 解决跨域问题
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 65535,
		}
	)
	fmt.Println(room)
	if conn, err = upgrader.Upgrade(ctx.Writer, ctx.Request, nil); err != nil {
		panic("升级为websocket失败，" + err.Error())
		return
	}
	if nick == "" {
		nick = "guest"
	}

	// 创建一个客户端
	client = NewClient(h, conn, []byte(nick), []byte(room), []byte(ctx.ClientIP()))

	// 注册一个客户端
	h.Register <- client
	// 写消息
	go client.Write()
	// 读消息 走的接口发消息 不用读取 websocket 消息
	//go client.Read()
	// 禁言
	//go client.Ban()
}

// 2=禁言，3=拉黑
//func (c *Client) Ban() {
//	//ticker := time.NewTicker(time.Second * 1)
//	//for range ticker.C {
//	//	if s := redis.Get(c.Token); s == "2" {
//	//		c.Mute = true
//	//	} else {
//	//		c.Mute = false
//	//	}
//	//}
//}

//func Authorize(ctx *gin.Context) bool {
//	var (
//		b = false
//		// host   = conf.Configuration.Domain
//		// origin = ""
//		token = ctx.Query("token")
//		room  = ctx.Param("room")
//	)
//	if (room != "") && (token != "" && redis.Get(token) == "1") {
//		b = true
//	}
//	// headerHost := ctx.GetHeader("Host")
//	// headerOrigin := ctx.GetHeader("Origin")
//	// if host == headerHost && origin == headerOrigin {
//	// 	b = true
//	// }
//	return b
//}
