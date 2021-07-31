package operation

import "reflect"

func ToMap(arr []interface{}, keyMapper func(interface{}) interface{}, valueMapper func(interface{}) interface{}) map[interface{}]interface{} {
	var result map[interface{}]interface{} = make(map[interface{}]interface{})
	for _, item := range arr {
		result[keyMapper(item)] = valueMapper(item)
	}
	return result
}

func ToTypedMap(arr []interface{}, t reflect.Type, keyMapper func(interface{}) interface{}, valueMapper func(interface{}) interface{}) reflect.Value {
	result := reflect.MakeMap(t)
	for _, item := range arr {
		result.SetMapIndex(reflect.ValueOf(keyMapper(item)), reflect.ValueOf(valueMapper(item)))
	}
	return result
}

func ToTypedArray(arr []interface{}, t reflect.Type) reflect.Value {
	result := reflect.MakeSlice(reflect.SliceOf(t), 0, len(arr))
	for _, item := range arr {
		result = reflect.Append(result, reflect.ValueOf(item))
	}
	return result
}
