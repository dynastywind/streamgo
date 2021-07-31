package operation

import "github.com/dynastywind/go-stream/util"

func MaxOrMin(arr []interface{}, less func(interface{}, interface{}) bool, max bool) *util.Optional {
	sorted := util.MergeSort(arr, less)
	length := len(sorted)
	if length == 0 {
		return util.OfEmpty()
	}
	if max {
		return util.OfNillable(sorted[length-1])
	} else {
		return util.OfNillable(sorted[0])
	}
}
