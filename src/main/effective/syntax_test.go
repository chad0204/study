package effective

import (
	"fmt"
	"testing"
)

func TestFor(t *testing.T) {

	arr := []int{1, 2, 3}

	//style 1
	for i := 0; i < len(arr); i++ {
		fmt.Println(i, arr[i])
	}

	//style 2
	for i, v := range arr {
		fmt.Println(i, v)
	}

}
