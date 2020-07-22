package core

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/linger1216/jelly-schedule/etcdv3"
	"os"
	"time"
)
import "github.com/valyala/fasttemplate"

var (
	EtcdWorkerFormat = fasttemplate.New(`/schedule/worker/{name}`, "{", "}")
	EtcdLeaderKey    = `/schedule/leader`
	TTL              = int64(5)
)

type Worker struct {
	name     string
	dir      string
	discover *etcdv3.Etcd
	role     WorkerRole
}

func NewWorker(name string, discover *etcdv3.Etcd) *Worker {
	ret := &Worker{name: name, discover: discover}
	ret.dir, _ = os.Getwd()
	ret.discover = discover

	ret.changeRole(Unknown)
	err := ret.register()
	if err != nil {
		panic(err)
	}

	err = ret.watchElectionLeader()
	if err != nil {
		panic(err)
	}

	ret.startElection()

	go func() {
		for {
			l.Debugf("worker %s", ret.String())
			time.Sleep(time.Second * 3)
		}
	}()
	return ret
}

func (w *Worker) register() error {
	s := EtcdWorkerFormat.ExecuteString(map[string]interface{}{
		"name": w.name,
	})
	err := w.discover.TxKeepaliveWithTTL(s, "true", 10)
	return err
}

func (w *Worker) startElection() {
	err := w.discover.TxKeepaliveWithTTL(EtcdLeaderKey, w.name, TTL)
	if err != nil {
		l.Debugf("startElection err:%s", err.Error())
		// 选举不成功的话,默认为Follower
		// 如果此时已经有了leader, 当worker运行时, 是不会受到put消息的, 所以没法在watch中判定role
		// 刚开始只有通过是否成功竞选才能判定role
		w.changeRole(Follower)
	} else {
		w.changeRole(Leader)
	}
}

func (w *Worker) watchElectionLeader() error {
	return w.discover.WatchWithPrefix(EtcdLeaderKey, func(event *clientv3.Event) {
		switch event.Type {
		case mvccpb.PUT:
			if string(event.Kv.Value) == w.name {
				w.changeRole(Leader)
			} else {
				w.changeRole(Follower)
			}
		case mvccpb.DELETE:
			time.Sleep(time.Second)
			l.Debugf("need election leader again")
			w.startElection()
		}
	})
}

func (w *Worker) changeRole(role WorkerRole) {
	w.role = role
}

func (w *Worker) String() string {
	return fmt.Sprintf("host:%s dir:%s role:%s", w.name, w.dir, getWorkerRoleDescription(w.role))
}
