package effective

import (
	"fmt"
	"testing"
)

//数组也可以看成一种结构体, 下标相当于字段, 都是值类型
type Demo struct {
	filed1 string
	filed2 string
	filed3 string
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

}
