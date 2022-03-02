set GOOS=darwin
set GOARCH=amd64
go build -gcflags=-m -ldflags="-w -s" -o tmp/docg.darwin main.go
upx tmp/docg.darwin

