package core

type WorkFlow struct {
	Id          string
	Name        string
	description string
	JobIds      []string
	CreateTime  int64
	UpdateTime  int64
}
