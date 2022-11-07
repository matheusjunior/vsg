package logger

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func init() {
	l, _ := zap.NewProduction()
	logger = l.Sugar()
}

func Sync() {
	logger.Sync()
}

func Logger() *zap.SugaredLogger {
	return logger
}
