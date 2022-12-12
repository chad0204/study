package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"study/src/main/demo/rpc/api"
	"sync"
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

func callWatch() {
	conn, e := net.Dial("tcp", "localhost:2333")
	if e != nil {
		return
	}
	client := rpc.NewClient(conn)

	var wg = sync.WaitGroup{}

	go func() {
		wg.Add(1)
		for i := 0; i < 10; i++ {
			var keyChanged string
			err := client.Call(api.KvStoreServiceName+".Watch", 3, &keyChanged)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("watch:", keyChanged)
		}
		wg.Done()
	}()

	time.Sleep(1e9)

	go func() {
		wg.Add(1)
		for i := 0; i < 10; i++ {
			e = client.Call(api.KvStoreServiceName+".Set", [2]string{"abc", strconv.Itoa(i)}, new(struct{}))
			if e != nil {
				log.Fatal(e)
			}
		}
		wg.Done()
	}()

	wg.Wait()

	var reply string
	e = client.Call(api.KvStoreServiceName+".Get", "abc", &reply)
	if e != nil {
		log.Fatal(e)
	}
	fmt.Println(reply)

}

func main() {

	callWatch()
}
