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
