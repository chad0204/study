package main

import (
	"sync"
	"time"
)

type Singleton struct {
	name string
}

var (
	once     sync.Once
	INSTANCE *Singleton
	lock     sync.RWMutex
)

func getInstanceUnsafe() *Singleton {
	if INSTANCE == nil {
		INSTANCE = &Singleton{"2333"}
	}
	return INSTANCE
}

func getInstanceLock() *Singleton {
	//读和读不互斥
	lock.RLock()
	//读写互斥, 保证不会发生指令重排, 返回初始化不完全的对象
	if INSTANCE != nil {
		lock.RUnlock()
		return INSTANCE
	}
	lock.RUnlock()

	lock.Lock()
	//再次判断防止重复初始化
	if INSTANCE == nil {
		INSTANCE = &Singleton{name: "2333"}
	}
	lock.Unlock()

	return INSTANCE
}

func getInstanceOnce() *Singleton {
	once.Do(func() {
		INSTANCE = &Singleton{name: "2333"}
	})
	return INSTANCE
}

func main() {
	go getInstanceOnce()
	go getInstanceOnce()

	time.Sleep(1e9)
}
