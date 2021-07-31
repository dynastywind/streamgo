package operation

func DoMap(arr []interface{}, mapper func(interface{}) interface{}) []interface{} {
	var result []interface{}
	for _, item := range arr {
		result = append(result, mapper(item))
	}
	return result
}

type GeneralResultWrapper struct {
	index int
	data  interface{}
}

func DoMapInParallel(arr []interface{}, num int, mapper func(interface{}) interface{}, ordered bool) []interface{} {
	length := len(arr)
	orderedResult := make([]interface{}, length)
	unorderedResult := make([]interface{}, 0)
	ch := make(chan GeneralResultWrapper, length)
	for i, item := range arr {
		go func(index int, data interface{}) {
			ch <- GeneralResultWrapper{
				index: index,
				data:  mapper(data),
			}
		}(i, item)
		if (i+1)%num == 0 || i == length-1 {
			for j := 0; j <= i%num; j++ {
				wrapper := <-ch
				if ordered {
					orderedResult[wrapper.index] = wrapper.data
				} else {
					unorderedResult = append(unorderedResult, wrapper.data)
				}
			}
		}
	}
	if ordered {
		return orderedResult
	}
	return unorderedResult
}

type ArrayResultWrapper struct {
	index int
	data  []interface{}
}

func DoFlatMap(arr []interface{}, mapper func(interface{}) []interface{}) []interface{} {
	var result []interface{}
	for _, item := range arr {
		result = append(result, mapper(item)...)
	}
	return result
}

func DoFlatMapInParallel(arr []interface{}, num int, mapper func(interface{}) []interface{}, ordered bool) []interface{} {
	length := len(arr)
	container := make([][]interface{}, length)
	ch := make(chan ArrayResultWrapper, length)
	var result []interface{}
	for i, item := range arr {
		go func(index int, data interface{}) {
			ch <- ArrayResultWrapper{
				index: index,
				data:  mapper(data),
			}
		}(i, item)
		if (i+1)%num == 0 || i == length-1 {
			for j := 0; j <= i%num; j++ {
				wrapper := <-ch
				if ordered {
					container[wrapper.index] = wrapper.data
				} else {
					result = append(result, wrapper.data...)
				}
			}
		}
	}
	if ordered {
		for _, item := range container {
			result = append(result, item...)
		}
	}
	return result
}
