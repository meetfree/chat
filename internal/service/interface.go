package service

import "github.com/gin-gonic/gin"

type IChat interface {
	Run(ctx *gin.Context)
	// 重新加载字典
	Reload(ctx *gin.Context)
}

type INotice interface {
	Run(ctx *gin.Context)
}

type ICli interface {
	Run(ctx *gin.Context)
	Stop(ctx *gin.Context)
}
