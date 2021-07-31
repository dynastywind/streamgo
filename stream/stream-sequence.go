package stream

import (
	"reflect"

	"github.com/Workiva/go-datastructures/set"
	"github.com/dynastywind/go-stream/stream/operation"
	"github.com/dynastywind/go-stream/util"
)

type SequencialStream struct {
	data        func() []interface{}
	operation   func() []interface{}
	descriptors []OperationDescriptor
}

// Of returns a sequential stream from given data items
//
// @param i	Data items
// @return	A sequential stream
func Of(i ...interface{}) Stream {
	f := func() []interface{} {
		var result []interface{}
		result = append(result, i...)
		return result
	}
	return &SequencialStream{
		data:      f,
		operation: f,
	}
}

// FromArray returns a sequential stream from an interface array
//
// @param arr	An interface array
// @return		A sequential stream
func FromArray(arr []interface{}) Stream {
	f := func() []interface{} {
		return arr
	}
	return &SequencialStream{
		data:      f,
		operation: f,
	}
}

// FromTypedArray returns a sequential stream from a typed array
//
// @param arr	A typed array
// @return		A sequential stream
func FromTypedArray(arr interface{}) Stream {
	f := func() []interface{} {
		return FromTypedArrayToInterfaceArray(arr)
	}
	return &SequencialStream{
		data:      f,
		operation: f,
	}
}

// Concat returns a sequential stream from several streams, either sequential or parallel
//
// @param s	Several streams
// @return	A sequential stream
func Concat(s ...Stream) Stream {
	f := func() []interface{} {
		var result []interface{}
		for _, stream := range s {
			result = append(result, stream.ToArray()...)
		}
		return result
	}
	return &SequencialStream{
		data:      f,
		operation: f,
	}
}

func (s *SequencialStream) AsParallel(routines int) Stream {
	return Transform(OfParallel(routines, s.data()...), s.descriptors)
}

func (s *SequencialStream) AsSequence() Stream {
	return s
}

func (s *SequencialStream) AllMatch(predict func(interface{}) bool) bool {
	arr := s.ToArray()
	if len(arr) == 0 {
		return false
	}
	for _, item := range arr {
		if !predict(item) {
			return false
		}
	}
	return true
}

func (s *SequencialStream) AnyMatch(predict func(interface{}) bool) bool {
	for _, item := range s.ToArray() {
		if predict(item) {
			return true
		}
	}
	return false
}

func (s *SequencialStream) Count() int {
	return len(s.ToArray())
}

func (s *SequencialStream) Distinct(hash func(interface{}) string) Stream {
	return &SequencialStream{
		data: s.data,
		operation: func() []interface{} {
			set := set.New()
			var result []interface{}
			for _, item := range s.ToArray() {
				id := hash(item)
				if !set.Exists(id) {
					set.Add(id)
					result = append(result, item)
				}
			}
			return result
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    DISTINCT,
			params: []interface{}{hash},
		}),
	}
}

func (s *SequencialStream) Filter(filter func(interface{}) bool) Stream {
	return &SequencialStream{
		data: s.data,
		operation: func() []interface{} {
			return operation.Filter(s.ToArray(), filter)
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    FILTER,
			params: []interface{}{filter},
		}),
	}
}

func (s *SequencialStream) FilterOrdered(filter func(interface{}) bool) Stream {
	return &SequencialStream{
		data: s.data,
		operation: func() []interface{} {
			return operation.Filter(s.ToArray(), filter)
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    FILTER_ORDERED,
			params: []interface{}{filter},
		}),
	}
}

func (s *SequencialStream) FindAny() *util.Optional {
	return operation.FindAny(s.ToArray())
}

func (s *SequencialStream) FindFirst() *util.Optional {
	return operation.FindFirst(s.ToArray())
}

func (s *SequencialStream) FlatMap(mapper func(interface{}) []interface{}) Stream {
	return &SequencialStream{
		data: s.data,
		operation: func() []interface{} {
			return operation.DoFlatMap(s.ToArray(), mapper)
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    FLAT_MAP,
			params: []interface{}{mapper},
		}),
	}
}

func (s *SequencialStream) FlatMapOrdered(mapper func(interface{}) []interface{}) Stream {
	return &SequencialStream{
		data: s.data,
		operation: func() []interface{} {
			return operation.DoFlatMap(s.ToArray(), mapper)
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    FLAT_MAP_ORDERED,
			params: []interface{}{mapper},
		}),
	}
}

