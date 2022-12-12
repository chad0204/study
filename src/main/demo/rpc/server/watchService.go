package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"study/src/main/demo/rpc/api"
	"sync"
	"time"
)

type DefaultKvStoreService struct {
	cache   map[string]string
	filters map[string]func(key string) // 监听方法的过滤函数,
	mu      sync.Mutex
}

func NewDefaultKvStoreService() *DefaultKvStoreService {
	return &DefaultKvStoreService{
		cache:   make(map[string]string),
		filters: make(map[string]func(key string)),
	}
}

func (p *DefaultKvStoreService) Set(kv [2]string, reply *struct{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	key, value := kv[0], kv[1]

	if oldValue := p.cache[key]; oldValue != value {
		// value有变动, 执行所有的filter函数
		for _, filter := range p.filters {
			filter(key)
		}
	}
	p.cache[key] = value
	return nil
}

func (p *DefaultKvStoreService) Get(key string, value *string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if v, ok := p.cache[key]; ok {
		*value = v
		return nil
	}
	return fmt.Errorf("not found key: %v", key)
}

func (p *DefaultKvStoreService) Watch(timeout int, keyChanged *string) error {
	watchId := fmt.Sprintf("watch-%s-%03d", time.Now(), rand.Int())

	calls := make(chan string, 100) // buffer 防止filter导致Set方法阻塞

	fmt.Printf("filter set before: %v \n", watchId)
	p.mu.Lock()
	p.filters[watchId] = func(key string) {
		calls <- key
	}
	p.mu.Unlock()
	fmt.Printf("filter set after: %v \n", watchId)

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		fmt.Printf("timeout")
	case key := <-calls:
		fmt.Printf("select key changed: %v \n", key)
		*keyChanged = key
		return nil
	}
	return nil
}

func main() {
	api.RegisterKvStoreService(NewDefaultKvStoreService())

	//也可以使用http协议, http.HandleFunc
	listen, err := net.Listen("tcp", ":2333")
	if err != nil {
		log.Fatal("ListenTCP error:", err) // 打印并退出
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal("Accept  error:", err) // 打印并退出
		}
		go rpc.ServeConn(conn)
	}
}
