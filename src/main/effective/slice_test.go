package effective

import (
	"fmt"
	"testing"
)

// 1. 数组是值传递、切片是引用传递。 数组传递需要拷贝值, 切片只需要传递引用, 效率更高。
// 2. 数组需要定义长度、切片长度可以动态变化。数组的长度是类型的一部分, [3]int和[2]int不是一个数组类型。
// 3. 多个slice表示同一数组, 这些slice可以共享存储。数组是切片的构建块。当你有个数组arr需要函数传递，最好创建一个切片arr[:],传递这个切片。
// go/src/runtime/slice.go

// 声明、零值、类型、值传递
func TestArray(t *testing.T) {

	//声明, 需要定义长度, 有初始化默认值
	var arr [5]int
	arr = [5]int{1, 2, 3, 4, 5}
	fmt.Printf("Array is %v\n", arr)

	// 初始化 ...表示长度就是元素个数 arr1、2、3等价
	//var arr1 = [5]int{5, 6, 7, 8, 22}
	var arr2 = [...]int{5, 6, 7, 8, 22}
	//var arr3 = [...]int{0: 5, 1: 6, 2: 7, 3: 8, 4: 22}

	//var arr4 = [5]string{0: "赋值索引0的字符串", 4: "赋值索引4的字符串"}
	//var arr5 = [...]int{99: 2333} //声明一个长度为100的数组, 最后一个索引99的元素初始化为2333

	// 没有定义长度 这是一个切片
	//var slice1 = []int{5, 6, 7, 8, 22}

	//是值类型
	fmt.Printf("arr address:%p \n", &arr2)
	exchange(arr2)
	//exchangeRef(&arr1)

	fmt.Printf("result: %v \n", arr2)

}

func exchange(arr [5]int) {
	//为什么这样不行？ 这样v也是copy了值, 修改不会影响原数组
	//for _,v := range arr {
	//	v = v + 1
	//}
	fmt.Printf("arr address:%p \n", &arr)
	for i, v := range arr {
		arr[i] = v + 1
	}
}

func exchangeRef(arr *[5]int) {
	for i, v := range arr {
		arr[i] = v + 1
	}
}

/*------------------------------------------------slice---------------------------------------------------------*/

/*
*
声明
零值 nil, new返回的是len=cap=0的切片
*/
func TestSlice(t *testing.T) {

	//声明 一个切片在未初始化之前默认为 nil，长度为 0
	var slice []string

	//定义一个原始数组
	arr := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	//初始化
	var slice1 []int = arr[0:9]
	//var slice2 []int = arr[:] //完整
	//var slice3 []int = arr[2:] //2 ~ len
	//var slice4 []int = arr[:4] //0 ~ 4

	//修改原始数组衍生的切片, 会影响原始数组的值
	slice1[2] = 2333
	fmt.Printf("arr %T, %v \n", arr, arr)

	//声明 + 初始化
	slice2 := []int{1, 2, 3}

	//如果原始数组不存在, 无法通过数组获取切片, 可以使用make
	s := make([]int, 10) //cap(s) == len(s) == 10
	fmt.Printf("s %d, %d, %v \n", len(s), cap(s), s)

	fmt.Printf("array %T, %v \n", [...]int{1, 2, 3}, [...]int{1, 2, 3})
	fmt.Printf("slice %T, %v \n", slice, slice)
	fmt.Printf("slice1 %T, %v \n", slice1, slice1)
	fmt.Printf("slice2 %T, %v \n", slice2, slice2)

	fmt.Printf("slice address:%p \n", slice2)
	exchangeSlice(slice2)
	fmt.Printf("exchange slice2: %T, %v \n", slice2, slice2)

}

// 引用传递
func exchangeSlice(arr []int) {
	fmt.Printf("slice address:%p \n", arr)
	for i, v := range arr {
		arr[i] = v + 1
	}
}

