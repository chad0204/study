package algorithm

func bubbleSort(nums []int) {
	var complete bool
	for i := 0; i < len(nums)-1; i++ { //次数
		complete = true
		for j := 0; j < len(nums)-1-i; j++ { //交换
			if nums[j] > nums[j+1] {
				tmp := nums[j]
				nums[j] = nums[j+1]
				nums[j+1] = tmp
				complete = false
			}
		}
		if complete == true {
			return
		}
	}
}

func mergeSort(nums []int) []int {
	if len(nums) <= 1 {
		return nums
	}

	left := mergeSort(nums[:len(nums)/2])
	right := mergeSort(nums[len(nums)/2:])

	var res []int
	lx := 0
	rx := 0
	for lx < len(left) && rx < len(right) {
		if left[lx] < right[rx] {
			res = append(res, left[lx])
			lx++
		} else {
			res = append(res, right[rx])
			rx++
		}
	}

	for lx < len(left) {
		res = append(res, left[lx])
		lx++
	}
	for rx < len(right) {
		res = append(res, right[rx])
		rx++
	}
	return res
}

func quickSort(start, end int, nums []int) {
	if start >= end {
		return
	}
	lx := start
	rx := end
	p := nums[lx]
	for lx < rx {
		for nums[rx] >= p && lx < rx {
			rx--
		}
		for nums[lx] <= p && lx < rx {
			lx++
		}
		tmp := nums[lx]
		nums[lx] = nums[rx]
		nums[rx] = tmp
	}

	nums[start] = nums[lx]
	nums[lx] = p

	quickSort(start, lx-1, nums)
	quickSort(lx+1, end, nums)
}

func insertionSort(nums []int) {

}
