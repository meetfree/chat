package service

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"

	"im-service/pkg/logger"
)

type CliService struct {
}

func (*CliService) Run(ctx *gin.Context) {
	var (
		out []byte
		err error
		cmd *exec.Cmd
	)
	cmd = exec.Command("sh", "./s8game/start.sh")
	if out, err = cmd.Output(); err != nil {
		logger.ZapLogger.Warn(err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{"res": "服务已启动" + string(out)})
}

func (*CliService) Stop(ctx *gin.Context) {
	var (
		out []byte
		err error
		cmd *exec.Cmd
	)
	cmd = exec.Command("sh", "./s8game/stop.sh")
	if out, err = cmd.Output(); err != nil {
		logger.ZapLogger.Warn(err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{"res": "服务已停止" + string(out)})
}
