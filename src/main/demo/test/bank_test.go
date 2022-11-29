package test

import (
	"fmt"
	"testing"
)

var (
	deposit = make(chan int)
	balance = make(chan int)
)

func Deposit(amount int) {
	deposit <- amount
}

func Balance() int {
	return <-balance
}

func TestV(t *testing.T) {
	go func() {
		var b = 0
		for {
			select {
			case amount := <-deposit:
				b = b + amount
			case balance <- b:
			}
		}
	}()

	done := make(chan int)

	go func() {
		Deposit(100)
		done <- 0
	}()

	go func() {
		Deposit(200)
		done <- 0
	}()

	<-done
	<-done

	fmt.Println(Balance())

}
