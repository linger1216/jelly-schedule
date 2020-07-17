package core

import (
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	l      *zap.SugaredLogger
)

func init() {
	logger, _ = zap.NewDevelopment()
	l = logger.Sugar()
}
