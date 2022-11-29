package main

type result struct {
	err   error
	value interface{}
}

type Memo struct {
	f     Func
	cache map[string]result
}

type Func func(key string) (interface{}, error)

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

func (mem *Memo) Get(key string) (interface{}, error) {
	res, ok := mem.cache[key]
	if !ok {
		res.value, res.err = mem.f(key)
		mem.cache[key] = res
	}
	return res.value, res.err
}
