package algorithm

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBubbleSort(t *testing.T) {

	arr := []int{1, 5, 3, 4}

	fmt.Println(reflect.TypeOf(arr))

	bubbleSort(arr) // 值传递不行的

	for i, v := range arr {
		fmt.Println(i, v)
	}
}
