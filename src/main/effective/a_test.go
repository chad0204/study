package effective

import (
	"fmt"
	"testing"
)

func TestLock(t *testing.T) {

	c := make(chan string, 10)
	c <- "1"
	c <- "2"
	c <- "3"
	c <- "4"
	c <- "5"
	c <- "6"
	c <- "7"
	c <- "8"
	c <- "9"
	c <- "10"
	c <- "11"

	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println(<-c)

}
