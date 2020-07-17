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
	TTL              = int64(10)
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
	ret.role = Follower
	l.Debugf("%s", ret.String())

	err := ret.register()
	if err != nil {
		panic(err)
	}
	err = ret.watchVote()
	if err != nil {
		panic(err)
	}

	ret.startElection()
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
	}
}

func (w *Worker) watchVote() error {
	return w.discover.WatchWithPrefix(EtcdLeaderKey, func(event *clientv3.Event) {
		switch event.Type {
		case mvccpb.PUT:
			if event.IsCreate() {
				l.Debugf("leader create:%s", string(event.Kv.Value))
			} else {
				l.Debugf("leader update:%s", string(event.Kv.Value))
			}
		case mvccpb.DELETE:
			l.Debugf("delete leader")
			time.Sleep(time.Second)
			l.Debugf("prepare election leader")
			w.startElection()
		}
	})
}

func (w *Worker) String() string {
	return fmt.Sprintf("host:%s dir:%s role:%s", w.name, w.dir, getWorkerRoleDescription(w.role))
}
