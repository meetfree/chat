package util

import (
	"fmt"
	"runtime"

	"im-service/pkg/logger"
)

// RecoverPanic 恢复panic
func RecoverPanic() {
	if err := recover(); err != nil {
		logger.ZapLogger.Error(GetStackInfo())
	}
}

// GetStackInfo 获取Panic堆栈信息
func GetStackInfo() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return fmt.Sprintf("%s", buf[:n])
}
