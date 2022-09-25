package algorithm

import (
	"fmt"
	"testing"
)

func TestBubbleSort(t *testing.T) {

	arr := []int{1, 5, 3, 4, 2, 99, 88, 77}
	bubbleSort(arr) // 值传递不行的
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
