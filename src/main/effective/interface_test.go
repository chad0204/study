package effective

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sort"
	"testing"
	"time"
)

/**
永远不要使用一个指针指向一个接口类型，因为它已经是一个指针
*/

// 如果接口只有一个方法, 命名规范是加er、able后缀。如Reader
type Shaper interface {
	Area() int

	//other() int//结构体必须实现接口的所有方法才算实现该接口
}

type Square struct {
	side int
}

func (d *Square) Area() int {
	return d.side * d.side
}

type Triangle struct {
	high   int
	bottom int
}

func (t *Triangle) Area() int {
	return t.high * t.bottom / 2
}

func (t *Triangle) getHigh() int {
	return t.high
}

func TestInterface(t *testing.T) {
	var s Shaper   // 本质是一个指针, 没有赋值前不能直接使用
	s = &Square{5} // 将结构体类型的指针赋值给接口类型的变量
	fmt.Println(s.Area())

	//多态
	shapers := make([]Shaper, 2)
	shapers[0] = &Triangle{3, 4}
	shapers[1] = &Square{3}
	printArea(shapers[0])
	printArea(shapers[1])

}

func printArea(shaper Shaper) {
	//多态
	fmt.Println(shaper.Area())
}

type ReadWrite interface {
	Read(b bytes.Buffer) bool
	Write(b bytes.Buffer) (bool, error)
}

type Lock interface {
	Lock()
	UnLock()
}

// File接口 包含匿名接口 嵌套接口
type File interface {
	ReadWrite
	Lock
	Close()
}

type DefaultFile struct {
	buffer bytes.Buffer
}

func (d *DefaultFile) Write(b bytes.Buffer) (bool, error) {
	return false, fmt.Errorf("")
}

func (d *DefaultFile) Read(b bytes.Buffer) bool {
	return false
}

func (d *DefaultFile) Lock() {

}

func (d *DefaultFile) UnLock() {

}

func (d *DefaultFile) Close() {

}

// 嵌套接口
func TestNested(t *testing.T) {
	var file File
	file = new(DefaultFile)
	file.Close()
}

// 断言检查
func TestTypeAssert(t *testing.T) {

	//varI.(T) varI必须是个接口变量。T要不要用*取决于实现时的接收者类型

	//判断接口的具体实现
	shapers := []Shaper{&Square{5}, &Triangle{1, 4}}
	for _, s := range shapers {
		if square, ok := s.(*Square); ok {
			fmt.Println("this is a square", square)
		}

		if _, ok := s.(*Triangle); ok {
			fmt.Println("this is a triangle")
		}
	}

	//判断类型是否实现某个接口

	var square Shaper
	square = &Square{9}
	//triangle := &Triangle{2, 4}

	if s, ok := square.(Shaper); ok {
		fmt.Printf("this is a Shaper, area = %v \n", s.Area())
	}
}

// type-switch 通过type判断
func TestTypeSwitch(t *testing.T) {
	var f float32 = 3.3
	typeSwitch("stringParam", Square{9}, f)

}

func TestWriter(t *testing.T) {

}

func writeHeader(w io.Writer, contentType string) error {
	//需要copy []byte
	if _, err := w.Write([]byte(contentType)); err != nil {
		return err
	}
	//无需copy []byte
	if _, err := writeString(w, contentType); err != nil {
		return err
	}
	return nil
}

func writeString(w io.Writer, s string) (n int, err error) {
	//type stringWriter interface {
	//	WriteString(string) (n int, err error)
	//}
	//判断w是不是stringWriter类型, 是的话直接执行WriteString, 无需将字符串转未byte数组
	if sw, ok := w.(io.StringWriter); ok {
		return sw.WriteString(s) // avoid a copy
	}
	return w.Write([]byte(s)) // allocate temporary copy
}

func typeSwitch(args ...interface{}) {
	for i, v := range args {
		switch v.(type) {
		case Square:
			fmt.Printf("Param %d is a Square \n", v)
		case float32:
			fmt.Printf("Param %d is a float32 \n", v)
		case string:
			fmt.Printf("Param %d is a string \n", v)
		case nil:
			fmt.Printf("Param %d is a nil \n", v)
		default:
			fmt.Printf("Param %d is unknown\n", i)
		}
	}

}

type Lener interface {
	Len() int
}

type Appender interface {
	Append(int)
}

// 类型List实现了Lener和Appender接口
type List []int

func (l List) Len() int {
	return len(l)
}

func (l *List) Append(val int) {
	*l = append(*l, val) //如果接收者用值, 赋值操作时IDE会提示Assignment to the method receiver doesn't propagate to other calls
}

// 在接口上调用方法的方法 （接口本质是一个指针）
func CountInto(a Appender, start, end int) {
	for i := start; i <= end; i++ {
		a.Append(i)
	}
}
func LongEnough(l Lener) bool {
	fmt.Println(l)
	return l.Len()*10 > 42
}

