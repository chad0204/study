package effective

import (
	"fmt"
	"testing"
)

func TestFor(t *testing.T) {

	arr := []int{1, 2, 3}

	//style 1
	for i := 0; i < len(arr); i++ {
		fmt.Println(i, arr[i])
	}

	//style 2
	for i, v := range arr {
		fmt.Println(i, v)
	}

}

type V struct {
	i   int
	str string
}

func TestVar(t *testing.T) {
	//零值
	var s string
	var i int
	var f float64
	var v V
	var p *V // 任何类型的指针的零值都是nil
	fmt.Println(s)
	fmt.Println(i)
	fmt.Println(f)
	fmt.Println(v.i)
	fmt.Println(p)

	//类型推导
	//var a, b, c = 1, 1.0, "str"
	//简短声明
	a, b, c := 1, 1.0, "str"

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
}

/**


引用, 声明出来用来表示内存数据的名称, 也称为变量。引用本身不产生副本, 只是给特定地址的内存起一个别名。对引用的操作就是对原内存数据的操作。
引用的值, 内存存储的数据。引用的地址, 内存的地址。可分为指针引用和值引用。

指针是有副本的, 指针也可以被引用。指针引用表示的是内存地址而不是内存数据。

指针的类型是 *T, 指针包含了一个变量的地址, 函数可以通过指针参数改变变量的值
& 取址符 对变量取址就是指针。结果是指向的对象的地址。对指针取址, 还是指针, 结果是指向的指针的地址。
* 取值 对指针进行取值, 就是通过指针存储的地址, 找到地址上的内存数据。

指针作为参数和值作为参数的区别：
不管指针引用还是值引用, 作为参数在函数中传递时都会copy一份。
1.值引用copy一份值, 两个值副本地址不同, 副本的内容相同。
2.指针引用copy一份指针, 两个指针副本地址不同, 副本内容相同。但是由于指针可以进行取值,
而取值操作是通过指针保存的地址进行的, 两个指针副本存储的地址是相同的, 指向同一份内存。
所以对指针引用保存的地址指向的内存进行操作会影响原来的内存（就是同一块内存）。


go中 变量是变量 指针是指针, 变量保存的是引用的对象的地址（变量本身并没有地址）, 指针保存的是指针自己的地址。
java中 变量就是指针, 指针就是变量, 都是保存引用的对象的地址。java




参数变量和返回值变量, 都是局部变量。如果分配在栈上, 随着方法执行结束弹出。
对象如果没有发生逃逸(没有被外部引用或者全局变量引用), 可以直接分配在栈上。 发生逃逸的对象, 需要分配在堆上。
使用指针作为入参或者返回值, 会被认为是逃逸变量, 外部可以通过指针访问到这个变量的, 所以会被分配到堆上。

如果将指向短生命周期对象的指针保存到具有长生命周期的对象中，特别是保存到全局变量时，会阻止对短生命周期对象的垃圾回收（从而可能影响程序的性能）


*/
func TestPointVar(t *testing.T) {
	var i = 999
	var p *int = &i   //p 是i指向的999的地址
	var pp **int = &p // pp是999的地址的地址

	fmt.Printf("original p address: %v \n", p)
	fmt.Printf("original pp address: %v \n", pp)
	fmt.Println(pp)
	//changeVar(i)
	changePoint(&i)

	//point依然指向变量n的地址, 由于getPoint退出, 栈内存释放, 所以必须将变量n移动到堆上
	point := getPoint()
	//逃逸变量
	*point = 2

}

func getPoint() *int {
	n := 1
	return &n
}

/**

变量i保存的值是999, 被copy, i和i'地址不同(&i), 保存的内容相同都是999. i'分配在栈上, 不用GC

赋值操作, 就是通过i'直接修改i'上存储的999, 不会改变i的999

*/
func changeVar(i int) {
	var p *int = &i
	var pp **int = &p
	fmt.Printf("value args p address: %v \n", p)   //值i的地址
	fmt.Printf("value args pp address: %v \n", pp) //值i的地址的地址
	i = 1000
}

/**
变量p保存的值是变量i的地址,被copy, p和p'地址不同(&p), 保存的内容相同都是变量i的地址. p'会将i移动到堆上, 需要GC

赋值操作, 首先需要取值, 就是通过p'存储的地址找到变量i的存储的999, 修改i的999
*/
func changePoint(p *int) {
	var pp **int = &p
	fmt.Printf("point args p address: %v \n", p)   //指针的值（变量i的地址）
	fmt.Printf("point args pp address: %v \n", pp) //指针的地址
	*p = 1000
}

//new的返回值是指针类型
// 用new创建对象和指针, 区别在于省去了变量名i_。new函数只是一个语法糖。
func TestNew(t *testing.T) {

	var i_ int
	p_ := &i_
	fmt.Println(*p_)

	p := new(int)
	fmt.Println(*p)

	v := new(V)
	fmt.Println(*v)
}

//赋值
func TestAssignment(t *testing.T) {
	n := 1
	//编译错误
	//m := n++

	n++
	m := n
	fmt.Println(m)

	//元组赋值
	x, y := 1, 999
	x, y = y, x
	x, y = x+y, x
	fmt.Println(x)

}

type Celsius float64    //摄氏度
type Fahrenheit float64 //华氏度
type Unknower string    //未知

const (
	AbsoluteZeroC Celsius = -273.15 // 绝对零度
	FreezingC     Celsius = 0       // 结冰点温度
	BoilingC      Celsius = 100     // 沸水温度
)

func TestType(t *testing.T) {
	var c Celsius
	var f Fahrenheit
	//var un Unknower
	//fmt.Println(c == f)//error
	fmt.Println(CToF(c) == f) // ture/false 可以比较
	fmt.Println(float64(c))
	//fmt.Println(float64(un)) // error

}

func CToF(celsius Celsius) Fahrenheit {
	//T(x), 将x转为T类型, 只有x和T的底层类型一样才能转换
	//fahrenheit := float64(celsius)
	fahrenheit := Fahrenheit(celsius)
	return fahrenheit + 100
}

const (
	Sunday = iota //0
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

const (
	B = 1 << (10 * iota)
	KB
	MB
	GB
	TB
)

func TestIota(t *testing.T) {
	fmt.Println(Saturday)
	fmt.Println(MB)

}
