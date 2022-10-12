package inner

import "testing"

//结构体公开
type Inner struct {
	Name    string
	address string
}

// Inner的public方法
func (in *Inner) GetName() string {
	//值类型 不会把修改带到外部；指针类型，会把值带到外部。
	return in.Name
}

// Inner的private方法
func (in *Inner) getName() string {
	return in.Name
}

func TestTestVisibility(t *testing.T) {

	i := new(Inner)
	i.getName()
	i.GetName()
}
