package solid

import (
	"fmt"
	"testing"
)

// 匿名结构体实现继承 C -> B -> A
type A struct {
}

type B struct {
	A
}

type C struct {
	B
}

type BaseType struct {
}

func (t *BaseType) GetValue() *B {
	fmt.Println("BaseType GetValue")
	return new(B)
}

type SubTypeA struct {
	BaseType
}

func (t *SubTypeA) GetValue() *A {
	fmt.Println("SubTypeA GetValue")
	return new(A)
}

type SubTypeB struct {
	BaseType
}

func (t *SubTypeB) GetValue() *C {
	fmt.Println("SubTypeB GetValue")
	return new(C)
}

func TestLsp(t *testing.T) {
	//b := new(BaseType)
	//b := new(SubTypeA)
	b := new(SubTypeB)
	b.GetValue()

}
