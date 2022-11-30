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

// Get 写和读map互斥, 问题是并发获取相同key时无法防止重复调用, 虽然覆盖不影响结果, 但是导致资源浪费。
func (mem *Memo) Get(key string) (interface{}, error) {
	mem.mu.Lock()
	res, ok := mem.cache[key]
	mem.mu.Unlock()
	if !ok {
		//repeat call
		//duplicate suppression（重复抑制/避免）
		res.value, res.err = mem.f(key)

		mem.mu.Lock()
		mem.cache[key] = res
		mem.mu.Unlock()
	}
	return res.value, res.err
}
