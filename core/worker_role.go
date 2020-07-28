package core

type WorkerRole uint8

const (
	Wait WorkerRole = iota
	Follower
	Leader
)

func getWorkerRoleDescription(role WorkerRole) string {
	switch role {
	case Follower:
		return "Follower"
	case Leader:
		return "Leader"
	}
	return "Wait"
}
