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
		//g1 ok是false, 调用f后正在设置map, ok是true, 但value还没设置。g2 ok是true, 但是value是nil
		mem.cache[key] = res
	}
	return res.value, res.err
}