func (s *SequencialStream) ForEach(consumer func(interface{})) {
	for _, item := range s.ToArray() {
		consumer(item)
	}
}

func (s *SequencialStream) IsParallel() bool {
	return false
}

func (s *SequencialStream) Limit(limit int) Stream {
	return &SequencialStream{
		data: s.data,
		operation: func() []interface{} {
			return operation.Limit(s.ToArray(), limit)
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    LIMIT,
			params: []interface{}{limit},
		}),
	}
}

func (s *SequencialStream) Map(mapper func(interface{}) interface{}) Stream {
	return &SequencialStream{
		data: s.data,
		operation: func() []interface{} {
			return operation.DoMap(s.ToArray(), mapper)
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    MAP,
			params: []interface{}{mapper},
		}),
	}
}

func (s *SequencialStream) MapOrdered(mapper func(interface{}) interface{}) Stream {
	return &SequencialStream{
		data: s.data,
		operation: func() []interface{} {
			return operation.DoMap(s.ToArray(), mapper)
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    MAP_ORDERED,
			params: []interface{}{mapper},
		}),
	}
}

func (s *SequencialStream) Max(less func(interface{}, interface{}) bool) *util.Optional {
	return operation.MaxOrMin(s.ToArray(), less, true)
}

func (s *SequencialStream) Min(less func(interface{}, interface{}) bool) *util.Optional {
	return operation.MaxOrMin(s.ToArray(), less, false)
}

func (s *SequencialStream) NoneMatch(predict func(interface{}) bool) bool {
	for _, item := range s.ToArray() {
		if predict(item) {
			return false
		}
	}
	return true
}

func (s *SequencialStream) Peek(peeker func(interface{})) Stream {
	return &SequencialStream{
		data: s.data,
		operation: func() []interface{} {
			arr := s.ToArray()
			for _, item := range arr {
				peeker(item)
			}
			return arr
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    PEEK,
			params: []interface{}{peeker},
		}),
	}
}

func (s *SequencialStream) Reduce(init interface{}, reducer func(acc, cur interface{}) interface{}) interface{} {
	return operation.Reduce(s.ToArray(), init, reducer)
}

func (s *SequencialStream) ReduceOptional(reducer func(acc, cur interface{}) interface{}) *util.Optional {
	return operation.ReduceOptional(s.ToArray(), reducer)
}

func (s *SequencialStream) ReduceCombine(init interface{}, reducer func(interface{}, interface{}) interface{}, combiner func(interface{}, interface{}) interface{}) interface{} {
	return operation.ReduceCombine(s.ToArray(), init, reducer, combiner)
}

func (s *SequencialStream) Reverse() Stream {
	return &SequencialStream{
		data: s.data,
		operation: func() []interface{} {
			arr := s.ToArray()
			length := len(arr)
			for i := 0; i < length>>1; i++ {
				arr[i], arr[length-i-1] = arr[length-i-1], arr[i]
			}
			return arr
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag: REVERSE,
		}),
	}
}

func (s *SequencialStream) Skip(skip int) Stream {
	return &SequencialStream{
		data: s.data,
		operation: func() []interface{} {
			return operation.Skip(s.ToArray(), skip)
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    SKIP,
			params: []interface{}{skip},
		}),
	}
}

func (s *SequencialStream) Sorted(less func(prev, next interface{}) bool) Stream {
	return &SequencialStream{
		data: s.data,
		operation: func() []interface{} {
			return util.HeapSort(s.ToArray(), less)
		},
		descriptors: append(s.descriptors, OperationDescriptor{
			tag:    SORTED,
			params: []interface{}{less},
		}),
	}
}

func (s *SequencialStream) ToArray() []interface{} {
	return s.operation()
}

func (s *SequencialStream) ToMap(keyMapper func(interface{}) interface{}, valueMapper func(interface{}) interface{}) map[interface{}]interface{} {
	return operation.ToMap(s.ToArray(), keyMapper, valueMapper)
}

func (s *SequencialStream) ToTypedArray(t reflect.Type) reflect.Value {
	return operation.ToTypedArray(s.ToArray(), t)
}

func (s *SequencialStream) ToTypedMap(t reflect.Type, keyMapper func(interface{}) interface{}, valueMapper func(interface{}) interface{}) reflect.Value {
	return operation.ToTypedMap(s.ToArray(), t, keyMapper, valueMapper)
}
