package effective

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestReflect(t *testing.T) {
	var val float64 = 3.2
	fmt.Println("type: ", reflect.TypeOf(val))

	x := reflect.ValueOf(val)
	fmt.Println("value: ", x)
	fmt.Println("kind: ", x.Kind())
	fmt.Println("kind is float64: ", x.Kind() == reflect.Float64)
	fmt.Println("value: ", x.Float())

	i := x.Interface()
	fmt.Println("interface: ", i)
	fmt.Println("interface value: ", i.(float64))

}

//反射修改值
func TestReflectWrite(t *testing.T) {

	var val float64 = 3.2
	of := reflect.ValueOf(val)
	fmt.Printf("of type: %T, of value: %v \n", of, of)
	//无法修改, of只是val值的一个拷贝reflect.Value.SetFloat using unaddressable value
	if of.CanSet() {
		of.SetFloat(3.99)
	}

	v := reflect.ValueOf(&val)
	fmt.Printf("v type: %T, v value: %v \n", v, v)
	if v.CanSet() {
		v.SetFloat(3.98)
	}
	v = v.Elem()
	fmt.Printf("v type: %T, v value: %v \n", v, v)
	if v.CanSet() {
		v.SetFloat(3.97)
	}

	fmt.Println(val)

}

//具体的类型
type SpecificType struct {
	str string
	I   int
	f   float64
}

// 大写 值转递才可以反射
func (s SpecificType) String() string {
	return s.str + strconv.Itoa(s.I)
}

var secretVal interface{} = SpecificType{"strField", 1, 3.2}

// 反射结构体
func TestReflectStruct(t *testing.T) {

	fmt.Println("type: ", reflect.TypeOf(secretVal))

	x := reflect.ValueOf(secretVal)
	fmt.Println("value: ", x)
	fmt.Println("kind: ", x.Kind()) // struct
	fmt.Println("kind is struct: ", x.Kind() == reflect.Struct)
	fmt.Println("value: ", x.Field(0))

	i := x.Interface()
	fmt.Println("interface: ", i)
	fmt.Println("interface value: ", i.(SpecificType))

	//打印field
	for i := 0; i < x.NumField(); i++ {
		fmt.Printf("filed%v is %v \n", i, x.Field(i))
	}

	//反射方法
	for i := 0; i < x.NumMethod(); i++ {
		fmt.Printf("method%v is %v \n", i, x.Method(i).Call(nil))
	}

}

func TestReflectWriteStruct(t *testing.T) {

	val := SpecificType{"str", 1, 2.3}

	elem := reflect.ValueOf(&val).Elem()

	//using value obtained using unexported field 小写的字段域不可修改
	//elem.Field(0).SetString("change")
	elem.Field(1).SetInt(999)

	fmt.Println(val)
	fmt.Println(elem)

}

//反射与printf
func TestReflectAndPrintf(t *testing.T) {

	val := SpecificType{"value", 1, 3.2}

	fmt.Printf("value is %v \n", val)

}
