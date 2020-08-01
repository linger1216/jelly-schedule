rm -fr build
mkdir -p build/bin

go build -o build/bin/api cmd/api/api.go
go build -o build/bin/executor cmd/executor/executor.go
go build -o build/bin/echo-job example/echo-job/server/main.go