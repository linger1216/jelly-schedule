
rm -fr build
mkdir -p build/bin

go build -o build/bin/api cmd/api/main.go
go build -o build/bin/echo-job example/echo-job/server