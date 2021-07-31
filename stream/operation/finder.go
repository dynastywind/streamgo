package operation

import (
	"math/rand"

	"github.com/dynastywind/go-stream/util"
)

func FindAny(arr []interface{}) *util.Optional {
	length := len(arr)
	if length == 0 {
		return util.OfEmpty()
	}
	return util.OfNillable(arr[rand.Intn(length)])
}

func FindFirst(arr []interface{}) *util.Optional {
	length := len(arr)
	if length == 0 {
		return util.OfEmpty()
	}
	return util.OfNillable(arr[0])
}
