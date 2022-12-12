package api

import "net/rpc"

// 标记唯一名 , api不能放在main包,
const HelloServiceName = "src/main/demo/rpc/api/api.HelloService"
const KvStoreServiceName = "src/main/demo/rpc/api/api.KvStoreService"

type HelloService interface {
	Hello(request string, replay *string) error
}

func RegisterHelloService(svc HelloService) error {
	return rpc.RegisterName(HelloServiceName, svc)
}

type KvStoreService interface {
	// kv 就是k和v的数组. reply 空返回值(go 内置rpc规则要有个reply)
	Set(kv [2]string, reply *struct{}) error
	Get(key string, value *string) error
	Watch(timeout int, keyChanged *string) error
}

func RegisterKvStoreService(svc KvStoreService) error {
	return rpc.RegisterName(KvStoreServiceName, svc)
}