// 0 <= len(s) <= cap(s)
func TestLenAndCap(t *testing.T) {

	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// len是元素个数4, cap是最大容量8, 因为是从2开始切割也就是2到10 = 8
	slice := arr[2:6] // 左闭右开

	fmt.Printf("slice len = %d, cap = %d, %v \n", len(slice), cap(slice), slice)

	//这里[s:e], e不能超过8. TODO 为什么设置索引cap不能自动扩容？ 需要通过append(s1, 1)来扩容
	slice = slice[0:7]

	fmt.Printf("slice len = %d, cap = %d, %v \n", len(slice), cap(slice), slice)

	s := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k'}

	//	修改s1/2/3/4, 会作用到s TODO  append之后, 还是这样指向同一个数组嘛？
	s1 := s[1:4]
	s2 := s[2:]
	s3 := s[0:2]
	s4 := s[0:5] //end不能超过cap

	s4[0] = '$'

	// s = arr[start:end] cap = len(arr) - start, len = end - start
	fmt.Printf("s len = %d, cap = %d, address = %p, %c \n", len(s), cap(s), &s, s)
	fmt.Printf("s1 len = %d, cap = %d, address = %p, %c \n", len(s1), cap(s1), &s1, s1)
	fmt.Printf("s2 len = %d, cap = %d, address = %p, %c \n", len(s2), cap(s2), &s2, s2)
	fmt.Printf("s3 len = %d, cap = %d, address = %p, %c \n", len(s3), cap(s3), &s3, s3)
	fmt.Printf("s4 len = %d, cap = %d, address = %p, %c \n", len(s4), cap(s4), &s4, s4)

}

// make(T) 返回一个类型为 T 的初始值，它只适用于3种内建的引用类型：切片、map 和 channel
// new(T) 为每个新的类型T分配一片内存，初始化为 0 并且返回类型为*T的内存地址：这种方法 返回一个指向类型为 T，值为 0 的地址的指针，它适用于值类型如数组和结构体,它相当于 &T{}
func TestMake(t *testing.T) {
	//当没有前提数组时, 使用make创建切片
	s := make([]int, 0)

	fmt.Printf("s len = %d, cap = %d, %v \n", len(s), cap(s), s)

	s1 := make([]int, 10, 20)
	fmt.Printf("s1 len = %d, cap = %d, %v \n", len(s1), cap(s1), s1)

	//new返回的是值对象的地址(引用), make返回的是引用对象（map\slice\channel）的值
	s2 := new([]int)

	fmt.Printf("s2 len = %d, cap = %d, %v \n", len(*s2), cap(*s2), *s2)

}

func TestForRange(t *testing.T) {

	s := make([]int, 10)
	for i := 0; i < 10; i++ {
		s[i] = i
	}

	fmt.Printf("s len = %d, cap = %d, %v \n", len(s), cap(s), s)

	//for _,v := range s {
	//这里v只是一个索引值的拷贝, 不能用于修改切片的索引值
	//	v = v + 1
	//}

	for i, v := range s {
		s[i] = v + 1
	}
	fmt.Printf("s len = %d, cap = %d, %v \n", len(s), cap(s), s)

}

func TestReSlice(t *testing.T) {
	s := make([]int, 0, 10)

	// len不能超过cap
	// s[start:end],
	for i := 0; i < cap(s); i++ {
		s = s[0 : i+1]
		s[i] = i
		fmt.Printf("s len = %d, cap = %d, %v \n", len(s), cap(s), s)
	}
}

func TestCopy(t *testing.T) {

	src := []int{1, 2, 3}
	dst := make([]int, 10)
	res := copy(dst, src)

	fmt.Printf("src len = %d, cap = %d, %v \n", len(src), cap(src), src)
	fmt.Printf("dst len = %d, cap = %d, %v \n", len(dst), cap(dst), dst)
	fmt.Printf("copy res = %d \n", res)

	//修改目标, 不会改变源数组
	dst[1] = 9999

	fmt.Printf("src len = %d, cap = %d, %v \n", len(src), cap(src), src)
	fmt.Printf("dst len = %d, cap = %d, %v \n", len(dst), cap(dst), dst)

}

