package effective

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
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

	//本次排序结束, 输入结束标志, 让上一层结束
	chanSend <- 0
}

var values = []string{"a", "b", "c", "d", "e"}

// 循环里面闭包使用goroutine, 一定要copy到新的变量
func TestGoroutine(t *testing.T) {

	//version1
	for index := range values {
		func() {
			fmt.Print(index, " ")
		}()
	}
	fmt.Println()

	//version2 error
	//for index := range values {
	//	// index是单一变量,在所有goroutine中共享, 因为for很快执行完, 当协程执行到print时, index可能已经是最后一个索引值了
	//	go func() {
	//		fmt.Print(index, " ")
	//	}()
	//	//time.Sleep(1e8)
	//}
	//fmt.Println()

	//version2_1 操作value也是error
	//for _, value := range values {
	//	go func() {
	//		fmt.Print(value, " ")
	//	}()
	//}

	//version3 right
	for index := range values {
		go func(idx int) {
			fmt.Print(idx, " ")
		}(index)
	}
	time.Sleep(1e8)
	fmt.Println()

	//version3_1 right one
	for index := range values {
		val := values[index]
		go func() {
			fmt.Print(val, " ")
		}()
	}
	time.Sleep(1e8)
	fmt.Println()

	time.Sleep(1e9)
}
