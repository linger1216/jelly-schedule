package core

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasttemplate"
	"log"
	"os"
	"time"
)

var (
	WorkerPrefix     = `/schedule/worker`
	EtcdWorkerFormat = fasttemplate.New(WorkerPrefix+`/{Name}`, "{", "}")
	EtcdLeaderKey    = `/schedule/leader`
	TTL              = int64(5)
)

type WorkerStats struct {
	Name string     `json:"name"`
	Path string     `json:"path"`
	Role WorkerRole `json:"role"`
}

func (w WorkerStats) String() string {
	return fmt.Sprintf("host:%s Path:%s Role:%s", w.Name, w.Path, getWorkerRoleDescription(w.Role))
}

func WorkerKey(name string) string {
	s := EtcdWorkerFormat.ExecuteString(map[string]interface{}{
		"Name": name,
	})
	return s
}

func LeaderKey() string {
	return EtcdLeaderKey
}

func MarshalWorker(w *WorkerStats) ([]byte, error) {
	return jsoniter.ConfigFastest.Marshal(w)
}

func UnMarshalWorker(buf []byte) (*WorkerStats, error) {
	s := &WorkerStats{}
	err := jsoniter.ConfigFastest.Unmarshal(buf, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

type Worker struct {
	stats    *WorkerStats
	etcd     *Etcd
	leaderId clientv3.LeaseID
	ticker   *time.Ticker
}

func NewWorker(name string, discover *Etcd) *Worker {
	ret := &Worker{stats: &WorkerStats{}}
	ret.stats.Name = name
	ret.stats.Path, _ = os.Getwd()
	ret.etcd = discover
	ret.changeRole(Wait)

	// 一上来就发起选举, 万一中了呢?
	// 当前监视leader的动作还没有产生, 这时候主动发起一次选举
	// 给自己确定角色
	err := ret.startElection()
	if err != nil {
		panic(err)
	}

	// 监视leader, 如果发现leader不在了, 立即发起选举
	err = ret.watchLeader()
	if err != nil {
		panic(err)
	}

	// run ticker by TTL/2
	// 1. keep alive lease if leader
	// 2. keep alive lease if worker
	// 3. update stats if worker
	ticker := time.NewTicker(time.Duration(TTL/2) * time.Second)

	go func() {
		for {
			l.Debugf("worker %s", ret.Stats())
			time.Sleep(time.Second * 3)
		}
	}()
	return ret
}

func (w *Worker) xxx() {
	for {
		select {
		case <-w.ticker.C:
		default:
		}
	}
}

func (w *Worker) startElection() error {
	if w.leaderId > 0 {
		_, err := w.etcd.RevokeLease(w.leaderId)
		if err != nil {
			return err
		}
		w.leaderId = 0
	}

	leaderId, err := w.etcd.GrantLease(TTL)
	if err != nil {
		return err
	}

	jsonBuf, _ := MarshalWorker(w.stats)
	err = w.etcd.InsertKV(context.Background(), LeaderKey(), string(jsonBuf), leaderId)
	if err != nil {
		// 发生了错误, 无论是真实的错误, 还是已经选举出来了leader
		// 无论如何, 为此准备的lease不起作用了, 随即释放掉
		_, _ = w.etcd.RevokeLease(leaderId)
		if err == ErrKeyAlreadyExists {
			// 已经有了leader
			w.changeRole(Follower)
			return nil
		} else {
			// 发生了错误
			return err
		}
	} else {
		// 保留leader Id
		w.leaderId = leaderId
		w.changeRole(Leader)
	}
	return nil
}

func (w *Worker) watchLeader() error {
	return w.etcd.WatchWithPrefix(EtcdLeaderKey, func(event *clientv3.Event) error {
		switch event.Type {
		case mvccpb.PUT:
			ws, err := UnMarshalWorker(event.Kv.Value)
			if err != nil {
				return err
			}
			if ws.Name == w.stats.Name {
				w.changeRole(Leader)
			} else {
				w.changeRole(Follower)
			}
		case mvccpb.DELETE:
			time.Sleep(time.Second)
			l.Debugf("need election leader again")
			_ = w.startElection()
		}
		return nil
	})
}

func (w *Worker) changeRole(role WorkerRole) {
	if w.stats.Role != role {
		w.stats.Role = role
	}
}

func (w *Worker) Stats() string {
	return w.stats.String()
}

func (w *Worker) Close() {
	w.ticker.Stop()
}
