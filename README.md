# jelly-schedule

### linux
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api_x cmd/api/main.go

### pprof
go tool pprof -http=:6060 --seconds 30 http://localhost:6060

go build -o api cmd/api/main.go && ./api --etcd 172.3.0.122:2379 --port 23808


## 设计概念
### 调度:
- cron (定时调度)
- fixed (cron不支持的调度时间)
> 由于Crontab必须被60整除，如果需要每隔40分钟执行一次调度，则Cron无法支持。Fixed rate专门用来做定期轮询，可以解决该问题，且表达式简单，但不支持秒级别

- api (user自定义触发)
- replay (回放)
> 如果您的业务发生变更，如数据库增加一个字段或者上一个月数据有错误，需要把过去一段时间的任务重新执行一遍，可以重刷调度任务数据。(底层还是委托fake cron来实现)


### trggier
- cron
- ticker
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
- ID
- 名称
- 描述
- cron
- 实例并发数


### 编程模型
- 单机 一个任务实例只会随机触发到一台Worker上
- 广播执行表示一个任务实例会广播到该分组所有Worker上执行，当所有Worker都执行完成，该任务才算完成。任意一台Worker执行失败，都算该任务失败
- map 执行map方法可以把一批子任务分布式到多台机器上执行, 子任务入口需要从context里获取taskName，自己判断是哪个子任务，进行相应的逻辑处理 (相当于任务分散执行了, 然后统一返回)
- mapreduce map后会执行reducer模块, 任务结果会缓存在Master节点，内存压力较大，建议子任务个数和Result不要太多
- 分片模式, 还要想下 todo


### 命名空间
概念
- worker leader
所有worker只有1台

- worker follower
普通节点


replace go.etcd.io/etcd => github.com/etcd-io/etcd v3.3.22+incompatible

参考:
1. https://github.com/betterde/ects
2. https://github.com/busgo/forest


/schedule/worker/ip -> worker node (不使用lease实现, 使用定时器TTL/2, 尝试实现以下)
/schedule/leader -> worker node (使用lease实现)




./api --etcd 172.3.0.122:2379 --postgres "postgres://lid.guan:@localhost:15432/schedule?sslmode=disable" --port 23808

