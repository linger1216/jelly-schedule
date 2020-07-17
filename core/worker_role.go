package core

type WorkerRole uint8

const (
	Unknown WorkerRole = iota
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
	return "Unknown"
}
