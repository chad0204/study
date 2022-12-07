package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"study/src/main/demo/rpc/api"
)

type HelloServiceClient struct {
	*rpc.Client
}

func DailHelloService(network, address string) (*HelloServiceClient, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	clientCodec := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	return &HelloServiceClient{clientCodec}, nil
}

func (hello *HelloServiceClient) hello(request string, replay *string) error {
	return hello.Client.Call(api.HelloServiceName+".Hello", request, replay)
}

func main() {
	service, err := DailHelloService("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Dial error: ", err)
	}

	var reply string
	err = service.hello("pc", &reply)
	if err != nil {
		log.Fatal("Call error: ", err)
	}

	fmt.Println(reply)
}
