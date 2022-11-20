package effective

import (
	"fmt"
	"sync"
	"testing"
)

var balance int

func Deposit(amount int) {
	balance = balance + amount
}

func Balance() int {
	return balance
}

func clear() {
	balance = 0
}

func TestUnSafeBank(t *testing.T) {

	countNum := 0 //执行次数
	errNum := 0   //错误次数
	for {
		clear() //每次开始需要重新置为0
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

func Clear() {
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
		Clear()
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
		mu.Unlock()
		mu.RUnlock()
	}()
	//sema <- 0
	mu.RUnlock()
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
