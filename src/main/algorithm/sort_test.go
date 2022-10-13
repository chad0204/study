package algorithm

import (
	"fmt"
	"testing"
)

func TestBubbleSort(t *testing.T) {

	arr := []int{1, 5, 3, 4, 2, 99, 88, 77} //没有定义长度，是slice, 引用类型
	BubbleSort(arr)                         // 数组值传递不行的，这里入参数组没带长度，是slice
	for _, v := range arr {
		fmt.Println(v)
	}
}
func TestMergeSort(t *testing.T) {
	arr := []int{1, 5, 3, 4, 2, 99, 88, 77}
	sorted := mergeSort(arr)
	for _, v := range sorted {
		fmt.Println(v)
	}
}

func TestQuickSort(t *testing.T) {
	arr := []int{1, 5, 3, 4, 2, 99, 88, 77}
	quickSort(0, len(arr)-1, arr)
	for _, v := range arr {
		fmt.Println(v)
	}
}
