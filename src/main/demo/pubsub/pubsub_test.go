package pubsub

import (
	"fmt"
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	ch := make(chan int)

	go func() {
		time.Sleep(10e9)
		ch <- 1
	}()
	select {
	case <-ch:
		fmt.Println("case 1")
	case <-time.After(5e9):
		fmt.Println("case 2")
	default:
		fmt.Println("default")
	}

}
