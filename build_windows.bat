set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -gcflags=-m -ldflags="-w -s" -o tmp/docg.exe main.go
upx tmp/docg.exe