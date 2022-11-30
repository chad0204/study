package main

import "sync"

type entry struct {
	res   result
	ready chan struct{} //标记是否写入完成
}

type result struct {
	err   error
	value interface{}
}

type Func func(key string) (interface{}, error)

type Memo struct {
	f     Func
	mu    sync.Mutex
	cache map[string]*entry
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]*entry)}
}

func (mem *Memo) Get(key string) (interface{}, error) {
	mem.mu.Lock()
	e, ok := mem.cache[key]
	if !ok {
		//不管三七二十一, 先把entry放进去, 尽早是否锁。其他goroutine会走下面的分支
		e = &entry{ready: make(chan struct{})}
		mem.cache[key] = e
		mem.mu.Unlock()

		//开始call
		e.res.value, e.res.err = mem.f(key)

		close(e.ready)
	} else {
		//说明已经有其他goroutine在写相同的key, 直接释放锁, 让其他读请求进来
		mem.mu.Unlock()
		//阻塞到ready 关闭, 写结束。 这里会多次从关闭通道中接收, 无所谓
		<-e.ready
	}
	return e.res.value, e.res.err
}
