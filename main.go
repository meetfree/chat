package main

import (
	_ "net/http/pprof"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"im-service/conf"
	"im-service/internal/logic"
	"im-service/internal/scheduler"
	"im-service/internal/service"
	"im-service/pkg/db/mgo"
	"im-service/pkg/db/redis"
	"im-service/pkg/filter"
	"im-service/pkg/mq"
)

// init 初始化配置
func init() {
	gin.SetMode(gin.ReleaseMode)
	conf.Init()
	redis.Open()
	mgo.Open()
	mq.Open()
	filter.Init()
}

func main() {
	var (
		chat service.IChat = new(service.ChatService)
		// cli  service.ICli  = new(service.CliService)
	)
	router := gin.New()
	router.Use(gin.Recovery())
	// router := gin.Default()
	// 设置跨域
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"x-xq5-jwt"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(conf.Load())
	h := logic.NewHub()
	go h.Run()
	go scheduler.Count(h.Counter)
	go scheduler.Clear(h.Rooms)
	router.GET("/websocket/:room", func(ctx *gin.Context) {
		h.Websocket(ctx)
	})
	// router.POST("/chat", func(ctx *gin.Context) {
	// 	chat.Run(ctx)
	// })
	router.GET("/reload", func(ctx *gin.Context) {
		chat.Reload(ctx)
	})
	_ = router.Run(conf.TomlConfig.GetListenAddr())
}
