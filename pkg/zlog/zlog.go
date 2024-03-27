package zlog

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger = InitZlog()

func InitZlog() *zap.Logger {
	coreList := make([]zapcore.Core, 0)
	// 输出到 终端
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	coreList = append(coreList, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.InfoLevel))

	// 输出到文件
	// infoJsonEncoder := &lumberjack.Logger{
	// 	Filename:   "logs/error.log",
	// 	MaxSize:    10,
	// 	MaxBackups: 10,
	// 	MaxAge:     7,
	// 	Compress:   true,
	// }
	// coreList = append(coreList, zapcore.NewCore(encoder, zapcore.AddSync(infoJsonEncoder), zap.ErrorLevel))
	core := zapcore.NewTee(coreList...)
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))
}

func Zlog() *zap.Logger {
	return logger
}
