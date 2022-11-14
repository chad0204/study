package effective

import (
	"fmt"
	"testing"
)

//当调用一个函数的时候，函数的每个调用参数将会被赋值给函数内部的参数变量，所以函数参数变量接收的是一个复制的副本，并不是原始调用的变量。
//形参是引用类型（指针、slice、map、chan）, 实参和形参副本都可以取到原来的变量。形参数是值类型, 那形参就是实参值的副本, 修改形参不影响实参。

//必要时给返回值取合适的名字, 让函数更清晰
func Split(path string) (dir, file string, err error) {
	//return dir, file
	return //返回值取名后, 可以省略
}

func TestReturn(t *testing.T) {
	fmt.Println(Split(""))
}

/**
异常处理的4种策略:
	1. 向上抛, 或者append后向上抛
	2. 重试并打印日志
    3. 终止程序打印日志
    4. 不终止程序打印日志
*/

func preHandler(arg string) {
	fmt.Println("pre: " + arg)
}

func postHandler(arg string) {
	fmt.Println("post: " + arg)
}

func exec(arg string, pre, post func(args string)) {
	if pre != nil {
		pre(arg)
	}

	fmt.Println("exec...")

	if post != nil {
		post(arg)
	}
}

/**
ref: https://www.sulinehk.com/post/golang-closure-details/#1-%E5%A3%B0%E6%98%8E%E6%96%B0%E5%8F%98%E9%87%8F
*/

//函数值(也叫闭包) 就是声明一个变量 指向一个匿名的函数. 函数值属于引用类型, 函数值不可比较
func TestFuncVar(t *testing.T) {
	var before = preHandler
	var after = postHandler

	exec("a", before, after)

	f := func(i int) int {
		return i * i
	}

	g := func(i int) int {
		return i * i
	}

	fmt.Println(f(2))
	fmt.Println(g(2))

}

//返回类型是func() int
func squares() func() int {
	var i int // 变量 x 的值是被保留的, 不管外部函数squares退出与否,  说明匿名函数是引用类型
	return func() int {
		i++
		return i * i
	}
}

//将函数值作为返回值
func TestReturnFunc(t *testing.T) {
	s := squares()
	//squares返回后，变量i仍然隐式的存在于s中
	//s持有i
	fmt.Println(s())
	fmt.Println(s())
	fmt.Println(s())
	fmt.Println(s())
	//fmt.Println(squares()())//这里每次都是新的函数调用, 新的i值

	n := 1
	f := func() {
		//这里的n就是外层的n
		fmt.Println(n)
	}
	f()
}

func Test1(t *testing.T) {
	//申明一个"函数值"数组
	var funcs []func()

	//add func var
	for i := 0; i < 10; i++ {
		funcs = append(funcs, func() {
			fmt.Println(i)
		})
	}

	//执行函数数组的每个函数, 输出的都是10（不是9, i++之后才退出）
	for i := 0; i < 10; i++ {
		funcs[i]()
	}

}

func method(i int) (res int) {
	defer func() {
		res++
	}()
	return i * i
}

func execute() {
	doExec()
	fmt.Println("exec over")
}

func doExec() {
	//如果没有recover, 那么程序就会在panic后面结束
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("defer print: %v \n", p)
		}
	}()
	/*
		do biz
		...
	*/

	//error
	panic("error")
}

//defer func and recover
func TestDefer(t *testing.T) {
	fmt.Println(method(10))

	//有panic异常时, go程序会退出，使用recover可以恢复
	execute()
}
