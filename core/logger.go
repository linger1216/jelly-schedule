package core

import (
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	l      *zap.SugaredLogger
)

const (
	// logger const field
	ProjectKey   = "project"
	ProjectValue = "schedule"

	LM = "module"
	LE = "executor"
	LW = "workflow"
	LJ = "job"
)

func init() {
	logger, _ = zap.NewDevelopment()
	//logger, _ = zap.NewProduction()
	l = logger.Sugar().With(ProjectKey, ProjectValue)
}

// func (s *SugaredLogger) Debugf(template string, args ...interface{}) {
func _M(name string) *zap.SugaredLogger {
	return l.With(LM, name)
}

func _W(name string) *zap.SugaredLogger {
	return l.With(LW, name)
}

func _J(name string) *zap.SugaredLogger {
	return l.With(LJ, name)
}
