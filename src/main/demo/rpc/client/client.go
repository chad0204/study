package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"study/src/main/demo/rpc/api"
	"time"
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

func callHello() {
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

func Watch(timeout int, client *rpc.Client) {
	var keyChanged string
	//阻塞调用
	err := client.Call(api.KvStoreServiceName+".Watch", timeout, &keyChanged)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("watch:", keyChanged)
}

func Set(k string, v string, client *rpc.Client) {
	err := client.Call(api.KvStoreServiceName+".Set", [2]string{k, v}, new(struct{}))
	if err != nil {
		log.Fatal(err)
	}
}

func Get(k string, client *rpc.Client) {
	var reply string
	err := client.Call(api.KvStoreServiceName+".Get", k, &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("result: %v \n", reply)
}

func main() {
	conn, e := net.Dial("tcp", "localhost:2333")
	if e != nil {
		return
	}
	client := rpc.NewClient(conn)

	//注册watch

	for i := 0; i < 100; i++ {
		go func() {
			Watch(100, client)
		}()
	}

	time.Sleep(time.Duration(5) * time.Second)

	//for i := 0; i < 100; i++ {
	Set("key", strconv.Itoa(235), client)
	//}

	Get("key", client)

	time.Sleep(time.Duration(1000) * time.Second)

}
