env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/bin/api cmd/api/api.go
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/bin/executor cmd/executor/executor.go
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/bin/echo-job example/echo-job/main.go
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/bin/shell-job example/shell-job/main.go
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/bin/http-job example/http-job/main.go