package effective

import "fmt"

//对于一些组件, 不确定会出现什么问题, 所以用recover捕获, 当作普通地解析错误抛出

type Syntax struct {
}

func ParseJSON(input string) (s *Syntax, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("JSON: internal error: %v", p)
		}
	}()
	//parser...
	return nil, err
}
