package main

/*
g1 -> requests
g2 -> requests
g3 -> requests

monitor <- requests

cache 只有monitor可以访问
*/
type result struct {
	err   error
	value interface{}
}

type entry struct {
	ready chan struct{}
	res   result
}

// 封装调用Get的每一次请求
type request struct {
	response chan<- result
	key      string
}

// Memo requests作为Get函数和monitor goroutine的通道
type Memo struct {
	requests chan request
}

type Func func(key string) (interface{}, error)

func New(f Func) *Memo {
	requests := make(chan request)
	m := &Memo{requests: requests}
	go m.server(f)
	return m
}

// Get 遍历所有的请求进行处理
func (mem *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range mem.requests {
		e := cache[req.key]
		if e == nil {
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			//异步处理, 让其他请求进来. call和deliver通过ready通道同步
			go e.call(f, req.key)

		}
		go e.deliver(req.response)
	}
}

func (mem *Memo) Close() {
	close(mem.requests)
}

func (e *entry) deliver(response chan<- result) {
	//阻塞到cache写入完成
	<-e.ready
	//发送给response
	response <- e.res
}

func (e *entry) call(f Func, key string) {
	e.res.value, e.res.err = f(key)
	close(e.ready)
}

func (mem *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	//将请求发送给monitor goroutine
	mem.requests <- request{response: response, key: key}
	//阻塞等待结果
	res := <-response
	return res.value, res.err
}