func TestAppend(t *testing.T) {

	s := []int{1, 2, 3, 4, 5}
	fmt.Printf("s len = %d, cap = %d, address = %p, %v \n", len(s), cap(s), &s, s)

	//len超过cap, 扩容一倍
	for i := 0; i < 100; i++ {
		//更新slice变量不仅对调用append函数是必要的，实际上对应任何可能导致长度、容量或底层数组变化的操作都是必要的
		//你无法知道append操作之后, 底层数组的内存是否发生变化。
		s = append(s, i)
		fmt.Printf("s len = %d, cap = %d, address = %p, %v \n", len(s), cap(s), &s, s)
	}

}

func TestAppendAddress(t *testing.T) {
	//结论： 两个s指向同一个数组, 但是是不同的切片
	s := make([]int, 3, 10) //数组不变
	//s := []int{1, 2, 3} //由于扩容, 数组也变化了
	fmt.Printf("指向的数组的地址: %p, 切片的地址: %p\n", s, &s)
	add(s)
	fmt.Println(s)
}

func add(s []int) {
	//这样更新是不能改变函数外面的值的, 地址发生变化
	s = append(s, 999)
	fmt.Printf("指向的数组的地址: %p, 切片的地址: %p\n", s, &s)
}

func TestAppendChar(t *testing.T) {

	s := "string"
	beforeChars := []byte(s) // 也可以用copy转数组

	fmt.Printf("before len = %d, cap = %d, address = %p, %v \n", len(beforeChars), cap(beforeChars), &beforeChars, beforeChars)
	afterChars := AppendChar(beforeChars, 'a', 'b', 'c')

	fmt.Printf("after len = %d, cap = %d, address = %p, %v \n", len(afterChars), cap(afterChars), &afterChars, afterChars)

	fmt.Println(string(afterChars))

}

// 用copy实现append, append可以理解为java list的add
func AppendChar(slice []byte, ele ...byte) []byte {
	oldLen := len(slice)
	newLen := oldLen + len(ele)

	newCap := cap(slice)

	if newLen > newCap {
		newCap = newLen * 2
		newSlice := make([]byte, newCap)
		copy(newSlice, slice)
		slice = newSlice
	}

	slice = slice[0:newLen]
	copy(slice[oldLen:newLen], ele)
	return slice
}

// 所有的Go语言函数应该以相同的方式对待nil值的slice和0长度的slice
func TestReverse(t *testing.T) {
	arr := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	reverse(nil)
	reverse(arr[1:9])
	fmt.Println(arr)
}

func reverse(s []int) {
	l := len(s)
	for i := 0; i < l; i++ {
		l--
		s[i], s[l] = s[l], s[i]
	}
}

type Stack struct {
	top []int
}

func (s *Stack) pop() int {
	l := len(s.top) - 1
	head := s.top[l]
	s.top = s.top[0:l]
	return head
}

func (s *Stack) push(v int) {
	s.top = append(s.top, v)
}

func (s *Stack) len() int {
	return len(s.top)
}

func TestStack(t *testing.T) {
	s := new(Stack)

	for i := 0; i < 100; i++ {
		s.push(i)
	}
	for i := 0; i < 100; i++ {
		fmt.Println(s.pop())
	}
	fmt.Printf("len := %v \n", s.len())
}

func remove(slice []int, index int) []int {
	copy(slice[index:], slice[index+1:])
	return slice[:len(slice)-1]
}

// 删除某个元素, 保持原有顺序, 后面的元素依次向前移动一位
func TestRemove(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	slice = remove(slice, 3)
	fmt.Println(slice)
}
