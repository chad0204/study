package effective

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 1. 互斥锁
var total struct {
	mu    sync.Mutex
	value int
}

func work(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		total.mu.Lock()
		total.value += i
		total.mu.Unlock()
	}
}

// 2. 原子变量
var totalV2 int32

func workAtomic(wg *sync.WaitGroup) {
	defer wg.Done()
	var i int32
	for i = 0; i < 100; i++ {
		atomic.AddInt32(&totalV2, i)
	}
}

func TestTotal(t *testing.T) {
	for {
		var wg sync.WaitGroup
		wg.Add(2)
		go workAtomic(&wg)
		go workAtomic(&wg)
		wg.Wait()
		if totalV2 != 9900 {
			fmt.Println(totalV2)
		}
		totalV2 = 0
	}
}

// 3. channel
var operates = make(chan int)
var res = make(chan int)

func monitor() {
	var t int
	for ope := range operates {
		t += ope
	}
	res <- t
}

func workChan() {
	go monitor()

	d := make(chan struct{})
	go func() {
		for i := 0; i < 100; i++ {
			operates <- i
		}
		d <- struct{}{}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			operates <- i
		}
		d <- struct{}{}
	}()
	<-d
	<-d
	close(operates)
	fmt.Println(<-res)

}

func TestChanCount(t *testing.T) {
	workChan()
}

// 3. 互斥 + 原子变量标志 单例
type singleton struct {
}

var (
	instance    *singleton
	initialized uint32 //用一个无符号4字节整数作为初始化标志
	mux         sync.Mutex
)

func Instance() *singleton {
	if atomic.LoadUint32(&initialized) == 1 {
		return instance
	}

	mux.Lock()
	defer mux.Unlock()

	if instance == nil {
		defer atomic.StoreUint32(&initialized, 1)
		instance = &singleton{}
	}
	return instance
}

// 4. sync.once, 上面的方式是once的实现方式

type Once struct {
	mu   sync.Mutex
	done uint32
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

// 基于once的单例
var one Once

func InstanceV2() *singleton {
	one.Do(func() {
		instance = &singleton{}
	})
	return instance
}

// 一种异常
var a string
var dones = make(chan bool)
var done bool

func setup() {
	a = "hello world"
	dones <- true
	//done = true
	close(dones) //效果一样
}

func TestSetup(t *testing.T) {
	go setup()
	//1. setup结束了, 但是done还在寄存器没感知到, 所以无法退出
	//2. setup指令重排, done为true, 但是a没初始化
	//for !done {}
	<-dones
	fmt.Println(a)
}

//对应两个producer
var over = make(chan struct{})

func producer(factor int, containers chan<- interface{}) {
	for i := 0; ; i++ {
		select {
		case containers <- i * factor:
		case <-over:
			fmt.Println("p over")
			break
		}
	}
}

func consumer(containers <-chan interface{}) {
	for v := range containers {
		fmt.Printf("consume value: %v \n", v)
	}
	fmt.Println(" c over")
}

//producer and consumer
func TestVChannel(t *testing.T) {
	containers := make(chan interface{}, 100)

	go producer(3, containers)
	go producer(5, containers)
	go consumer(containers)

	time.Sleep(1e8)
	over <- struct{}{}
	over <- struct{}{}

}
