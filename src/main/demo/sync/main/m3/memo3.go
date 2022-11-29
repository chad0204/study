package m3

import "sync"

type result struct {
	err   error
	value interface{}
}

type Func func(key string) (interface{}, error)

type Memo struct {
	f     Func
	mu    sync.Mutex
	cache map[string]result
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

func (mem *Memo) Get(key string) (interface{}, error) {
	mem.mu.Lock()
	res, ok := mem.cache[key]
	mem.mu.Unlock()
	if !ok {
		//repeat call
		res.value, res.err = mem.f(key)

		mem.mu.Lock()
		mem.cache[key] = res
		mem.mu.Unlock()
	}
	return res.value, res.err
}
