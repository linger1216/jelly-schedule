package core

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasttemplate"
	"os"
	"time"
)

var (
	WorkerPrefix     = `/schedule/worker`
	EtcdWorkerFormat = fasttemplate.New(WorkerPrefix+`/{Name}`, "{", "}")
	EtcdLeaderKey    = `/schedule/leader`
	TTL              = int64(10)
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

func MarshalLeader(w *WorkerStats) ([]byte, error) {
	return []byte(w.Name), nil
}

func UnMarshalLeader(buf []byte) (*WorkerStats, error) {
	s := &WorkerStats{}
	s.Name = string(buf)
	return s, nil
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
	stats         *WorkerStats
	etcd          *Etcd
	leaderLeaseId clientv3.LeaseID
	roleLeaseId   clientv3.LeaseID
	ticker        *time.Ticker
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
	err := ret.tryElection()
	if err != nil {
		panic(err)
	}

	// 监视leader, 如果发现leader不在了, 立即发起选举
	err = ret.watchLeader()
	if err != nil {
		panic(err)
	}

	// 自身也作为worker节点
	err = ret.register()
	if err != nil {
		panic(err)
	}

	// 续租
	// run ticker by TTL/2
	// 1. keep alive lease if leader
	// 2. keep alive lease if worker
	// 3. update stats if worker
	ticker := time.NewTicker(time.Duration(TTL/2) * time.Second)
	ret.ticker = ticker
	go ret.handleTicker()
	return ret
}

func (w *Worker) handleTicker() {
	for {
		select {
		case <-w.ticker.C:
			// worker register
			err := w.register()
			if err != nil {
				l.Debugf("register err:%s", err.Error())
			}
			// keep alive lease if worker
			if err := w.etcd.RenewLease(context.Background(), w.roleLeaseId); err != nil {
				l.Debugf("renew worker lease err:%s", err.Error())
			} else {
				l.Debugf("renew worker lease ok")
			}
			// keep alive lease if leader
			if w.leaderLeaseId > 0 {
				if err := w.etcd.RenewLease(context.Background(), w.leaderLeaseId); err != nil {
					l.Debugf("renew leader lease err:%s", err.Error())
				} else {
					l.Debugf("renew leader lease ok")
				}
			}
		}
	}
}

func (w *Worker) register() error {
	if w.roleLeaseId == 0 {
		roleLeaseId, err := w.etcd.GrantLease(TTL)
		if err != nil {
			return err
		}
		w.roleLeaseId = roleLeaseId
	}
	jsonBuf, _ := MarshalWorker(w.stats)
	return w.etcd.InsertKV(context.Background(), WorkerKey(w.stats.Name), string(jsonBuf), w.roleLeaseId)
}

func (w *Worker) tryElection() error {
	if w.stats.Role == Leader {
		return nil
	}

	// 不是leader如果id
	if w.leaderLeaseId > 0 {
		_, err := w.etcd.RevokeLease(w.leaderLeaseId)
		if err != nil {
			return err
		}
		w.leaderLeaseId = 0
	}

	if w.leaderLeaseId == 0 {
		leaderId, err := w.etcd.GrantLease(TTL)
		if err != nil {
			return err
		}
		w.leaderLeaseId = leaderId
	}

	jsonBuf, _ := MarshalLeader(w.stats)
	err := w.etcd.InsertKVNoExisted(context.Background(), LeaderKey(), string(jsonBuf), w.leaderLeaseId)
	if err != nil {
		// 发生了错误, 无论是真实的错误, 还是已经选举出来了leader
		// 无论如何, 为此准备的lease不起作用了, 随即释放掉
		_, _ = w.etcd.RevokeLease(w.leaderLeaseId)
		w.leaderLeaseId = 0
		if err == ErrKeyAlreadyExists {
			//l.Debugf("try election leader already exists")
			w.changeRole(Follower)
			return nil
		} else {
			//l.Debugf("try election err:%s", err.Error())
			return err
		}
	} else {
		//l.Debugf("try election be leader")
		w.changeRole(Leader)
	}

	if w.stats.Role != Leader {
		_, err := w.etcd.RevokeLease(w.leaderLeaseId)
		if err != nil {
			return err
		}
		w.leaderLeaseId = 0
	}
	return nil
}

func (w *Worker) watchLeader() error {
	return w.etcd.WatchWithPrefix(EtcdLeaderKey, func(event *clientv3.Event) error {
		switch event.Type {
		case mvccpb.PUT:
			ws, err := UnMarshalLeader(event.Kv.Value)
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
			_ = w.tryElection()
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
