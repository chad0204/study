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

	//默认是无缓存通道。发送完成,没有被接收之前当前协程阻塞。 接收协程也是阻塞,直到有其他协程发送
	//接收和发送都会阻塞协程。
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

//无缓冲通道, 发送一个数据后当前协程就会block. 接收和发送需要异步处理
func TestDeadLock(t *testing.T) {
	c := make(chan string)

	//go func() {
	//	c <- "data"
	//	//block until receive...
	//}()
	//go func() {
	//	fmt.Println(<- c)
	//  //block until send...
	//}()

	//一个协程是做不了同时发送和接收的
	c <- "data"      // block当前协程
	fmt.Println(<-c) //永远无法执行。除非上面一行在另一个协程中执行, 或者c是个有缓冲通道

	time.Sleep(1e9)

}

//缓冲通道具备弹性, 只有空和满会阻塞
func TestBufferChannel(t *testing.T) {

	ch := make(chan string, 100)

	// 但协程下 i < 100, 否则deadlock
	for i := 0; i < 101; i++ {
		ch <- strconv.Itoa(i)
	}

	for i := 0; i < 101; i++ {
		fmt.Println(<-ch)
	}

}

// 想让主协程在子协程完成后退出
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
	/*
		三个协程
		主协程开启两个协程, 两个子协程执行结束后向通道中塞值, 主协程取值
	*/

	//int channel
	chSemaphore := make(chan int)
	//模拟快排
	arr := make([]int, 100)
	pivot := 5

	go func(nums []int) {
		time.Sleep(2e9)
		chSemaphore <- 0
		fmt.Println("sort1 done...")
	}(arr[:pivot])

	go func(nums []int) {
		time.Sleep(5e9)
		chSemaphore <- 0
		fmt.Println("sort2 done...")
	}(arr[pivot:])

	//阻塞到两个协程结束
	<-chSemaphore
	<-chSemaphore
}

// 有缓存channel实现信号量
func TestSemaphoreV2(t *testing.T) {

	N := 10
	chSemaphore := make(chan int, N) //10个

	for i := 0; i < N; i++ {
		go func() {
			time.Sleep(2e9)
			chSemaphore <- 0
			fmt.Println("process done...")
		}()
	}

	//阻塞到所有协程结束
	for i := 0; i < N; i++ {
		<-chSemaphore
	}
}

// foreach channel 从channel中读取数据, 直到通道关闭（会自动检测）才会向下执行
func TestChanFor(t *testing.T) {

	ch := make(chan int, 3)

	//go func() {
	//	//从channel中读取数据, 直到通道关闭 才往下执行
	//	for v := range ch {
	//		fmt.Printf("%v \n", v)
	//	}
	//	fmt.Println("exec after closed")//关闭channel后才会执行
	//}()

	go func() {
		ch <- 0
		ch <- 1
		ch <- 2
		time.Sleep(2e9)
		close(ch)                      //主动关闭
		fmt.Println("channel closing") //关闭channel
	}()

	for v := range ch {
		fmt.Printf("%v \n", v)
	}

	//等价于
	//for {
	//	v, open := <-ch
	//	if !open {
	//		break
	//	}
	//	fmt.Printf("%v \n", v)
	//}

	fmt.Println("exec after closed") //关闭channel后才会执行

	time.Sleep(10e9)

}

//返回一个只读的channel
func pump() <-chan int {
	ch := make(chan int)
	//协程中执行, 不然阻塞
	go func() {
		for i := 0; ; i++ {
			ch <- i
		}
	}()
	return ch
}

//取
func suck(ch <-chan int) {
	go func() {
		for {
			fmt.Println(<-ch)
		}
	}()
}

//取
func suckForeach(ch <-chan int) {
	go func() {
		for v := range ch {
			fmt.Println(v)
		}
	}()
}

//管道工厂
func TestChannelFactory(t *testing.T) {
	suckForeach(pump())
	time.Sleep(1e9)
}

//通道方向
func TestSendRecvOnly(t *testing.T) {
	//只能向通道发送数据
	sendOnly := make(chan<- string)
	//只接收的通道（<-chan T）无法关闭, 准备的说是不必关闭, 关闭表示不能向通道发送数据。
	recvOnly := make(<-chan int)

	go func(send chan<- string, recv <-chan int) {
		for i := range recv {
			result := strconv.Itoa(i) + ">>>"
			send <- result
		}

	}(sendOnly, recvOnly)

}

func generate() chan int {
	ch := make(chan int)
	go func() {
		for i := 2; ; i++ {
			ch <- i
		}
	}()
	return ch
}

//从in channel过滤出不能被prime整除的数 输出到outCh
func filter(in chan int, prime int) chan int {
	outCh := make(chan int)
	go func() {
		for {
			//3 4 5 6 7 8 9 10 11 12 13..  %2
			//5 7 9 11 13.. %3
			//7 11 13.. %5
			//..
			if i := <-in; i%prime != 0 {
				outCh <- i
			}
		}
	}()
	return outCh
}

