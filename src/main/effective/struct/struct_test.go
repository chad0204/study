package _struct

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"study/src/main/effective/struct/inner"
	"testing"
	"unsafe"
)

/**


结构体有零值, 零值就是各个成员变量的零值


*/

// 数组也可以看成一种结构体, 下标相当于字段, 都是值类型
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
						|	value:		 |	    |   id: 234    | |   value:           | |  vp | | next  |
					    | {data: "fp"}	 |	    |		       | |  {data: "second"}  | |     | |       |
						+————————————————+		+ —————————————+ +——————————————————— + + ————+ +—————— +


*/

// 递归结构体
type Node struct {
	id    int
	value Value
	vp    *Value
	//next Node//递归结构体, 不能嵌套自己。否则无法计算node的大小, 而指针的大小是固定的。（该限制同样适用于数组）
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

func TestFactory(t *testing.T) {
	//结构体没有构造器, 可以通过工厂方法来创建, 通过小写不能包外导出, 禁止new(Type)
	instance := inner.NewInstance(233, "pc")
	fmt.Println(instance)
	fmt.Println(instance.GetId())
}

type Foo map[string]int

type Bar struct {
	key   string
	value int
}

func TestNewAndMake(t *testing.T) {
	// make用于引用类型, new用于值类型
	//b := make(Bar)//make用于值类型编译报错
	f := make(Foo)
	f["key"] = 10
	fmt.Println(f)

	bar := new(Bar)
	bar.key = "key"
	bar.value = 10
	fmt.Println(bar)

	foo := new(Foo) //foo是一个指向nil的指针, 仍然需要*foo = map[string]int{}来分配内存
	fmt.Printf("%v, %p \n", foo, foo)
	//(*foo)["key"] = 10// new用于引用类型，运行时报错assignment to entry in nil map
}

type TagFoo struct {
	field1 string `tag:"name"`
	field2 int    `tag:"id"`
}

func TestTag(t *testing.T) {

	tf := new(TagFoo)
	tf.field1 = "tag1"
	tf.field2 = 233
	fmt.Printf("tf = %v,%p, *tf = %v,%p \n", tf, tf, *tf, &*tf)

	//TypeOf的入参是值。 反射获取tag
	field1 := reflect.TypeOf(*tf).Field(0)
	field2 := reflect.TypeOf(*tf).Field(1)

	fmt.Printf("%v, %v", field1.Tag, field2.Tag)
}

type Employee struct {
	id   int
	name string
}

func EmpById(id int) *Employee {
	employee := &Employee{234, "point_emp"}
	return employee
}

func EmpByIdV2(id int) Employee {
	employee := Employee{234, "value_emp"}
	fmt.Printf("%v", &employee)
	return employee
}

// 返回值
func TestReturn(t *testing.T) {

	employee := EmpById(1)
	employee.name = "_changed"
	fmt.Println(&employee)

	employee1 := EmpByIdV2(1)
	employee1.name = "_changed"
	fmt.Println(&employee1)

}

/**-----------------------------------面向对象---------------------------------------**/

type Father struct {
	firstName string
	address   string
}

type Son struct {
	id     int
	string `json:"匿名字段"`
	Father `json:"内嵌结构体"`
	f      Father
}

//内嵌结构体, 可以实现继承！！！

// 匿名字段, 内嵌匿名结构体
func TestAnonymous(t *testing.T) {
	s := new(Son)
	s.id = 10
	s.string = "匿名字段"
	s.firstName = "abc" //直接访问内嵌结构体的字段
	s.address = "LA"
	s.f.firstName = "not anonymous" // 非匿名就不是内嵌结构体
	fmt.Println(s)

	son := &Son{
		id:     1,
		string: "anonymousStr",
		Father: Father{
			firstName: "abc",
			address:   "LA",
		},
		f: Father{
			firstName: "abc",
			address:   "LA",
		},
	}
	fmt.Println(son)
}

type V struct {
	name string
}

type K struct {
	V     // 内嵌类型不需要指针
	vp *V //指针类型的域
	v  V  // 结构体内存是连续的, 通过k指针修改v的内容, 是有效的
}

// 指针结构体字段、值结构体字段、和内嵌结构体字段
func TestStructField(t *testing.T) {
	k := &K{V{"anonymous field"}, &V{"point field"}, V{"value field"}}
	fmt.Printf("before: %v \n", k)
	//changeK(k)// 都生效
	changeV(k.V)   // 无效
	changeV(k.v)   // 无效
	changeV(*k.vp) // 无效

	fmt.Printf("after : %v \n", k)

}

func changeV(v V) {
	v.name = v.name + "change"
}

func changeK(k *K) {
	// 通过外结构体的指针, 可以修改内部值、指针、匿名的结构体体域
	k.name = "anonymous field change" //ok
	k.vp.name = "point field change"  //ok
	//k.v.name = "value field change"   //ok
	(&(k.v)).name = "value field change"
}

type A struct {
	a string
}

type B struct {
	a, b string
}

// 同级别命名冲突
type C struct {
	A
	B
}

// 内外层命名冲突
type D struct {
	B
	b string
}

// 内嵌结构体, 命名冲突
func TestDuplicateName(t *testing.T) {
	//同级别命名冲突
	c := &C{
		A{"a"},
		B{"a", "b"},
	}
	fmt.Println(c.A.a)
	//fmt.Println(c.a)//compiler error

	d := &D{
		B{
			"a", "out b",
		},
		"in b",
	}
	fmt.Println(d.b)   // 内层会覆盖外层
	fmt.Println(d.B.b) //可以这样

}

type Local struct {
	x int
	y int
}

func (l *Local) String() string {
	return "{x = " + strconv.Itoa(l.x) + ", y = " + strconv.Itoa(l.y) + "}"
}

// toString ,fmt.Printf()的%v和fmt.Printf()都会调用类型的String()方法
func TestString(t *testing.T) {

	l := new(Local)
	l.x = 10
	l.y = 15
	fmt.Println(l)
	fmt.Printf("%v \n", l)
	fmt.Printf("%T \n", l)
	fmt.Printf("%#v \n", l) // 生成的结构可以作为go代码生成

}

type StackNode struct {
	prev *StackNode
	next *StackNode
	data int
}

type Stack struct {
	size int
	head *StackNode
	tail *StackNode
}

// push 到尾巴
func (s *Stack) push(data int) int {
	//init
	if s.head == nil || s.tail == nil {
		s.head = &StackNode{}
		s.tail = &StackNode{}
		s.head.next = s.tail
		s.tail.prev = s.head
	}
	//append
	node := &StackNode{
		prev: s.tail.prev,
		next: s.tail,
		data: data,
	}
	s.tail.prev.next = node
	s.tail.prev = node
	s.size++
	return s.size
}

// 从尾巴取
func (s *Stack) pop() int {
	if s.tail.prev == s.head {
		return 0
	}

	data := s.tail.prev.data

	prev := s.tail.prev
	s.tail.prev = prev.prev
	prev.prev.next = s.tail
	prev.prev = nil
	prev.next = nil
	s.size--

	return data
}

func TestStack(t *testing.T) {
	stack := new(Stack)
	stack.push(1)
	stack.push(2)
	stack.push(3)
	fmt.Println(stack.size)

	fmt.Println(stack.pop())
	fmt.Println(stack.pop())
	fmt.Println(stack.pop())
	fmt.Println(stack.size)

}

type Point struct {
	x, y int
}

type Circle struct {
	Point
	Radius int
}

type Wheel struct {
	Circle
	Spokes int // 轮子辐条数
}

func (c *Circle) Len() float64 {
	return math.Pi * float64(2) * float64(c.Radius)
}

func TestWheel(t *testing.T) {
	wheel := Wheel{
		Circle: Circle{
			Point: Point{
				x: 1,
				y: 2,
			},
			Radius: 3,
		},
		Spokes: 2,
	}
	change(wheel)
	fmt.Println(wheel.Circle.y)
	fmt.Printf("fmt = %#v, v = %v \n", wheel, wheel)
}

func change(wheel Wheel) {
	wheel.Circle.y = 233
}
