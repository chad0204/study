package main

import "sync"

type result struct {
	err   error
	value interface{}
}

type Memo struct {
	f     Func
	mu    sync.Mutex
	cache map[string]result
}

type Func func(key string) (interface{}, error)

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

// Get 相当于串行
func (mem *Memo) Get(key string) (interface{}, error) {
	mem.mu.Lock()
	res, ok := mem.cache[key]
	if !ok {
		res.value, res.err = mem.f(key)
		mem.cache[key] = res
	}
	mem.mu.Unlock()
	return res.value, res.err
}
