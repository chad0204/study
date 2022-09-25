package effective

import (
	"fmt"
	"math/rand"
	"testing"
)

func Test(t *testing.T) {

	var arr []int
	for i := 0; i < 1000; i++ {
		arr = append(arr, rand.Intn(10000))
	}

	chanReceive := make(chan int)

	go quickSort(0, len(arr)-1, arr, chanReceive)

	<-chanReceive

	fmt.Printf("arr len = %d, cap = %d, address = %p \n", len(arr), cap(arr), &arr)
	for _, v := range arr {
		fmt.Printf("%d,", v)
	}
}

func quickSort(start, end int, nums []int, chanSend chan int) {
	if start >= end {
		//输入结束标志
		chanSend <- 0
		return
	}
	lx := start
	rx := end
	p := nums[lx]
	for lx < rx {
		for nums[rx] >= p && lx < rx {
			rx--
		}
		for nums[lx] <= p && lx < rx {
			lx++
		}
		tmp := nums[lx]
		nums[lx] = nums[rx]
		nums[rx] = tmp
	}

	nums[start] = nums[lx]
	nums[lx] = p

	chanReceive := make(chan int)

	go quickSort(start, lx-1, nums, chanReceive)
	go quickSort(lx+1, end, nums, chanReceive)

	<-chanReceive
	<-chanReceive

	//输入结束标志
	chanSend <- 0
}
