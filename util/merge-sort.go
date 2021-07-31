package util

func MergeSort(data []interface{}, less func(prev, cur interface{}) bool) []interface{} {
	length := len(data)
	if length < 2 {
		return data
	}
	half := length >> 1
	ch1 := make(chan []interface{})
	ch2 := make(chan []interface{})
	go func(arr []interface{}) {
		ch1 <- MergeSort(arr, less)
	}(data[:half])
	go func(arr []interface{}) {
		ch2 <- MergeSort(arr, less)
	}(data[half:])
	return merge(<-ch1, <-ch2, less)
}

func merge(a, b []interface{}, less func(prev, cur interface{}) bool) []interface{} {
	var result []interface{}
	i := 0
	j := 0
	lenA := len(a)
	lenB := len(b)
	for i < lenA && j < lenB {
		if less(a[i], b[j]) {
			result = append(result, a[i])
			i++
		} else {
			result = append(result, b[i])
			j++
		}
	}
	for ; i < lenA; i++ {
		result = append(result, a[i])
	}
	for ; j < lenB; j++ {
		result = append(result, b[j])
	}
	return result
}
