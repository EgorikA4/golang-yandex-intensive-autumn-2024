package logger

import (
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func InitLogger() (*zap.Logger, error) {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		return &zap.Logger{}, err
	}
	return logger, err
}

func GetLogger() *zap.Logger {
	return logger
}
