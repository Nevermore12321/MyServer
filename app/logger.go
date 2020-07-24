package app

import (
	"MyServer/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func loggerConfigure() *zap.Logger {
	//  获取 writer
	lumerjackWriteSyncer := getLoggerWriter()

	//  获取 编码 格式配置
	encoder := getEncoder()

	//  读取配置文件， logLevel 参数
	logLevel := config.ConfigViper.GetString("log.logLevel")
	var zapLogLevel zap.AtomicLevel

	switch logLevel {
	case "debug":
		zapLogLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
		break
	case "info":
		zapLogLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
		break
	case "warn":
		zapLogLevel = zap.NewAtomicLevelAt(zap.WarnLevel)
		break
	case "error":
		zapLogLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
		break
	case "fatal":
		zapLogLevel = zap.NewAtomicLevelAt(zap.FatalLevel)
		break
	default:
		zapLogLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	//  创建 新的  zap core
	//  func NewCore(enc Encoder, ws WriteSyncer, enab LevelEnabler) Core
	//  第二个参数为 写入到 XXX, 可以设置多个 ,最后一个参数为  日志的级别
	myCore := zapcore.NewCore(encoder, lumerjackWriteSyncer, zapLogLevel)

	//  开启 开发者模式,  堆栈跟踪 zap.AddCaller() zap.Development()
	//  堆栈 关键词 为  caller
	//  实例化 logger
	logger := zap.New(myCore, zap.AddCaller(), zap.Development())

	return logger
}

//  WriteSyncer  log日志的 写入的 实现了 writer 接口
//  使用 Lumberjack 来 实现 写入 日志文件
func getLoggerWriter() zapcore.WriteSyncer {
	//  读取 log 的配置信息
	logFilename := config.ConfigViper.GetString("log.logFilename")
	logFileFullPath := "logs/" + logFilename

	//  lumberjackLogger 其实是一个 io.Writter
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFileFullPath, // 日志文件的位置
		MaxSize:    10,              // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxAge:     3,               // 保留旧文件的最大天数
		MaxBackups: 28,              // 保留旧文件的最大个数
		Compress:   false,           // 是否压缩/归档旧文件
	}

	//  AddSync将io.Writer转换为WriteSyncer.
	//  NewMultiWriteSyncer  可以同时 写入多个地方， 同时将日志写入 stdout 和 日志文件
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberjackLogger))
}

//  对 zapcore 中 的日志 进行配置， 例如: 日志级别，编码格式等
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	//  设置事件格式为 UTC
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//  设置  日志 级别 的 格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	//  返回  终端打印 格式
	return zapcore.NewConsoleEncoder(encoderConfig)
}
