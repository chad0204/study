package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"study/src/main/demo/rpc/api"
)

type DefaultHelloService struct{}

// 规则: 方法只能有两个可序列化的参数，其中第二个参数是指针类型，并且返回一个error类型，同时必须是公开的方法
func (helloService *DefaultHelloService) Hello(request string, reply *string) error {
	*reply = "hello " + request
	return nil
}

func main_____() {
	api.RegisterHelloService(new(DefaultHelloService))

	//也可以使用http协议, http.HandleFunc
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err) // 打印并退出
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal("Accept  error:", err) // 打印并退出
		}
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
