package logger

import (
	"errors"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"ohedu/config"
)

// NewLogger 用于创建日志对象
func NewLogger(config *config.LoggerConfig) (*zap.SugaredLogger, error) {
	tmp, err := NewLoggerCore(config)
	if err != nil {
		return nil, err
	}

	return zap.New(
		tmp,
		zap.AddCaller(),
	).Sugar(), nil
}

// NewLoggerCore 用于创建日志对象
func NewLoggerCore(config *config.LoggerConfig) (zapcore.Core, error) {

	// 初始化日志
	var encoderConfig = zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02-15:04:05")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 获取编码器，NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	var encoder = zapcore.NewConsoleEncoder(encoderConfig)

	if config.File == "" {
		return nil, errors.New("未配置日志目录")
	}
	if config.MaxSize <= 0 {
		config.MaxSize = 8
	}
	if config.MaxBackups <= 0 {
		config.MaxBackups = 8
	}

	// 创建自动分割的日志文件
	var fileWriteSyncer = zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.File,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		Compress:   config.Compress,
	})

	return zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(fileWriteSyncer, zapcore.AddSync(os.Stdout)),
		getLoggerLevel(config.Level),
	), nil
}

// NewTestingLogger 用于创建用于单元测试的日志对象
func NewTestingLogger() zapcore.Core {
	// 初始化日志
	var encoderConfig = zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02-15:04:05")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 获取编码器，NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	var encoder = zapcore.NewConsoleEncoder(encoderConfig)

	return zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stderr),
		zapcore.DebugLevel,
	)
}

func getLoggerLevel(lvl string) zapcore.Level {
	var level zapcore.Level
	err := level.UnmarshalText([]byte(lvl))
	if err != nil {
		return zapcore.InfoLevel
	}
	return level
}
