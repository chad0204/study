package effective

import (
	"fmt"
	"testing"
)

//指针的值是地址, 指向对象。引用也有自己的地址。
//对象的值是对象属性。对象的地址是指针。

/**

类型: 引用类型和值类型（基本类型）
传递: 针对函数调用,
	值传递, 实参通过拷贝将自身内容传递给形参, 形参实际是实参的一个副本, 对这个副本的修改不影响原先的实参内容;
	引用传递, 实参在函数调用时只将自己的地址传递给形参, 通过形参的内容（即实参的地址）,可以操作实参的内容。


引用类型或者传递指针，函数调用时开销比较小，因为copy一份指针副本，比copy整个对象的开销小。而且传递对象，无法改变原来的对象内容。

指针：引用传递, 需要函数传递后,重新赋值对象的属性。会走gc（gc只关心指针，扫码对象是否包含指针）, 且可以修改指针指向的对象的属性, 带来复杂性（需要明确是否可以被修改）。
不用指针：不用走gc。需要copy整个对象的内存。如果对象比较小, 不用指针可以减少gc, 如果用了指针, gc需要扫描的路径会变长。


*/
func TestPoint(t *testing.T) {
	//基本类型都是值类型, java也是

	//值传递
	str1 := "original"
	changeStr(str1)
	fmt.Println("基本类型是值传递: original -> " + str1)

	//值传递
	str2 := "original"
	changeStrV2(&str2)
	fmt.Println(str2)
	fmt.Println("基本类型是值传递: original -> " + str1)
}

func TestStruct(t *testing.T) {
	//值传递
	p := Person{
		name: "original",
	}
	fmt.Printf("before p address:%p, p value:%v \n", &p, p)
	// java传递p时就是引用传递（内外地址一样）, 不需要使用指针, 所以java是引用传递。
	changeName(p)
	fmt.Println("p是值传递: original -> " + p.name)

	//值传递
	p1 := Person{
		name: "original",
	}
	changeNameV2(p1)
	fmt.Println("p1是值传递: original -> " + p.name)

	//⭐值传递（指针本身是值传递，但是修改引用指向的对象的内存, 会改变原来对象的属性）
	// map和slice是引用类型, 不需要地址符。而java所有对象都是引用类型。
	p2 := Person{
		name: "original",
	}
	point := &p
	fmt.Printf("before point address:%p, point value:%p \n", &point, point)
	changeNameRef(point)
	fmt.Println("p2的指针是值传递，但是修改当前指针的指向的对象内存会改变原指针指向的对象: original -> " + p2.name)

	//值传递（引用地址（指针）作为值是值传递, 所以修改这个指针指向其他地址，并不会改变原来的指针）
	p3 := Person{
		name: "original",
	}
	changeNameRefV2(&p3)
	fmt.Println("p3的指针是值传递，修改指针的地址值不会改变原来的指针地址: original -> " + p2.name)

	//总结： java基本类型是值传递,对象类型是引用传递; go基本类型是值传递, 但对象类型也是值传递（如果需要引用传递可以通过指针）,
	//go也提供了一些引用类型的对象, 比如slice、map、chan

}

func TestMapChange(t *testing.T) {
	m := make(map[string]int)
	m["key1"] = 999
	fmt.Printf("map before info. pointAddress: %p, pointValue: %p \n", &m, m)
	changeMapValue(m)
	fmt.Printf("map info: %v \n", m)
}

func changeMapValue(m map[string]int) {
	//传递的是引用, 引用copy了一份副本, 但是引用的值（也就是指向的地址）是同一个。java自定义类型都是引用类型, go自定义类型是值类型。
	fmt.Printf("map after  info. pointAddress: %p, pointValue: %p \n", &m, m)
	m["key1"] = 111
}

func changeStr(str string) {
	str = "changed"
}

func changeStrV2(str *string) {
	//copy了引用的值（地址值）, 修改引用指向的值不会改变外部引用的地址值。
	changed := "changed"
	str = &changed
}

func changeName(p Person) {
	fmt.Printf("after p address:%p, p value:%v \n", &p, p)
	//这里copy出一份对象副本, 形参和实参的地址不一样, 也就是形参的p和实参的p已经不是一个对象, 修改当前对象不会改变外面的对象
	p.name = "changed"
}

//演示直接赋值而不是修改属性
func changeNameV2(p Person) {
	//这里copy出一份对象副本, 形参和实参的地址不一样, 修改当前副本不会改变外面的实参对象
	person := Person{
		name: "changed",
	}
	p = person
}

func changeNameRef(point *Person) {
	//这里会copy一份指针的副本（两个指针的地址不一样）, 虽然形参是实参的副本, 但是两个指针的值一样, 指向同一个对象的地址, 修改引用地址指向的内存会改变改内存
	fmt.Printf("after  point address:%p, point value:%p \n", &point, point)
	point.name = "changed"
}

//演示直接赋值而不是修改属性
func changeNameRefV2(p *Person) {
	//这里会copy一份指针的副本, 由于是两个不同的指针, 形参指向新的地址, 不会改变原来的实参指向的地址
	p = &Person{
		name: "changed",
	}
	p.name = p.name + "aaa"
}

type Person struct {
	name string
}