func sieve() chan int {
	out := make(chan int)
	go func() {
		numberCh := generate()
		for {
			prime := <-numberCh
			numberCh = filter(numberCh, prime)
			out <- prime //每次输出一个素数
		}
	}()
	return out
}

//输出素数 prime number
func TestPrimeNumber(t *testing.T) {
	primeCh := sieve()

	for prime := range primeCh {
		fmt.Println(prime)
	}
}

//close

//关闭通道 只有发送者需要关闭通道 表示告诉接收者不会再有新的值了
func generateV2() chan int {
	ch := make(chan int)
	go func() {
		for i := 2; i < 100; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

//从in channel过滤出不能被prime整除的数 输出到outCh
func filterV2(in chan int, prime int) chan int {
	outCh := make(chan int)
	go func() {
		for {
			i, ok := <-in
			if !ok {
				close(outCh)
				break
			}
			if i%prime != 0 {
				outCh <- i
			}
		}

	}()
	return outCh
}

func sieveV2() chan int {
	out := make(chan int)
	go func() {
		numberCh := generateV2()
		for {
			prime, ok := <-numberCh
			if !ok {
				close(out)
				break
			}
			numberCh = filterV2(numberCh, prime)
			out <- prime //每次输出一个素数
		}
	}()
	return out
}

//输出素数 prime number 带关闭
func TestPrimeNumberWithClose(t *testing.T) {
	primeCh := sieveV2()

	//会自动校验channel是否关闭
	for prime := range primeCh {
		fmt.Println(prime)
	}
}

func pump1(ch chan<- int) {
	for i := 0; ; i++ {
		time.Sleep(10)
		if i%2 == 0 {
			ch <- i
		}
	}
}

func pump2(ch chan<- int) {
	for i := 0; ; i++ {
		time.Sleep(5)
		if i%2 != 0 {
			ch <- i
		}
	}
}

func selectSuck(c1, c2 chan int) {
	for {
		// 都阻塞 select阻塞. 如果有default不阻塞, 执行default
		// 都有值 随机执行
		//
		select {
		case v := <-c1:
			fmt.Printf("receiver from channel1, value = %v \n", v)
		case v := <-c2:
			fmt.Printf("receiver from channel2, value = %v \n", v)
		default:
			fmt.Println("暂时无数据...")
		}
	}
}

func TestSelect(t *testing.T) {
	c1 := make(chan int)
	c2 := make(chan int)
	go pump1(c1)
	go pump2(c2)
	go selectSuck(c1, c2)

	time.Sleep(1e9)
}

//使用tick进行限速
func TestTick(t *testing.T) {
	// 每秒执行10次
	rate_per_sec := 10
	var dur = time.Duration(1e9 / rate_per_sec)
	//返回的tick是一个只接收通道（保证外部只能进行读取操作）, 每dur会写入一个值
	tick := time.Tick(dur)

	for {
		<-tick
		fmt.Println("exec request")
	}

}

//Timer只设置一次时间
func TestTimer(t *testing.T) {

	tick := time.Tick(1e9)  //1s一次
	over := time.After(5e9) //5s后执行一次

	for {
		select {
		case <-tick:
			fmt.Println("tick")
		case <-over:
			fmt.Println("over")
			return
		default:
			//没有准备好的通道时, 每秒执行2次
			fmt.Println("    .")
			time.Sleep(5e8)
		}
	}
}

type Conn struct {
	replica int
}

func (c *Conn) doQueryDB() string {
	return "query result"
}

func TestQuery(t *testing.T) {

	conns := []Conn{{1}, {2}, {3}}

	//无缓存channel必须异步执行, 也就是存和取‘同时’发生。有缓存channel可以一个协程执行
	res := make(chan string, 1)

	//从多个副本数据库中查询数据, 第一个返回的就是结果
	for _, conn := range conns {
		//创建多个协程同时执行
		go func(c Conn) {
			select {
			case res <- c.doQueryDB():
			default:

			}
		}(conn)
	}
	fmt.Println(<-res)
}

func server(workChan chan string) {
	for work := range workChan {
		go safeDo(work)
	}

}

func safeDo(work string) {
	//defer在return之后函数返回之前执行
	//recover仅在defer中有效
	//使用recover()可以捕获异常 使其他协程继续执行
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Work failed with %v in %v \n", err, work)
		}
	}()
	doWork(work)
}

func doWork(work string) {
	if work == "error" {
		panic("fail") // 模拟协程错误
	}
	fmt.Println(work + " done.")
}

//停掉了服务器内部一个失败的协程而不影响其他协程的工作
func TestRecover(t *testing.T) {

	//有 panic 没 recover，程序宕机
	//有 panic 也有 defer recover，程序不会宕机，执行完对应的 defer 后，从宕机点退出当前函数后继续执行

	ch := make(chan string)

	go func() {
		ch <- "work1"
		ch <- "work2"
		ch <- "work3"
		ch <- "error"
		ch <- "work4"
		ch <- "work5"
	}()

	go server(ch)

	time.Sleep(10e9)

}
