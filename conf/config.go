package conf

import (
	"flag"

	"github.com/gin-gonic/gin"

	"im-service/pkg/logger"
)

var (
	tomlFile = flag.String("config", "config.toml", "config file")
	// tomlFile      = flag.String("config", "/im/config.toml", "config file")
	TomlConfig *Toml
)

// 获取toml配置信息
func Load() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("config", TomlConfig)
		ctx.Next()
	}
}

func Init() {
	c, err := UnmarshalConfig(*tomlFile)
	if err != nil {
		logger.Info("配置文件解析错误: err:\n" + err.Error())
		return
	}
	TomlConfig = c
}
