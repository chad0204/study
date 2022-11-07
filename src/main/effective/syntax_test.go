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
	fmt.Println(s)
	fmt.Println(i)
	fmt.Println(f)
	fmt.Println(v.i)

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
& 取址符 对值取址就是地址, 对指针取址就是指针的地址
* 取值 对指针进行取值, 就是通过指针存储的地址, 找到地址上的内存数据。

指针作为参数和值作为参数的区别：
不管指针引用还是值引用, 作为参数在函数中传递时都会copy一份。
1.值引用copy一份值, 两个值副本地址不同, 副本的内容相同。
2.指针引用copy一份指针, 两个指针副本地址不同, 副本内容相同。但是由于指针可以进行取值,
而取值操作是通过指针保存的地址进行的, 两个指针副本存储的地址是相同的, 指向同一份内存。
所以对指针引用保存的地址指向的内存进行操作会影响原来的内存（就是同一块内存）。

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
