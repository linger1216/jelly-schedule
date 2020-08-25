package core

import "sync"

type SyncMap struct {
	m map[string]interface{}
	sync.RWMutex
}

func NewSyncMap() *SyncMap {
	return &SyncMap{m: make(map[string]interface{})}
}

func (s *SyncMap) Put(k string, v interface{}) {
	s.Lock()
	s.m[k] = v
	s.Unlock()
}

func (s *SyncMap) Get(k string) (interface{}, bool) {
	s.RLock()
	defer s.RUnlock()
	resp, ok := s.m[k]
	return resp, ok
}
