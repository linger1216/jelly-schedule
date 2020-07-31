go build -o echo-job example/echo-job/server/main.go && ./echo-job --etcd 172.3.0.122:2379


# ./worker --etcd 172.3.0.122:2379 --name a1