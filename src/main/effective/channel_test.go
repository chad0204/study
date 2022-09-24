package effective

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {

	//声名一个名为c1的channel, 传输string类型的数据（一个通道只能传输一种类型的数据）
	var c1 chan string
	c1 = make(chan string)
	//c2 := make(chan int) //声明并初始化内存

	//默认是无缓存通道。发送一个,没有接收则下次发送阻塞。接收完成, 没有发送则下次接收阻塞
	go sendData(c1)
	go receiveData(c1)

	time.Sleep(10e9)

}

func sendData(c chan string) {
	c <- "a"
	c <- "c"
	c <- "b"
	fmt.Println("send over")
}

func receiveData(c chan string) {
	for {
		input := <-c
		time.Sleep(2e9)
		fmt.Println(input)
	}
}

//无缓冲通道, 发送一个数据后当前协程就会block.
func TestDeadLock(t *testing.T) {
	c := make(chan string)

	//go func() {
	//	c <- "data"
	//}()
	//go func() {
	//	fmt.Println(<- c)
	//}()

	c <- "data"      // block当前协程
	fmt.Println(<-c) //永远无法执行。除非上面一行在另一个协程中执行, 或者c是个有缓冲通道

	time.Sleep(1e9)

}

//缓冲通道具备弹性, 只有空和满会阻塞
func TestBufferChannel(t *testing.T) {

	ch := make(chan string, 100)

	// i < 100, 否则deadlock
	for i := 0; i < 10; i++ {
		ch <- strconv.Itoa(i)
	}

	for i := 0; i < 10; i++ {
		fmt.Println(<-ch)
	}

}

//想让主协程在子协程完成后退出
func TestBlock(t *testing.T) {

	ch := make(chan int)
	go func() {
		fmt.Println("starting...")
		time.Sleep(2e9)
		fmt.Println("ending...")
		ch <- 0
	}()

	<-ch
}

// 主协程阻塞到多个协程执行完成 （no buffer）
func TestSemaphore(t *testing.T) {
	ch := make(chan int)
	arr := make([]int, 100)
	pivot := 5

	go func(nums []int) {
		time.Sleep(2e9)
		ch <- 0
		fmt.Println("sort1 done...")
	}(arr[:pivot])

	go func(nums []int) {
		time.Sleep(5e9)
		ch <- 0
		fmt.Println("sort2 done...")
	}(arr[pivot:])

	//阻塞到两个协程结束
	<-ch
	<-ch
}

// 信号量
func TestSemaphoreV2(t *testing.T) {

	N := 10

	ch := make(chan int, N) //10个

	for i := 0; i < N; i++ {
		go func() {
			time.Sleep(2e9)
			ch <- 0
			fmt.Println("process done...")
		}()
	}

	//阻塞到所有协程结束
	for i := 0; i < N; i++ {
		<-ch
	}
}
