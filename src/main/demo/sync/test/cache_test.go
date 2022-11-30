package test

import (
	"fmt"
	"testing"
	"time"
)

func TestClose(t *testing.T) {

	ch := make(chan struct{})

	go func() {
		time.Sleep(1e9)
		close(ch)
	}()

	<-ch
	<-ch
	<-ch
	<-ch

	fmt.Println("over")

}

func TestRecv(t *testing.T) {

	strings := make(chan string)

	close(strings)

	//从关闭的channel中接收 不会阻塞
	<-strings
	<-strings
	<-strings

	fmt.Println("over")

}
