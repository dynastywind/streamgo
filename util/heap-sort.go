package util

func HeapSort(arr []interface{}, less func(a, b interface{}) bool) []interface{} {
	for i := len(arr)/2 - 1; i >= 0; i-- {
		heapify(arr, i, len(arr), less)
	}
	for j := len(arr) - 1; j >= 0; j-- {
		arr[0], arr[j] = arr[j], arr[0]
		heapify(arr, 0, j, less)
	}
	return arr
}

func heapify(arr []interface{}, i int, length int, less func(a, b interface{}) bool) {
	pivot := arr[i]
	for k := 2*i + 1; k < length; k = 2*k + 1 {
		if k+1 < length && less(arr[k], arr[k+1]) {
			k++
		}
		if less(pivot, arr[k]) {
			arr[i] = arr[k]
			i = k
		} else {
			break
		}
	}
	arr[i] = pivot
}
