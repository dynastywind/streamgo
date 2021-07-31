package operation

import (
	"github.com/Workiva/go-datastructures/queue"
)

func Filter(arr []interface{}, filter func(interface{}) bool) []interface{} {
	var result []interface{}
	for _, item := range arr {
		if filter(item) {
			result = append(result, item)
		}
	}
	return result
}

type FilterResultWrapper struct {
	index  int
	result bool
	data   interface{}
}

func (wrapper FilterResultWrapper) Compare(other queue.Item) int {
	return wrapper.index - other.(FilterResultWrapper).index
}

func FilterInParallel(arr []interface{}, num int, filter func(interface{}) bool, ordered bool) []interface{} {
	length := len(arr)
	ch := make(chan FilterResultWrapper, length)
	result := make([]interface{}, 0)
	pq := queue.NewPriorityQueue(length, true)
	for i, item := range arr {
		go func(index int, data interface{}) {
			ch <- FilterResultWrapper{
				index:  index,
				result: filter(data),
				data:   data,
			}
		}(i, item)
		if (i+1)%num == 0 || i == length-1 {
			for j := 0; j <= i%num; j++ {
				wrapper := <-ch
				if wrapper.result {
					if ordered {
						pq.Put(wrapper)
					} else {
						result = append(result, wrapper.data)
					}
				}
			}
		}
	}
	if ordered {
		items, err := pq.Get(pq.Len())
		if err != nil {
			panic("Should succeed with priority queue get")
		}
		for _, wrapper := range items {
			result = append(result, wrapper.(FilterResultWrapper).data)
		}
	}
	return result
}
