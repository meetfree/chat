package logger

import (
	"fmt"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ZapLogger *zap.Logger
	err       error
)

func init() {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		// EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 这里可以指定颜色
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atom := zap.NewAtomicLevelAt(zap.InfoLevel)
	config := zap.Config{
		Level:            atom,                         // 日志级别
		Development:      true,                         // 开发模式，堆栈跟踪
		Encoding:         "console",                    // 输出格式 console 或 json
		EncoderConfig:    encoderConfig,                // 编码器配置
		InitialFields:    make(map[string]interface{}), // 初始化字段，如：添加一个服务器名称
		OutputPaths:      []string{"stdout"},           // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr"},
	}
	config.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder // 这里可以指定颜色
	// 构建日志
	ZapLogger, err = config.Build()
	if err != nil {
		panic(fmt.Sprintf("log 初始化失败: %v", err))
	}
}

// logPath 日志文件路径
// logLevel 日志级别 debug/info/warn/error
// maxSize 单个文件大小,MB
// maxBackups 保存的文件个数
// maxAge 保存的天数， 没有的话不删除
// compress 压缩
// jsonFormat 是否输出为json格式
// showLine 显示代码行
// logInConsole 是否同时输出到控制台
func initLogger(logPath string, logLevel string, maxSize, maxBackups, maxAge int, compress, jsonFormat, showLine, logInConsole bool) {
	hook := lumberjack.Logger{
		Filename:   logPath,    // 日志文件路径
		MaxSize:    maxSize,    // megabytes
		MaxBackups: maxBackups, // 最多保留300个备份
		Compress:   compress,   // 是否压缩 disabled by default
	}
	if maxAge > 0 {
		hook.MaxAge = maxAge // days
	}

	var syncer zapcore.WriteSyncer
	if logInConsole {
		syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	} else {
		syncer = zapcore.AddSync(&hook)
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	var encoder zapcore.Encoder
	if jsonFormat {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 设置日志级别,debug可以打印出info,debug,warn；info级别可以打印warn，info；warn只能打印warn
	// debug->info->warn->error
	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	core := zapcore.NewCore(
		encoder,
		syncer,
		level,
	)

	ZapLogger = zap.New(core)
	if showLine {
		ZapLogger = ZapLogger.WithOptions(zap.AddCaller())
	}
}

func Info(data ...interface{}) {
	ZapLogger.Sugar().Info(data)
	// initLogger("./logger.log", "info", 50, 10, 0, false, false, false, false)
}

func Warn(data ...interface{}) {
	ZapLogger.Sugar().Warn(data)
	// initLogger("./logger.log", "warn", 50, 10, 0, false, false, true, true)
}

func Error(data ...interface{}) {
	ZapLogger.Sugar().Error(data)
	// initLogger("./logger.log", "error", 50, 10, 0, false, false, true, true)
}

func color() {
	fmt.Println("")

	// 前景 背景 颜色
	// ---------------------------------------
	// 30  40  黑色
	// 31  41  红色
	// 32  42  绿色
	// 33  43  黄色
	// 34  44  蓝色
	// 35  45  紫红色
	// 36  46  青蓝色
	// 37  47  白色
	//
	// 代码 意义
	// -------------------------
	//  0  终端默认设置
	//  1  高亮显示
	//  4  使用下划线
	//  5  闪烁
	//  7  反白显示
	//  8  不可见

	for b := 40; b <= 47; b++ { // 背景色彩 = 40-47
		for f := 30; f <= 37; f++ { // 前景色彩 = 30-37
			for d := range []int{0, 1, 4, 5, 7, 8} { // 显示方式 = 0,1,4,5,7,8
				fmt.Printf(" %c[%d;%d;%dm%s(f=%d,b=%d,d=%d)%c[0m ", 0x1B, d, b, f, "", f, b, d, 0x1B)
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}
