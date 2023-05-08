package effective

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var balance int

func Deposit(amount int) {
	balance = balance + amount
}

func Balance() int {
	return balance
}

func Clear() {
	balance = 0
}

func TestUnSafeBank(t *testing.T) {

	countNum := 0 //执行次数
	errNum := 0   //错误次数
	for {
		Clear() //每次开始需要重新置为0
		countNum++

		//wait保证两个goroutine都执行完成再进行判断
		done := make(chan int)
		//g1
		go func() {
			Deposit(100)
			done <- 1
		}()
		//g2
		go func() {
			Deposit(200)
			done <- 1
		}()

		<-done
		<-done
		if Balance() != 300 {
			/**
			g1 read balance(0) + amount(100)
			g2 read balance(0) + amount(200)
			g2 write balance = 200 //丢失
			g1 write balance = 100

			最终balance不等于300
			*/

			errNum++
			fmt.Printf("balance = %v, count = %v; error = %v \n", Balance(), countNum, errNum)

		}
	}
}

/**-----------------------------------使用channel 进行线程安全改造--------------------------------------------------------*/

var deposits = make(chan int)       //存钱动作
var balances = make(chan int)       //金额
var clearChan = make(chan struct{}) //清空

func DepositSafe(amount int) {
	deposits <- amount
}

func BalanceSafe() int {
	return <-balances
}

func ClearSafe() {
	clearChan <- struct{}{}
}

func TestSafeBank(t *testing.T) {

	//先启动监听
	go func() {
		//balance变量被限制在了monitor goroutine中, 只有这个goroutine可以读写金额。
		// 其他goroutine只能通过chan来操作金额, chan是线程安全的
		var b int
		for {
			select {
			case amount := <-deposits: //有存钱动作 计算金额
				b = b + amount
			case balances <- b: //将金额设置到通道中

			case <-clearChan:
				b = 0
			}
		}
	}()

	countNum := 0 //执行次数
	errNum := 0   //错误次数
	for {
		ClearSafe()
		countNum++
		done := make(chan int)

		//g1
		go func() {
			DepositSafe(100)
			done <- 0
		}()
		//g2
		go func() {
			DepositSafe(200)
			done <- 0
		}()

		<-done
		<-done

		if BalanceSafe() != 300 {
			errNum++
			fmt.Printf("balance = %v, count = %v; error = %v \n", Balance(), countNum, errNum)
		} /*else {
			fmt.Println("ok")
		}*/
	}
}

/**----------------------------------使用chan控制只有一个goroutine访问共享变量---------------------------------------------*/

// 2.
var (
	//sema = make(chan int, 1)
	//mu sync.Mutex
	mu sync.RWMutex
)

func DepositV2(amount int) {
	defer func() {
		//<-sema
		mu.Unlock()
	}()
	//sema <- 0
	mu.Lock()
	balance = balance + amount
}

func BalanceV2() int {
	defer func() {
		//<-sema
		mu.RUnlock()
	}()
	//sema <- 0
	mu.RLock()
	return balance
}

func ClearV2() {
	defer func() {
		//<-sema
		mu.Unlock()
	}()
	//sema <- 0
	mu.Lock()
	balance = 0
}

func TestSafeBankV2(t *testing.T) {
	for {
		ClearV2() //每次开始需要重新置为0
		//wait保证两个goroutine都执行完成再进行判断
		done := make(chan int)
		//g1
		go func() {
			DepositV2(100)
			done <- 1
		}()
		//g2
		go func() {
			DepositV2(200)
			done <- 1
		}()

		<-done
		<-done
		fmt.Println(BalanceV2() != 300)
	}
}

/**------------------------------------------------------------------------------------------------------------------*/

type Singleton struct {
	name string
}

var (
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

var once sync.Once

//提供了一种更快的方式
func getInstanceOnce() *Singleton {
	once.Do(func() {
		INSTANCE = &Singleton{name: "2333"}
	})
	return INSTANCE
}

func TestSingleton(t *testing.T) {
	for i := 0; i < 10000; i++ {
		var s1 *Singleton
		var s2 *Singleton
		done := make(chan int, 2)
		go func() {
			s1 = getInstanceLock()
			done <- 0
		}()
		go func() {
			s2 = getInstanceLock()
			done <- 0
		}()
		<-done
		<-done
		if s1 != s2 {
			fmt.Printf("s1 = %v, s2 = %v, s1 == s2: %v \n", &s1, &s2, s1 == s2)
		}
		INSTANCE = nil
	}
}

func TestSpeed(t *testing.T) {
	start := time.Now() // 获取当前时间
	for i := 0; i < 10000000; i++ {
		//getInstanceLock()//135.4997ms
		getInstanceOnce() //18.9349ms
	}
	elapsed := time.Since(start)
	fmt.Println("used time: ", elapsed)
}
