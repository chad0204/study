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

func (hello *HelloServiceClient) hello(args string, reply *string) error {
	return hello.Client.Call(api.HelloServiceName+".Hello", args, reply)
}

func (hello *HelloServiceClient) helloAsync(args string, reply *string) {
	//异步调用
	call := hello.Client.Go(api.HelloServiceName+".Hello", args, reply, nil) //nil

	// do something else

	call = <-call.Done
	if err := call.Error; err != nil {
		log.Fatal(err)
	}
	a := call.Args.(string)
	r := call.Reply.(*string)
	fmt.Println(a, *r)
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

	service.helloAsync("pcpcpc233", &reply)

}
