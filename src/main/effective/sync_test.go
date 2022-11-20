package effective

import (
	"fmt"
	"testing"
)

var balance int

func Deposit(amount int) {
	balance = balance + amount
}

func Balance() int {
	return balance
}

func TestBankUnSafe(t *testing.T) {

	countNum := 0 //执行次数
	errNum := 0   //错误次数
	for {
		balance = 0 //每次开始需要重新置为0
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
		if balance != 300 {
			/**
			g1 read balance(0) + amount(100)
			g2 read balance(0) + amount(200)
			g2 write balance = 200 //丢失
			g1 write balance = 100

			最终balance不等于300
			*/

			errNum++
			fmt.Printf("balance = %v, count = %v; error = %v \n", balance, countNum, errNum)

		}
	}
}

//线程安全改造

var deposits = make(chan int) //发送存钱事件
var balances = make(chan int) //金额

func DepositSafe(amount int) {
	deposits <- amount
}

func BalanceSafe() int {
	return <-balances
}

func TestBankSafe(t *testing.T) {

	//先启动读监听
	go func() {
		var b int
		for {
			select {
			case amount := <-deposits:
				b = b + amount
			case balances <- b:

			}
		}
	}()

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

	fmt.Println(BalanceSafe() == 300)

}
