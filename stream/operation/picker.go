package operation

func Limit(arr []interface{}, limit int) []interface{} {
	if len(arr) > limit {
		return arr[:limit]
	}
	return arr
}

func Skip(arr []interface{}, skip int) []interface{} {
	if len(arr) > skip {
		return arr[skip:]
	}
	return []interface{}{}
}