// 值和类型
func TestPointAndValue(t *testing.T) {
	//众所周知 方法接收者不管是值还是指针, 指针和值方法都可以调用（编译器隐式转换）
	var l List

	//List实现了Appender, 却无法作为CountInto的Appender类型参数, 因为实现List.Append的接收者是List指针, 和CountInto形参类型是值不同
	//CountInto(l, 1, 10) // 可以使用&l

	//List实现了Lener, 可以作LongEnough的Lener类型参数, 因为实现List.Len的接收者是List值, 和LongEnough形参一致
	if LongEnough(l) {
		fmt.Printf("l is long enough\n")
	} else {
		fmt.Printf("l is not long enough\n")
	}

	pointL := new(List)

	//可以作为pointL可以作为Appender, 因为实现Append方法的接收者是指针
	CountInto(pointL, 1, 10)

	//奇怪的是这里也可以, LongEnough参数类型不是指针, 而实现Len方法的接收者是值。因为类型指针会被自动解引用（当形参类型是接口时）
	if LongEnough(pointL) {
		fmt.Printf("l is long enough\n")
	} else {
		fmt.Printf("l is not long enough\n")
	}

	//在接口上调用方法时，必须有和方法定义时相同的接收者类型或者是可以从具体类型 P 直接可以辨识的：
	//1 指针方法的对象, 作为接口类型参数时, 只能将指针赋值给接口参数 CountInto(&l, 1, 10)可以, CountInto(l, 1, 10)不行, 因为存储在接口中的值没有地址。
	//2 值方法的对象, 作为接口类型参数时, 指针或值都可以赋值给接口参数  LongEnough(l), LongEnough(&l)也可以, 因为指针会被解引用。

	//CountInto(l, 1, 10) 将一个值赋值给一个接口时，编译器会确保所有可能的接口方法都可以在此值上被调用，因此不正确的赋值在编译期就会失败。

}

type day struct {
	num       int
	shortName string
	longName  string
	workNum   int //工作日排序, 周一为第一天
}

// v1
type daySlice struct {
	data []day //引用快递
}

func (x daySlice) Len() int           { return len(x.data) }
func (x daySlice) Less(i, j int) bool { return x.data[i].num < x.data[j].num }
func (x daySlice) Swap(i, j int)      { x.data[i], x.data[j] = x.data[j], x.data[i] }

// v2 常用
type ByNum []*day //day可以用指针也可以用值, 指针更省空间

func (x ByNum) Len() int           { return len(x) }
func (x ByNum) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x ByNum) Less(i, j int) bool { return x[i].num < x[j].num }

type ByWorkNum []*day // 和ByNum有重复方法, 可以复用Len和Swap.  type days []*day, type ByNum struct {days}, ype ByWorkNum struct {days}

func (x ByWorkNum) Len() int           { return len(x) }
func (x ByWorkNum) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x ByWorkNum) Less(i, j int) bool { return x[i].workNum < x[j].workNum }

//通过内置函数, 实现时自定义排序
type customSort struct {
	days []*day
	less func(x, y *day) bool
}

func (c customSort) Len() int           { return len(c.days) }
func (c customSort) Swap(i, j int)      { c.days[i], c.days[j] = c.days[j], c.days[i] }
func (c customSort) Less(i, j int) bool { return c.less(c.days[i], c.days[j]) }

// 自定义weekday排序
func TestSored(t *testing.T) {

	Sunday := day{0, "SUN", "Sunday", 6}
	Monday := day{1, "MON", "Monday", 0}
	Tuesday := day{2, "TUE", "Tuesday", 1}
	Wednesday := day{3, "WED", "Wednesday", 2}
	Thursday := day{4, "THU", "Thursday", 3}
	Friday := day{5, "FRI", "Friday", 4}
	Saturday := day{6, "SAT", "Saturday", 5}

	// v1
	data := []day{Tuesday, Thursday, Wednesday, Sunday, Monday, Friday, Saturday}
	a := new(daySlice)
	a.data = data

	sort.Sort(a)
	if !sort.IsSorted(a) {
		fmt.Errorf("fail %v", a)
	}
	fmt.Printf("v1 sorted: ")
	for _, d := range data {
		fmt.Printf("%s ", d.longName)
	}
	fmt.Printf("\n")

	// v2
	dataV2 := []*day{&Tuesday, &Thursday, &Wednesday, &Sunday, &Monday, &Friday, &Saturday}
	v2 := ByNum(dataV2) // 语法糖
	sort.Sort(ByNum(dataV2))
	if !sort.IsSorted(v2) {
		fmt.Errorf("fail %v", v2)
	}
	fmt.Printf("v2 sorted: ")
	for _, d := range v2 {
		fmt.Printf("%s ", d.longName)
	}
	fmt.Printf("\n")

	//逆序
	v3 := ByWorkNum(dataV2)
	sort.Sort(sort.Reverse(v3)) //Reverse包一层, 反转了less
	if !sort.IsSorted(v3) {
		fmt.Errorf("fail %v", v3)
	}
	fmt.Printf("v3 sorted: ")
	for _, d := range v3 {
		fmt.Printf("%s ", d.longName)
	}
	fmt.Printf("\n")

	//customSort
	sort.Sort(customSort{dataV2, func(x, y *day) bool {
		if x.num != y.num {
			return x.num < y.num
		}
		if x.workNum != y.workNum {
			return x.workNum < y.workNum
		}
		return false
	}})
}

