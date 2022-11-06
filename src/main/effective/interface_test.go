package effective

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"testing"
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
	s = &Square{5} // 将结构体类型的指针赋值给接口类型的变了
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
		fmt.Printf("this is a square, area = %v \n", s.Area())
	}
}

// type-switch 通过type判断
func TestTypeSwitch(t *testing.T) {
	var f float32 = 3.3
	typeSwitch("stringParam", Square{9}, f)

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
	*l = append(*l, val) //如果接收者用值, 赋值操作会提示Assignment to the method receiver doesn't propagate to other calls
}

// 在接口上调用方法的方法 （接口本质是一个指针）
func CountInto(a Appender, start, end int) {
	for i := start; i <= end; i++ {
		a.Append(i)
	}
}
func LongEnough(l Lener) bool {
	return l.Len()*10 > 42
}

// 值和类型
func TestPointAndValue(t *testing.T) {
	//众所周知 方法接收者不管是值还是指针, 指针和值变量都可以调用

	var l List

	//list实现了Appender, 却无法作用Appender参数, 因为实现Appender.Append的接收者是list指针, 和CountInto形参不同
	//CountInto(l, 1, 10) // 可以使用&l

	//list实现了Lener, 可以作Lener参数, 因为实现Lener.Len的接收者是list值, 和LongEnough形参一致
	if LongEnough(l) {
		fmt.Printf("l is long enough\n")
	} else {
		fmt.Printf("l is not long enough\n")
	}

	pl := new(List)

	//可以作为pl可以作为Appender, 因为实现Append方法的接收者是指针
	CountInto(pl, 1, 10)

	//奇怪的是这里也可以, LongEnough参数类型不是指针, 而实现Len方法的接收者是值。因为类型指针会被自动解引用（形参是接口）
	if LongEnough(pl) {
		fmt.Printf("l is long enough\n")
	} else {
		fmt.Printf("l is not long enough\n")
	}

	//在接口上调用方法时，必须有和方法定义时相同的接收者类型或者是可以从具体类型 P 直接可以辨识的：
	//1 指针方法可以通过指针调用 CountInto(pl, 1, 10)
	//2 值方法可以通过值调用 LongEnough(l)
	//3 接收者是值的方法可以通过指针调用，因为指针会首先被解引用 LongEnough(pl)
	//4 接收者是指针的方法不可以通过值调用，因为存储在接口中的值没有地址 CountInto(l, 1, 10)

	//CountInto(l, 1, 10) 将一个值赋值给一个接口时，编译器会确保所有可能的接口方法都可以在此值上被调用，因此不正确的赋值在编译期就会失败。

	data := []int{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586}

	sort.Ints(data)

}

type day struct {
	num       int
	shortName string
	longName  string
}

// v1
type daySlice struct {
	data []*day
}

func (x *daySlice) Len() int           { return len(x.data) }
func (x *daySlice) Less(i, j int) bool { return x.data[i].num < x.data[j].num }
func (x *daySlice) Swap(i, j int)      { x.data[i], x.data[j] = x.data[j], x.data[i] }

// v2
type daySliceV2 []day //可以用指针 也可以用值

func (x daySliceV2) Len() int           { return len(x) }
func (x daySliceV2) Less(i, j int) bool { return x[i].num < x[j].num }
func (x daySliceV2) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// 自定义weekday排序
func TestSored(t *testing.T) {

	Sunday := day{0, "SUN", "Sunday"}
	Monday := day{1, "MON", "Monday"}
	Tuesday := day{2, "TUE", "Tuesday"}
	Wednesday := day{3, "WED", "Wednesday"}
	Thursday := day{4, "THU", "Thursday"}
	Friday := day{5, "FRI", "Friday"}
	Saturday := day{6, "SAT", "Saturday"}

	// v1
	data := []*day{&Tuesday, &Thursday, &Wednesday, &Sunday, &Monday, &Friday, &Saturday}
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
	dataV2 := []day{Tuesday, Thursday, Wednesday, Sunday, Monday, Friday, Saturday}
	v2 := daySliceV2(dataV2) // 语法糖
	sort.Sort(daySliceV2(dataV2))
	if !sort.IsSorted(v2) {
		fmt.Errorf("fail %v", v2)
	}
	fmt.Printf("v2 sorted: ")
	for _, d := range v2 {
		fmt.Printf("%s ", d.longName)
	}
	fmt.Printf("\n")

	//引用传递, 如果用指针 应该是防止copy
	changeSlice(v2)
	for _, d := range v2 {
		fmt.Printf("%s ", d.longName)
	}
	fmt.Printf("\n")
}

func changeSlice(days daySliceV2) {
	days[0] = day{7, "unKnow", "unKnow"}
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

//复制自定义类型slice到空接口类型slice
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
