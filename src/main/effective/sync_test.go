package effective

import (
	"fmt"
	"testing"
)

var balance int

func Deposit(amount int) {
	balance = balance + amount
}

func TestBankUnSafe(t *testing.T) {

	countNum := 0 //执行次数
	errNum := 0   //错误次数
	for {
		countNum++
		//wait保证两个goroutine都执行完成再进行判断
		wait := make(chan int, 2)
		//g1
		go func() {
			Deposit(100)
			wait <- 1
		}()
		//g2
		go func() {
			Deposit(200)
			wait <- 1
		}()

		<-wait
		<-wait
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
			balance = 0
		}
		if balance == 300 {
			balance = 0
		}
	}
}
