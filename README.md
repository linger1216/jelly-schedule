# jelly-schedule


这是一个轮子, 公司定时任务太多, 太繁杂, 需要一个轮子来管理

目前正在开发中ing

v0.1.0 基本版本

### Todo
- 每个任务增加是否立即实行,以及delay的时间参数 毕竟cron是要等下一个时间周期的
- 基于name的选择job执行(目前是基于id), 考虑到打洞
- 完善文档 (可能永远都不会完成)
- 对参数的强制限定
- k8s支持
- jobId 不存在的时候, 任务执行失败


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
a -> b -> c 顺序执行 = a and b and c 
a | b | c   并发执行 = a or  b or  c 



replace go.etcd.io/etcd => github.com/etcd-io/etcd v3.3.22+incompatible

参考:
1. https://github.com/betterde/ects
2. https://github.com/busgo/forest



