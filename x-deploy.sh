env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o worker main.go
scp worker root@172.3.0.122:/root
scp worker root@172.3.0.153:/root
scp worker root@172.2.0.21:/root
scp worker root@172.3.0.103:/root