package effective

import (
	"fmt"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {

	//map是引用类型
	//数组、切片和结构体不能作为 key，但是指针和接口类型可以。
	//如果要用结构体作为 key 可以提供 Key() 和 Hash() 方法，这样可以通过结构体的域计算出唯一的数字或者字符串的 key
	var m map[string]int
	fmt.Println(m == nil)
	m2 := map[string]int{}
	fmt.Println(m2 == nil)
	//不要使用 new，永远用 make 来构造 map
	newMap := new(map[string]float32)
	//newMap["key"] = 4.5//error null point
	fmt.Println(newMap)

	m = map[string]int{"a": 1, "b": 2, "c": 3}
	m["d"] = 10
	fmt.Println(m)

	m1 := make(map[string]string)
	m1["key1"] = "value1"
	fmt.Println(m1)
	changeMap(m1) //引用类型
	fmt.Println(m1)

}

func changeMap(m map[string]string) {
	m["key1"] = "value*"
	m["key2"] = "value1"
}

func TestCap(t *testing.T) {
	m := make(map[string]int, 10)
	fmt.Println(len(m))
	//超过10会扩容
	for i := 0; i < 11; i++ {
		m["key"+strconv.Itoa(i)] = i
	}

	//size for pair
	fmt.Println(len(m))
	fmt.Println(m["key10"])

	funcMap := map[int]func() string{
		1: func() string {
			return "func1"
		},
		2: func() string {
			return "func2"
		},
		3: func() string {
			return "func3"
		},
	}
	fmt.Println(funcMap)

	sliceMap := map[string][]int{
		"first":  {1, 2, 3},
		"second": {2, 3, 3},
	}

	sliceMap1 := make(map[string][]int)
	sliceMap1["first"] = []int{1, 2, 3}
	sliceMap2 := make(map[string]*[]int)
	sliceMap2["first"] = &[]int{1, 2, 3}

	sliceMap1["first"][0] = 999
	(*sliceMap2["first"])[0] = 999

	fmt.Println(sliceMap)
	fmt.Println(sliceMap1)
	fmt.Println(sliceMap2["first"])

	domainMap := map[string]Domain{
		"a": {name: "original"},
	}

	//无法修改值类型
	domain := domainMap["a"]
	domain.name = "changed"
	//domainMap["a"].name = "233" + domainMap["a"].name//error
	fmt.Println(domainMap["a"])

	domainMap1 := map[string]*Domain{
		"b": {name: "original"},
	}

	//可以修改引用值类型
	domain1 := domainMap1["b"]
	domain1.name = "changed"
	//上面的值类型无法下面操作
	domainMap1["b"].name = "233" + domainMap1["b"].name
	fmt.Println(domainMap1["b"])

}

type Domain struct {
	name string
}

func TestPresetAndDel(t *testing.T) {
	var stringMap = map[string]string{
		"key1": "",
		"key2": "value2",
		"key3": "value3",
	}

	//无法区分空值和不存在的值
	fmt.Println(stringMap["key1"])
	fmt.Println(stringMap["key4"] == stringMap["key1"])

	//使用多返回值
	v, ifPresent := stringMap["key1"]
	fmt.Printf("map contains %v: %v, value = %v \n", "key1", ifPresent, v)

	if _, ifPresent := stringMap["key1"]; ifPresent {
		delete(stringMap, "key1")
	}

	v1, ifPresent1 := stringMap["key1"]
	fmt.Printf("map contains %v: %v, value = %v \n", "key1", ifPresent1, v1)

}

func TestMapForRange(t *testing.T) {
	m := map[string]int{
		"k1": 1,
		"k2": 2,
		"k3": 3,
	}

	for k, v := range m {
		fmt.Printf("k = %v, v = %v \n", k, v)
	}

	for k := range m {
		fmt.Printf("k = %v \n", k)
	}

	for _, v := range m {
		fmt.Printf("v = %v \n", v)
	}
}

func TestMapSlice(t *testing.T) {

	mapSlice := []map[string]int{map[string]int{
		"k": 1,
	}, map[string]int{
		"k1": 1,
		"k2": 2,
	}, map[string]int{
		"k1": 1,
	}}

	fmt.Println(mapSlice[1])

	//使用make初始化
	mapSliceV2 := make([]map[string]int, 10)
	//不要使用_, v, 因为v是copy出的副本
	for i := range mapSliceV2 {
		mapSliceV2[i] = make(map[string]int, 1)
		mapSliceV2[i]["key"] = 999
	}

}
