package _struct

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

type InstanceService struct {
	name string

	// 方法没有和数据定义（结构体）混在一起：它们是正交的类型；表示（数据）和行为（方法）是独立的
	// not ok
	//func (i *InstanceService) f(name string) {
	//	i.name = name
	//}
}

// 值类型 不会把方法内的修改带到外部；指针类型，会把方法内的修改带到外部。不同的接收类型是一个方法集
func (i *InstanceService) modifyFieldPoint(name string) {
	i.name = name
}

func (_ *InstanceService) Query() string {
	return "result"
}

// 接收者不同, 方法名可以重复。一个接收者类型不能重载
func (i *Foo) modifyFieldPoint(name string) {

}

func (i InstanceService) modifyFieldValue(name string) {
	i.name = name
}

// 方法, 有接收者的特殊函数
func TestMethod(t *testing.T) {

	i := new(InstanceService)
	i.modifyFieldValue("value")
	fmt.Println(i)

	j := new(InstanceService)
	j.modifyFieldPoint("point")
	fmt.Println(j)

	//接口类型不能有方法

	//不能重载, 既不能有同名的方法。但接收者不同, 方法名可以相同

	//receiver类似oo语言的this

}

// 非结构体接收者, slice的别名类型。go可以给任何命名类型定义方法, 只要这个命名类型的低层类型(低层类型是指[]int, Slice是命名类型)不是指针或者interface
type Slice []int

func (s Slice) Sum() int {
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}

// map的别名类型, 只能用make
type Map map[string]int

func (m Map) appendKey() string {
	keys := ""
	for k, _ := range m {
		keys += k + ","
	}
	return keys
}

// 不能直接给string加方法, 但是可以给别名类型加方法
type MyString struct {
	string
}

// 此方法只针对别名MyString有效
func (s MyString) getLength() int {
	return len(s.string)
}

// 类型和作用在它上面定义的方法必须在同一个包里定义!!!
// 比如给int,string定义方法 编译报cannot define new methods on non-local type int/string
// 可以定义别名类型来解决
func TestAliasTypeMethod(t *testing.T) {

	//Slice是slice的别名类型引用类型, 不能new
	s := make(Slice, 3) //type must be slice, map, or channel
	s[0] = 1
	s[1] = 2
	set(s)
	fmt.Println(s.Sum())
	fmt.Println(Slice{1, 2, 3, 4}.Sum())

	m := make(Map)
	m["aa"] = 1
	m["bb"] = 1
	m["cc"] = 1
	fmt.Println(m.appendKey())

	s2 := new(MyString)
	s2.string = "123456"
	fmt.Println(s2.getLength())
}

func set(s Slice) {
	s[2] = 3
}

type Service struct {
	name string
}

func (s *Service) pointM() string {
	return "point"
}
func (s Service) valueM() string {
	return "value"
}

/**

1. 无论方法的receiver是*T还是T, 都可以通过指针和值进行调用, go编译器都会帮你做类型转换
2. receiver类型的场景, 使用*T指向的都是同一块内存, receiver相当于是调用者的一个别名, 修改会影响调用者。使用T不会影响调用者,
但是需要拷贝内存, 如果调用者本身特别大, 会有内存损耗。

*/
// 指针接收者和值接收者
func TestValueAndPoint(t *testing.T) {

	//如果指针和值都可以, 用指针性能更好, 拷贝一个值代价大概率比拷贝一个指针大
	//如果方法内部要改变接收者的数据, 必须用指针

	// 指针方法和值方法都可以在指针或非指针上被调用, 变量调用方法是不区分变量是值还是指针的, 只要可以寻址就行。（接口不是）

	//指针可以调用值方法和指针方法(go编译器帮我们隐式转换了)
	ps := &Service{}
	ps.pointM()
	ps.valueM() //指针调用值方法, 自动转换
	(*ps).valueM()
	(*ps).pointM() //值调用指针方法, 自动转换

	//值可以调用值方法和指针方法(go编译器帮我们隐式转换了)
	vs := Service{}
	vs.valueM()
	vs.pointM()
	(&vs).valueM()
	(&vs).pointM()

	Service{"biz"}.valueM()
	//Service{"biz"}.pointM() //compile error

}

// 通过内嵌结构体 模拟继承与多态
type Engine struct {
}

func (e *Engine) start() {
	fmt.Println("engine start")
}

func (e *Engine) stop() {
	fmt.Println("engine stop")
}

type Car struct {
	Engine
	wheelNum int
}

func (c *Car) numberOfWheel() int {
	return c.wheelNum
}

type Mercedes struct {
	Car //内嵌结构体不需要用指针
}

// 重写
func (m *Mercedes) start() {
	fmt.Println("Mercedes car start")
}

// 继承
func TestInherit(t *testing.T) {
	m := &Mercedes{Car{Engine{}, 4}}
	m.start()
	m.stop()
	fmt.Printf(strconv.Itoa(m.numberOfWheel()))

}

type Fly struct {
}

type Swim struct {
}

type Duck struct {
	Fly
	Swim
}

//定义结构体变量, 并直接初始化变量
var d = struct {
	Fly
	Swim
	name   string
	weight float64
}{
	name:   "duck",
	weight: 5.4,
}

func (f *Fly) fly() {
	fmt.Println("I can fly!")
}
func (s *Swim) swim() {
	fmt.Println("I can swim!")
}

// 多重继承
func TestMoreInherit(t *testing.T) {
	duck := Duck{}
	duck.fly()
	duck.swim()
	duck.Fly.fly()

	d.fly()
	fmt.Println(d.name)
}

// 和java相比, go的多态是通过组合实现的, 而不是继承, 更加灵活

//cache直接拥有了加锁的能力
var cache = struct {
	sync.Mutex
	mapping map[string]string
}{
	mapping: make(map[string]string),
}

func Lookup(key string) string {
	cache.Lock()
	defer cache.Unlock()
	val := cache.mapping[key]
	return val
}

type Person struct {
	id int
}

func (p *Person) SetId(id int) {
	p.id = id
}

func (p *Person) GetId() int {
	return p.id
}

func (r *Person) Launch() {
	r.id = 111
}

//方法值和方法表达式
func TestMethodV(t *testing.T) {
	person := Person{233}
	//set就是方法值
	set := person.SetId
	set(1999)
	fmt.Println(person.GetId())

	//setId就是函数表达式
	setId := (*Person).SetId
	//第一个参数是接收器, 需要区分指针和值
	setId(&person, 2000)
	fmt.Println(person.GetId())

}

func TestBitSlice(t *testing.T) {
	for i := 0; i < 100; i++ {
		//word是商, bit是取模
		word, bit := i/64, uint(i%64)
		fmt.Printf("i = %v, word = %v, bit = %v \n", i, word, bit)
	}

}
