package redis

import (
	"testing"

	"im-service/pkg/logger"
)

func TestConnect(t *testing.T) {
	Open()
	if err != nil {
		logger.ZapLogger.Error(err.Error())
	}
}
