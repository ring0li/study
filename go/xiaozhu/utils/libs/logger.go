package libs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

//日志还是不分开好，便于排查问题
var Logger *zap.Logger

func init() {
	// First, define our level-handling logic.
	// 仅打印Error级别以上的日志
	//highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	return lvl >= zapcore.ErrorLevel
	//})
	// 打印所有级别的日志

	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	hook := lumberjack.Logger{
		Filename:   "./",
		MaxSize:    1024, // megabytes
		MaxBackups: 3,
		MaxAge:     7,     //days
		Compress:   false, // disabled by default
	}

	fileWriter := zapcore.AddSync(&hook)

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)

	// Optimize the Kafka output for machine consumption and the console output
	// for human operators.
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.

	core := zapcore.NewTee(
		// 打印在控制台
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
		// 打印在文件中
		zapcore.NewCore(consoleEncoder, fileWriter, lowPriority),
	)

	// From a zapcore.Core, it's easy to construct a Logger.
	Logger = zap.New(core)
	defer Logger.Sync()
}
