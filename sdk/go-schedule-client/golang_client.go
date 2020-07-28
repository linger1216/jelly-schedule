package go_schedule_client

import "github.com/linger1216/jelly-schedule/core"

type GoScheduleClient struct {
	job core.Job
}

func (g *GoScheduleClient) Start() error {
	// 1. 检查本机是不是worker节点, 如果不是拒绝
	// 需要委托本机的restful服务, 使用默认的worker restful http端口
	//

	return nil
}
