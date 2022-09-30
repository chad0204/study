package effective

import (
	"fmt"
	"testing"
)

func TestPoint(t *testing.T) {

	str := "original"
	//值传递
	changeStr(str)
	fmt.Println(str)

	//值传递
	changeStrV2(&str)
	fmt.Println(str)

	p := Person{
		name: "original",
	}
	//值传递 (java这样就是引用传递了)
	changeName(p)
	fmt.Println(p.name)

	//引用传递
	changeNameRef(&p)
	fmt.Println(p.name)

	//总结： java基本类型是值传递,对象类型是引用传递; go基本类型是值传递, 但对象类型也是值传递（如果需要引用传递可以通过指针）,
	//go也提供了一些引用类型的对象, 比如slice、map

}

func changeStr(str string) {
	str = "changed"
}

func changeStrV2(str *string) {
	changed := "changed"
	str = &changed
}

func changeName(p Person) {
	p.name = "changed"
}

func changeNameRef(p *Person) {
	p.name = "changed"
}

type Person struct {
	name string
}
