package api

import "net/rpc"

//标记唯一名 , api不能放在main包,
const HelloServiceName = "src/main/demo/rpc/api/api."

type HelloService interface {
	Hello(request string, replay *string) error
}

func RegisterHelloService(svc HelloService) error {
	return rpc.RegisterName(HelloServiceName, svc)
}
