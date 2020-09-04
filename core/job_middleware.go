package core

type JobMiddlewareRequest struct {
	Name      string
	Parameter []string
}

type SyncRequest struct {
	SrcHost   string
	SrcFiles  []string
	DestHost  string
	DestFiles []string
}
