package effective

import (
	"fmt"
	"testing"
)

type Foo struct {
}

type InstanceService struct {
	name string

	// not ok
	//func (i *InstanceService) f(name string) {
	//	i.name = name
	//}
}

// 值类型 不会把方法内的修改带到外部；指针类型，会把方法内的修改带到外部。不同的接收类型是一个方法集
func (i *InstanceService) modifyFieldPoint(name string) {
	i.name = name
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

}
