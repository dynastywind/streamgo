package stream

import (
	"fmt"
	"reflect"
)

func Transform(stream Stream, descriptors []OperationDescriptor) Stream {
	for _, desc := range descriptors {
		switch desc.tag {
		case DISTINCT:
			stream = stream.Distinct(desc.params[0].(func(interface{}) string))
		case FILTER:
			stream = stream.Filter(desc.params[0].(func(interface{}) bool))
		case FLAT_MAP:
			stream = stream.FlatMap(desc.params[0].(func(interface{}) []interface{}))
		case LIMIT:
			stream = stream.Limit(desc.params[0].(int))
		case MAP:
			stream = stream.Map(desc.params[0].(func(interface{}) interface{}))
		case PEEK:
			stream = stream.Peek(desc.params[0].(func(interface{})))
		case REVERSE:
			stream = stream.Reverse()
		case SKIP:
			stream = stream.Skip(desc.params[0].(int))
		case SORTED:
			stream = stream.Sorted(desc.params[0].(func(interface{}, interface{}) bool))
		default:
			panic(fmt.Sprintf("Unsupported operation type found: %v", desc.tag))
		}
	}
	return stream
}

func FromTypedArrayToInterfaceArray(arr interface{}) []interface{} {
	if reflect.TypeOf(arr).Kind() != reflect.Slice {
		panic("arr should be of slice type")
	}
	v := reflect.ValueOf(arr)
	length := v.Len()
	result := make([]interface{}, length)
	for i := 0; i < length; i++ {
		result[i] = v.Index(i).Interface()
	}
	return result
}
