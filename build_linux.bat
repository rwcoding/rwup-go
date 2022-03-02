set GOOS=linux
set GOARCH=amd64
go build -gcflags=-m -ldflags="-w -s" -o tmp/docg.linux main.go
upx tmp/docg.linux
