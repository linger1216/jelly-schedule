package main

//
//import (
//	"context"
//	"fmt"
//	jsoniter "github.com/json-iterator/go"
//	"github.com/linger1216/jelly-schedule/core"
//	"github.com/linger1216/jelly-schedule/utils"
//	"gopkg.in/alecthomas/kingpin.v2"
//	"os/exec"
//)
//
//import _ "net/http/pprof"
//
//const DefaultConfigFilename = "/etc/config/schedule_config.yaml"
//
//var (
//	configFilename = kingpin.Flag("conf", "config file name").Short('c').Default(DefaultConfigFilename).String()
//)
//
//type ShellJob struct {
//}
//
//func NewShellJob() *ShellJob {
//	return &ShellJob{}
//}
//
//func (e *ShellJob) Name() string {
//	return "ShellJob"
//}
//
//func (e *ShellJob) Exec(ctx context.Context, req interface{}) (interface{}, error) {
//	cmdBuf, ok := req.(string)
//	if !ok {
//		return nil, fmt.Errorf("shell is not string")
//	}
//	fmt.Printf("shell:%s\n", cmdBuf)
//
//	cmds := make([]string, 0)
//	err := jsoniter.ConfigFastest.Unmarshal([]byte(cmdBuf), &cmds)
//	if err != nil {
//		return nil, err
//	}
//
//	var resp []byte
//	for _, cmd := range cmds {
//		command := exec.Command("/bin/sh", "-c", cmd)
//		resp, err = command.Output()
//		if err != nil {
//			return nil, err
//		}
//	}
//	return string(resp), nil
//}
//
//func init() {
//	kingpin.Version("0.1.0")
//	kingpin.Parse()
//}
//
//func main() {
//	fmt.Printf("rpc 参数是[]string的json序列化\n")
//	config, err := core.LoadScheduleConfig(*configFilename)
//	if err != nil {
//		panic(err)
//	}
//	end := make(chan error)
//	etcd := core.NewEtcd(&config.Etcd)
//	core.NewJobServer(etcd, NewShellJob())
//	go utils.InterruptHandler(end)
//	<-end
//}
