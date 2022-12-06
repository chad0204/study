package effective

import (
	"context"
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

//基于once的单例
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

func generateNums(ctx context.Context) <-chan int {
	nums := make(chan int)
	go func() {
		for i := 2; ; i++ {
			select {
			case nums <- i:
			case <-ctx.Done(): //通知结束
				fmt.Println("generate done")
				return
			}
		}
	}()
	return nums
}

func filterPrime(nums <-chan int, prime int, ctx context.Context) <-chan int {
	filters := make(chan int)
	go func() {
		for {
			if i := <-nums; i%prime != 0 {
				select {
				case <-ctx.Done():
					fmt.Println("filter done")
					return
				case filters <- i:
				}
			}
		}
	}()
	return filters
}

//打印素数
func TestPrime(t *testing.T) {

	ctx, cancelFunc := context.WithCancel(context.Background())

	nums := generateNums(ctx)

	//生成100个素数
	for i := 0; i < 100; i++ {
		prime := <-nums
		fmt.Println(prime)
		nums = filterPrime(nums, prime, ctx)
	}
	time.Sleep(10e9) // 这段时间 goroutine没有释放
	fmt.Println("start cancel goroutine")
	cancelFunc()
}

// select实现随机数
func TestSelectForRandom(t *testing.T) {
	random := make(chan int)

	//0到4的随机数
	go func() {
		for {
			//select会随机选择一个可以操作的chan进行操作
			select {
			case random <- 0:
			case random <- 1:
			case random <- 2:
			case random <- 3:
			case random <- 4:
			}
		}
	}()

	for rdm := range random {
		fmt.Println(rdm)
	}
}

func worker(wg *sync.WaitGroup, cancel chan bool) {
	defer wg.Done()
	for {
		select {
		default:
			//do something
			fmt.Println(".")
		case <-cancel:
			fmt.Println("cancel")
			return
		}
	}
}

// 使用select + close(chan)实现广播关闭goroutine （go没有提供关闭goroutine的方法）
func TestSelectForClose(t *testing.T) {
	cancel := make(chan bool)
	var wg sync.WaitGroup

	//如果要关闭多个goroutine, 不需要多个chan, 可以通过close(chan)来做到
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go worker(&wg, cancel)
	}

	time.Sleep(1 * time.Second)
	//cancel <- true
	close(cancel)

	//time.Sleep(10e9) //这个时间不好确定, 改用wg
	wg.Wait()
}

func workerCTX(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	for {
		select {
		default:
			//do something
			fmt.Println(".")
		case <-ctx.Done():
			fmt.Println("done")
			return ctx.Err()
		}
	}
}

func TestContextForClose(t *testing.T) {
	var wg sync.WaitGroup

	timeout, cancelFunc := context.WithTimeout(context.Background(), 10e9)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go workerCTX(&wg, timeout)
	}

	time.Sleep(1 * time.Second)
	//cancel <- true
	cancelFunc()

	wg.Wait()
}
