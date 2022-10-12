package _struct

import (
	"fmt"
	"testing"
	"unsafe"
)

//数组也可以看成一种结构体, 下标相当于字段, 都是值类型
type Demo struct {
	filed1 string
	filed2 string
	filed3 string
}

/**


node := &Node{
	id: 233,
	value: Value {
		data: "first",
	},
	vp: &Value {
		data: "fp",
	},
	next: &Node{
		id: 234,
		value: Value{data: "second"},
	},
}

结构体字段内存是连续的, 指针字段就是指针内存(指针内存固定), 值字段就是值内存
java的区别是,java引用都是栈内存,也就是下图中的value:{data:"first"}也是指针

node
 \
  \
  _\/
 + —————————————+ +————————————————— + + ————+ +—————— +
 |   id: 233    | |   value:         | |  vp | | next  |
 |              | |   {data:"first"} | |     | |       |
 + —————————————+ +————————————————— + + ————+ +—————— +
                                          /       \
                                         /         \
                                       \/_	       _\/
						+————————————————+		+ —————————————+ +——————————————————— + + ————+ +—————— +
						|	value:		 |	    |   id: 233    | |   value:           | |  vp | | next  |
					    | {data: "fp"}	 |	    |		       | |  {data: "second"}  | |     | |       |
						+————————————————+		+ —————————————+ +——————————————————— + + ————+ +—————— +


*/

//递归结构体
type Node struct {
	id    int
	value Value
	vp    *Value
	//next Node//递归结构体, 不能嵌套自己。否则无法计算Demo的大小,指针的大小是固定的。
	next *Node
}

type Value struct {
	data string
}

func TestStructDefine(t *testing.T) {

	//1 值类型
	var demo Demo
	demo.filed1 = "f1"
	demo.filed2 = "f2"
	changeDemo(demo)
	fmt.Println(demo)

	//2 new 返回值是指针
	d := new(Demo)
	d.filed3 = "f3"
	fmt.Println(d)
	changeDemoV2(d)
	fmt.Println(d)

	var d1 *Demo = new(Demo)
	var d2 *Demo
	d2 = new(Demo)
	changeDemoV2(d1)
	changeDemoV2(d2)
	fmt.Println(d1)
	fmt.Println(d2)

}

func changeDemo(demo Demo) {
	demo.filed1 = "changed"
}

func changeDemoV2(demo *Demo) {
	demo.filed1 = "changed"
}

func TestInit(t *testing.T) {

	demo := Demo{"f1", "f2", "f3"}
	fmt.Println(demo)

	//常规 &Type== new(Type)
	demo1 := &Demo{filed1: "f1", filed2: "f2", filed3: "f3"} //不必顺序
	demo2 := &Demo{"f1", "f2", "f3"}                         //按顺序
	demo3 := &Demo{filed3: "f3"}                             //部分字段
	fmt.Println(demo1)
	fmt.Println(demo2)
	fmt.Println(demo3)

	node := new(Node)
	node.id = 233
	node.next = new(Node)
	node.value = *new(Value)
	node.vp = new(Value)

	n := &Node{
		id: 233,
		value: Value{
			data: "233",
		},
		vp: &Value{
			data: "233",
		},
		next: &Node{
			id:    234,
			value: Value{data: "234"},
		},
	}
	fmt.Println(n)

	size := unsafe.Sizeof(*n)

	fmt.Println(size)

	//结构体没有构造器

}

type Instance struct {
	id   int
	name string
}

func TestFactory(t *testing.T) {
	//结构体没有构造器, 可以通过工厂方法来创建
	instance := NewInstance(233, "pc")

	//禁止使用new

	fmt.Println(instance)
}

func NewInstance(id int, name string) *Instance {
	return &Instance{
		id,
		name,
	}
}

func newInstance(id int, name string) *Instance {
	return &Instance{
		id,
		name,
	}
}
