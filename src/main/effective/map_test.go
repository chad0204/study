package effective

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {

	//数组、切片和结构体不能作为 key，但是指针和接口类型可以。
	//如果要用结构体作为 key 可以提供 Key() 和 Hash() 方法，这样可以通过结构体的域计算出唯一的数字或者字符串的 key
	var m map[string]int

	fmt.Println(m)

}
