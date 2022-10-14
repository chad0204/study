package inner

import "testing"

//结构体公开
type Inner struct {
	Name    string
	address string
}

type instance struct {
	//这里也可以小写, 但是要提供访问方法
	id   int
	name string
}

//Person 导出类型 不导出字段
type Person struct {
	name string
}

// Name getter
func (p *Person) Name() string {
	return p.name
}

// SetName setter
func (p *Person) SetName(name string) {
	p.name = name
}

//NewInstance factory
func NewInstance(id int, name string) *instance {
	return &instance{
		id,
		name,
	}
}

func (i *instance) GetId() int {
	return i.id
}

func (i *instance) GetName() string {
	return i.name
}

// GetName Inner的public方法
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
