package filter

import (
	"github.com/importcjj/sensitive"

	"im-service/conf"
	"im-service/pkg/logger"
)

var (
	filter *sensitive.Filter
	err    error
)

func Init() {
	filter = sensitive.New()
	if err = filter.LoadNetWordDict(conf.TomlConfig.Dict); err != nil {
		logger.Warn("字典加载错误，错误原因：" + err.Error() + "，加载位置：" + conf.TomlConfig.Dict)
	}
}

// 过滤
func Parse(s string) string {
	return filter.Replace(s, '*')
}

func Reload() {
	filter = nil
	Init()
}

// 追加敏感词
func Append(w string) {
	filter.AddWord(w)
}

// 移除敏感词
func Remove(w string) {
	filter.DelWord(w)
}