// interface {} 变量在内存中占据两个字长：一个用来存储它包含的类型，另一个用来存储它包含的数据或者指向数据的指针
type Any interface{}

// 自定义类型
type SpecialString string

// 空接口
func TestEmptyInterface(t *testing.T) {
	var val Any
	val = 1
	fmt.Printf("any is %v \n", val)
	val = "ss"
	fmt.Printf("any is %v \n", val)
	val = Square{1}
	fmt.Printf("any is %v \n", val)

	switch t := val.(type) {
	case int:
		fmt.Printf("any is %v, type is %T \n", t, t)
	case Square:
		fmt.Printf("any is %v, type is %T \n", t, t)
	case string:
		fmt.Printf("any is %v, type is %T \n", t, t)
	}

	testFunc := func(any Any) {
		switch any.(type) {
		case int:
			fmt.Printf("any %v is a int type \n", any)
		case string:
			fmt.Printf("any %v is a string type \n", any)
		case SpecialString:
			fmt.Printf("any %v is a SpecialString type \n", any)
		}
	}
	var str SpecialString = "hello"
	testFunc(str)
	var stri string = "world"
	testFunc(stri)

}

type Element interface{}

type Vector struct {
	element []Element
}

func (v *Vector) At(i int) Element {
	return v.element[i]
}

func (v *Vector) set(i int, element Element) {
	v.element[i] = element
}

// 通用类型的list
func TestElement(t *testing.T) {

	var list Vector
	ele := make([]Element, 100)
	list.element = ele
	list.set(1, 10)
	list.set(2, "AAA")
	fmt.Println(list.At(2))

}

// 复制自定义类型slice到空接口类型slice
func TestCopySliceToEmptyInterface(t *testing.T) {

	var s1 []Square = []Square{{1}, {2}, {3}}

	var s2 []interface{} = make([]interface{}, len(s1))

	//Cannot use 's1' (type []Square) as the type []interface{}
	//s2 = s1

	//只能遍历复制
	for i, v := range s1 {
		s2[i] = v
	}
}

type ReaderWriter struct {
	io.Reader
	io.Writer
}

// 通过内嵌结构体实现继承, 也可以实现接口继承
func TestMoreImpl(t *testing.T) {
	var rw ReaderWriter
	rw.Reader.Read([]byte{})
}

type ILock interface {
	TryLock() bool
	Release()
}

type ReentrantLock struct {
}

func (r ReentrantLock) TryLock() bool {
	fmt.Println("tryLock")
	return true
}

func (r ReentrantLock) Release() {
	fmt.Println("Release")
}

func checkAndSet(lock ILock) {
	if l := lock.TryLock(); l {
		lock.Release()
	}
}

// receiver
func TestReceiver(t *testing.T) {
	lock := ReentrantLock{}

	//如果结构体实现接口的方法都是值方法, 那么指针和值都可以作为接口参数。只要有一个方法不是值方法, 就只能用指针作为接口参数
	// 修改TryLock或Release的receiver为指针可以看到编译报错
	checkAndSet(lock)
	checkAndSet(&lock)
}

func testNil(lock ILock) {
	//lock包含的是值为nil的ReentrantLock指针, 但是lock不是nil接口。这里并不能起到保护作用
	//T = *effective.ReentrantLock,v = <nil>
	//T = <nil>,v = <nil>
	fmt.Printf("T = %T,v = %v \n", lock, lock)

	if lock != nil {
		lock.TryLock() //nil pointer
	}
}

const debug = false

// 包含nil指针的接口 不是nil接口
func TestNilInterface(t *testing.T) {
	//var lock *ReentrantLock //T = *effective.ReentrantLock,v = <nil>
	var lock ILock //T = <nil>,v = <nil>
	if debug {
		lock = &ReentrantLock{}
	}
	testNil(lock)
}

var period = flag.Duration("period", 5*time.Second, "sleep period")

// flag.Value: 给命令行标记定义新的符号
// 如果是个main方法, build后, ./sleep -period 50ms
func TestFlagValue(t *testing.T) {
	flag.Parse()
	fmt.Printf("Sleeping for %v... \n", *period)
	//time.Sleep(*period)
}

type ServerHandler struct {
}

//如果使用指针receiver, 接口做为参数时, 不能使用接口保存对象值。接口中的值没有地址, 只能用存指针。如果不是需要改变原变量, 直接用值做receiver
func (s *ServerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
}

func TestHttpHandler(t *testing.T) {
	var hdl http.Handler
	hdl = &ServerHandler{}
	http.ListenAndServe("localhost:8080", hdl)
}

//实现error接口
type MyError struct {
	msg string
}

func (e *MyError) Error() string {
	return e.msg
}

//工厂方法 一般用fmt.Errorf("EOF")
func New(msg string) error {
	return &MyError{msg: msg}
}

//error接口
func TestErrorInterface(t *testing.T) {
	fmt.Println(New("EOF") == New("EOF"))
	fmt.Println(fmt.Errorf("EOF") == fmt.Errorf("EOF"))
}
