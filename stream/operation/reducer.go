package operation

import "github.com/dynastywind/go-stream/util"

func Reduce(arr []interface{}, init interface{}, reducer func(acc, cur interface{}) interface{}) interface{} {
	length := len(arr)
	if length <= 0 {
		return init
	}
	result := init
	for i := 0; i < length; i++ {
		result = reducer(result, arr[i])
	}
	return result
}

func ReduceCombine(arr []interface{}, init interface{}, reducer func(interface{}, interface{}) interface{}, combiner func(interface{}, interface{}) interface{}) interface{} {
	length := len(arr)
	if length <= 0 {
		return init
	}
	result := init
	for i := 0; i < length; i++ {
		result = combiner(result, reducer(result, arr[i]))
	}
	return result
}

func ReduceOptional(arr []interface{}, reducer func(acc, cur interface{}) interface{}) *util.Optional {
	length := len(arr)
	if length == 0 {
		return util.OfEmpty()
	} else if length == 1 {
		return util.OfNillable(arr[0])
	}
	result := arr[0]
	for i := 1; i < length; i++ {
		result = reducer(result, arr[i])
	}
	return util.OfNillable(result)
}
