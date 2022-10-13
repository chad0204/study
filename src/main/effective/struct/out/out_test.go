package out

import (
	"study/src/main/effective/struct/inner"
	"testing"
)

//不可以引入xx_test
//不可以引入小写的成员(常量、变量、方法、结构体字段)
func TestVisibility(t *testing.T) {
	i := &inner.Inner{
		Name: "fdafdsa",
		//address: 不可见
	}
	i.GetName()
	//i.getName() 不可见

}
