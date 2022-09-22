package algorithm

func bubbleSort(arr []int) {

	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr)-1; j++ {
			if arr[j] < arr[i] {
				tmp := arr[i]
				arr[i] = arr[j]
				arr[j] = tmp
			}
		}
	}
	/**
	7845261
	7452618
	4526178
	5261478

	*/

}
