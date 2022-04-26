package logger

import (
	"go.uber.org/zap"
	"ohedu/config"
)

// Zaplog 日志对象
var Zaplog *zap.SugaredLogger
var zaplogCallerSkip1 *zap.SugaredLogger

// InitLogger 初始化日志
func InitLogger() error {
	var err error
	tmp, err := NewLoggerCore(&config.Config.Logger)
	Zaplog = zap.New(
		tmp,
		zap.AddCaller(),
	).Sugar()
	zaplogCallerSkip1 = zap.New(
		tmp,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	).Sugar()
	return err
}

func InitTestingLogger() {
	Zaplog = zap.New(
		NewTestingLogger(),
		zap.AddCaller(),
	).Sugar()
	zaplogCallerSkip1 = zap.New(
		NewTestingLogger(),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	).Sugar()
}

// Debug 输出日志，用于输出一些不重要的消息，或者是调试信息。
// 临时调试信息应该在推送前删除。
func Debug(args ...interface{}) {
	zaplogCallerSkip1.Debug(args...)
}

// Debugf 格式化输出日志，用于输出一些不重要的消息，或者是调试信息。
func Debugf(template string, args ...interface{}) {
	zaplogCallerSkip1.Debugf(template, args...)
}

// Warn 输出日志，用于输出警告信息。
func Warn(args ...interface{}) {
	zaplogCallerSkip1.Warn(args...)
}

// Warnf 格式化输出日志，用于输出警告信息。
func Warnf(template string, args ...interface{}) {
	zaplogCallerSkip1.Warnf(template, args...)
}

// Info 输出日志，用于输出通常日志信息。
func Info(args ...interface{}) {
	zaplogCallerSkip1.Info(args...)
}

// Infof 格式化输出日志，用于输出通常日志信息。
func Infof(template string, args ...interface{}) {
	zaplogCallerSkip1.Infof(template, args...)
}

// Error 输出日志，用于输出严重错误信息。
func Error(args ...interface{}) {
	zaplogCallerSkip1.Error(args...)
}

// Errorf 格式化输出日志，用于输出严重错误信息。
func Errorf(template string, args ...interface{}) {
	zaplogCallerSkip1.Errorf(template, args...)
}

// Fatal 输出日志并结束程序，用于再不可恢复的情况下输出错误日志。
func Fatal(args ...interface{}) {
	zaplogCallerSkip1.Fatal(args...)
}

// Fatalf 格式化输出日志并结束程序，用于再不可恢复的情况下输出错误日志。
func Fatalf(template string, args ...interface{}) {
	zaplogCallerSkip1.Fatalf(template, args...)
}
