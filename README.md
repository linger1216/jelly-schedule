# jelly-schedule

### linux
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api_x cmd/api/main.go

### pprof
go tool pprof -http=:6060 --seconds 30 http://localhost:6060

go build -o api cmd/api/main.go && ./api --etcd 172.3.0.122:2379 --port 23808


## 设计概念
### 调度:
- cron (定时调度)
- api (user自定义触发)

### trggier
- cron
- replay


### job:
代表一个独立逻辑
没有独立调度时间

- rpc job(需要编程)
- shell job 
- http job


### 工作流: workflow
支持数据的传递
工作流代表一个完整的任务, 每个任务,都有由1个或多个job组成


### 编程模型
[A,B,C]  A->B->C 顺序执行

[["a"],["a","b","c"], ["x.y.z"]]



replace go.etcd.io/etcd => github.com/etcd-io/etcd v3.3.22+incompatible

参考:
1. https://github.com/betterde/ects
2. https://github.com/busgo/forest



