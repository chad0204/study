package main

import (
	"fmt"
	"net"
)

// 进入文件路径
// .\server.exe
// go run .\server.go
func main() {
	fmt.Println("Starting the server to listen localhost 50000...")
	listen, err := net.Listen("tcp", "localhost:50000")
	if err != nil {
		fmt.Printf("Error listening, err:[%v] \n", err)
		return
	}
	//监听客户端
	for {
		accept, err := listen.Accept()
		if err != nil {
			fmt.Printf("Error accepting, err:[%v] \n", err)
			return // 终止程序
		}
		fmt.Printf("Accepting client connect, localAddr: [%v], remoteAddr: [%v] \n",
			accept.LocalAddr(),
			accept.RemoteAddr())
		// 为每一个客户端产生一个协程用来处理请求
		go connectHandler(accept)
	}

}

func connectHandler(accept net.Conn) {
	for {
		buf := make([]byte, 512)
		l, err := accept.Read(buf)
		if err != nil {
			fmt.Printf("Error reading, err:[%v] \n", err)
			return
		}
		fmt.Printf("Received data: %v \n", string(buf[:l]))
	}
}
