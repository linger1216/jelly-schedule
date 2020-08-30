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
	//logger, _ = zap.NewProduction()
	l = logger.Sugar().With(ProjectKey, ProjectValue)
}

// func (s *SugaredLogger) Debugf(template string, args ...interface{}) {
func WithModule(name string) *zap.SugaredLogger {
	return l.With(ModuleKey, name)
}

func WithWorflow(name string) *zap.SugaredLogger {
	return l.With(WorkFlowKey, name)
}

func WithJob(name string) *zap.SugaredLogger {
	return l.With(JobKey, name)
}
