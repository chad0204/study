go导包
首先生成go.mod文件。手动import相应的包。然后执行go mod tidy（也可以go get下载）
go build和go run生成可执行文件不在bin包，go install 生成可执行文件直接放入bin包