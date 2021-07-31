package stream

import (
	"reflect"

	"github.com/dynastywind/go-stream/stream/operation"
	"github.com/dynastywind/go-stream/util"
)

type ParallelStream struct {
	data        func() []interface{}
	operation   func() []interface{}
	descriptors []OperationDescriptor
	routines    int
}

// OfParallel returns a parallel stream from given data items
//
// @param routines	Number of goroutines
// @param i			Data items
// @return			A parallel stream
func OfParallel(routines int, i ...interface{}) Stream {
	if routines < 1 {
		panic("Parallel version need go routines greater than 1. Otherwise please use sequential version for better performance.")
	}
	f := func() []interface{} {
		var result []interface{}
		return append(result, i...)
	}
	return &ParallelStream{
		data:      f,
		operation: f,
		routines:  routines,
	}
}

// FromArrayParallel returns a parallel stream from an interface array
//
// @param routines	Number of goroutines
// @param arr		An interface array
// @return			A parallel stream
func FromArrayParallel(routines int, arr []interface{}) Stream {
	f := func() []interface{} {
		return arr
	}
	return &ParallelStream{
		data:      f,
		operation: f,
		routines:  routines,
	}
}

// FromTypedArrayParallel returns a parallel stream from a typed array
//
// @param routines	Number of goroutines
// @param arr		A typed array
// @return			A parallel stream
func FromTypedArrayParallel(routines int, arr interface{}) Stream {
	f := func() []interface{} {
		return FromTypedArrayToInterfaceArray(arr)
	}
	return &ParallelStream{
		data:      f,
		operation: f,
		routines:  routines,
	}
}

// ConcatAsParallel returns a parallel stream from several streams, either sequential or parallel
//
// @param routines	Number of goroutines
// @param s			Several streams
// @return			A parallel stream
func ConcatAsParallel(routines int, s ...Stream) Stream {
	f := func() []interface{} {
		var result []interface{}
		for _, stream := range s {
			result = append(result, stream.ToArray()...)
		}
		return result
	}
	return &ParallelStream{
		data:      f,
		operation: f,
		routines:  routines,
	}
}

func (s *ParallelStream) AsParallel(routines int) Stream {
	return s
}

func (s *ParallelStream) AsSequence() Stream {
	return Transform(Of(s.data()...), s.descriptors)
}

func (s *ParallelStream) AllMatch(predict func(interface{}) bool) bool {
	arr := s.ToArray()
	length := len(arr)
	if length == 0 {
		return false
	}
	ch := make(chan bool, length)
	for i, item := range arr {
		go func(data interface{}) {
			ch <- predict(data)
		}(item)
		if (i+1)%s.routines == 0 || i == length-1 {
			for j := 0; j <= i%s.routines; j++ {
				if !<-ch {
					return false
				}
			}
		}
	}
	return true
}

func (s *ParallelStream) AnyMatch(predict func(interface{}) bool) bool {
	arr := s.ToArray()
	length := len(arr)
	ch := make(chan bool, length)
	for i, item := range arr {
		go func(data interface{}) {
			ch <- predict(data)
		}(item)
		if (i+1)%s.routines == 0 || i == length-1 {
			for j := 0; j <= i%s.routines; j++ {
				if <-ch {
					return true
				}
			}
		}
	}
	return false
}

func (s *ParallelStream) Count() int {
	return len(s.ToArray())
}

type HashPair struct {
	hash string
	item interface{}
}

func (s *ParallelStream) Distinct(hash func(interface{}) string) Stream {
	return &ParallelStream{
		data: s.data,
		operation: func() []interface{} {
			arr := s.ToArray()
			length := len(arr)
			ch := make(chan HashPair, length)
			m := make(map[string]interface{})
			var result []interface{}
			for i, item := range arr {
				go func(data interface{}) {
					ch <- HashPair{
						hash: hash(data),
						item: data,
					}
				}(item)

				if (i+1)%s.routines == 0 || i == length-1 {
					for j := 0; j <= i%s.routines; j++ {
						pair := <-ch
						m[pair.hash] = pair.item
					}
				}
			}
			for _, v := range m {
				result = append(result, v)
			}
			return result
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    DISTINCT,
			params: []interface{}{hash},
		}),
		routines: s.routines,
	}
}

func (s *ParallelStream) Filter(filter func(interface{}) bool) Stream {
	return &ParallelStream{
		data: s.data,
		operation: func() []interface{} {
			return operation.FilterInParallel(s.ToArray(), s.routines, filter, false)
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    FILTER,
			params: []interface{}{filter},
		}),
		routines: s.routines,
	}
}

func (s *ParallelStream) FilterOrdered(filter func(interface{}) bool) Stream {
	return &ParallelStream{
		data: s.data,
		operation: func() []interface{} {
			return operation.FilterInParallel(s.ToArray(), s.routines, filter, true)
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    FILTER_ORDERED,
			params: []interface{}{filter},
		}),
		routines: s.routines,
	}
}

func (s *ParallelStream) FindAny() *util.Optional {
	return operation.FindAny(s.ToArray())
}

func (s *ParallelStream) FindFirst() *util.Optional {
	return operation.FindFirst(s.ToArray())
}

func (s *ParallelStream) FlatMap(mapper func(interface{}) []interface{}) Stream {
	return &ParallelStream{
		data: s.data,
		operation: func() []interface{} {
			return operation.DoFlatMapInParallel(s.ToArray(), s.routines, mapper, false)
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    FLAT_MAP,
			params: []interface{}{mapper},
		}),
		routines: s.routines,
	}
}

