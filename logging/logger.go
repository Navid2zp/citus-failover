package logging

import (
	"go.uber.org/zap"
)

type Logger struct {
	service string
	l       *zap.Logger
}

func NewLogger(service string) *Logger {
	//zap.NewProductionConfig()
	//cfg := zap.NewProductionConfig()
	//cfg.OutputPaths = []string{"stdout"}
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return &Logger{service, logger}
}
