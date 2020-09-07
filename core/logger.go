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

	_Mod          = "module"
	_Exec         = "executor"
	_Workflow     = "workflow"
	_Expr         = "expr"
	_Job          = "job"
	_AlternateJob = "AlternateJob"
	_SerialJob    = "SerialJob"
	_ParallelJob  = "ParallelJob"
)

func init() {
	logger, _ = zap.NewDevelopment()
	//logger, _ = zap.NewProduction()
	l = logger.Sugar().With(ProjectKey, ProjectValue)
}

// func (s *SugaredLogger) Debugf(template string, args ...interface{}) {
func _MOD(name string) *zap.SugaredLogger {
	return l.With(_Mod, name)
}

func _W(name string) *zap.SugaredLogger {
	return l.With(_Workflow, name)
}

func _J(name string) *zap.SugaredLogger {
	return l.With(_Job, name)
}