func (s *ParallelStream) FlatMapOrdered(mapper func(interface{}) []interface{}) Stream {
	return &ParallelStream{
		data: s.data,
		operation: func() []interface{} {
			return operation.DoFlatMapInParallel(s.ToArray(), s.routines, mapper, true)
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    FLAT_MAP_ORDERED,
			params: []interface{}{mapper},
		}),
		routines: s.routines,
	}
}

func (s *ParallelStream) ForEach(consumer func(interface{})) {
	operation.ForEachParallel(s.ToArray(), s.routines, consumer)
}

func (s *ParallelStream) IsParallel() bool {
	return true
}

func (s *ParallelStream) Limit(limit int) Stream {
	return &ParallelStream{
		data: s.data,
		operation: func() []interface{} {
			return operation.Limit(s.ToArray(), limit)
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    LIMIT,
			params: []interface{}{limit},
		}),
		routines: s.routines,
	}
}

func (s *ParallelStream) Map(mapper func(interface{}) interface{}) Stream {
	return &ParallelStream{
		data: s.data,
		operation: func() []interface{} {
			return operation.DoMapInParallel(s.ToArray(), s.routines, mapper, false)
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    MAP,
			params: []interface{}{mapper},
		}),
		routines: s.routines,
	}
}

func (s *ParallelStream) MapOrdered(mapper func(interface{}) interface{}) Stream {
	return &ParallelStream{
		data: s.data,
		operation: func() []interface{} {
			return operation.DoMapInParallel(s.ToArray(), s.routines, mapper, true)
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    MAP_ORDERED,
			params: []interface{}{mapper},
		}),
		routines: s.routines,
	}
}

func (s *ParallelStream) Max(less func(interface{}, interface{}) bool) *util.Optional {
	return operation.MaxOrMin(s.ToArray(), less, true)
}

func (s *ParallelStream) Min(less func(interface{}, interface{}) bool) *util.Optional {
	return operation.MaxOrMin(s.ToArray(), less, false)
}

func (s *ParallelStream) NoneMatch(predict func(interface{}) bool) bool {
	arr := s.ToArray()
	length := len(arr)
	ch := make(chan bool, length)
	for i, item := range arr {
		go func(data interface{}) {
			ch <- predict(item)
		}(item)
		if (i+1)%s.routines == 0 || i == length-1 {
			for j := 0; j <= i%s.routines; j++ {
				if <-ch {
					return false
				}
			}
		}
	}
	return true
}

func (s *ParallelStream) Peek(peeker func(interface{})) Stream {
	return &ParallelStream{
		data: s.data,
		operation: func() []interface{} {
			arr := s.ToArray()
			operation.ForEachParallel(arr, s.routines, peeker)
			return arr
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    PEEK,
			params: []interface{}{peeker},
		}),
		routines: s.routines,
	}
}

func (s *ParallelStream) Reduce(init interface{}, reducer func(interface{}, interface{}) interface{}) interface{} {
	return operation.Reduce(s.ToArray(), init, reducer)
}

func (s *ParallelStream) ReduceCombine(init interface{}, reducer func(interface{}, interface{}) interface{}, combiner func(interface{}, interface{}) interface{}) interface{} {
	return operation.ReduceCombine(s.ToArray(), init, reducer, combiner)
}

func (s *ParallelStream) ReduceOptional(reducer func(interface{}, interface{}) interface{}) *util.Optional {
	return operation.ReduceOptional(s.ToArray(), reducer)
}

func (s *ParallelStream) Reverse() Stream {
	return &ParallelStream{
		data: s.data,
		operation: func() []interface{} {
			arr := s.ToArray()
			length := len(arr)
			ch := make(chan int, length)
			half := length >> 1
			for i := 0; i < half; i++ {
				go func(index int) {
					arr[index], arr[length-index-1] = arr[length-index-1], arr[index]
					ch <- 1
				}(i)
				if (i+1)%s.routines == 0 || i == half-1 {
					for j := 0; j <= i%s.routines; j++ {
						<-ch
					}
				}
			}
			return arr
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag: REVERSE,
		}),
		routines: s.routines,
	}
}

func (s *ParallelStream) Skip(skip int) Stream {
	return &ParallelStream{
		data: s.data,
		operation: func() []interface{} {
			return operation.Skip(s.ToArray(), skip)
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    SKIP,
			params: []interface{}{skip},
		}),
		routines: s.routines,
	}
}

func (s *ParallelStream) Sorted(less func(interface{}, interface{}) bool) Stream {
	return &ParallelStream{
		data: s.data,
		operation: func() []interface{} {
			return util.MergeSort(s.ToArray(), less)
		},
		routines: s.routines,
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    SORTED,
			params: []interface{}{less},
		}),
	}
}

func (s *ParallelStream) ToArray() []interface{} {
	return s.operation()
}

func (s *ParallelStream) ToMap(keyMapper func(interface{}) interface{}, valueMapper func(interface{}) interface{}) map[interface{}]interface{} {
	return operation.ToMap(s.ToArray(), keyMapper, valueMapper)
}

func (s *ParallelStream) ToTypedArray(t reflect.Type) reflect.Value {
	return operation.ToTypedArray(s.ToArray(), t)
}

func (s *ParallelStream) ToTypedMap(t reflect.Type, keyMapper func(interface{}) interface{}, valueMapper func(interface{}) interface{}) reflect.Value {
	return operation.ToTypedMap(s.ToArray(), t, keyMapper, valueMapper)
}
