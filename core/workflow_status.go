package core

//
//// 主要记录workflow runtime 状态记录的
//
type WorkFlowStatus struct {
	Id                  string
	Executing           bool
	MaxExecuteCount     int64
	SuccessExecuteCount int64
	FailedExecuteCount  int64
	LastExecuteDuration int64
}

//
//type WorkFlowStatusCommandQueue struct {
//	in chan<- *WorkFlowStatus
//}
//
//func newWorkFlowStatusCommandQueue() *WorkFlowStatusCommandQueue {
//	ret := &WorkFlowStatusCommandQueue{in: make(chan<- *WorkFlowStatus, 1024)}
//	go ret.start()
//	return ret
//}
//
//func (w *WorkFlowStatusCommandQueue) append(cmd *WorkFlowStatus) {
//	w.in <- cmd
//}
//
//func (w *WorkFlowStatusCommandQueue) start() {
//	for cmd := range w.in {
//		l.Debugf("%v", cmd)
//	}
//}
