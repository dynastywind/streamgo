package stream_test

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/Workiva/go-datastructures/list"
	"github.com/dynastywind/go-stream/stream"
	"github.com/dynastywind/go-stream/util"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = ginkgo.Describe("Test if every function works well", func() {
	ginkgo.Context("Stream construction test", func() {
		ginkgo.When("Construct a stream with items", func() {
			ginkgo.It("should construct a parallel stream", func() {
				arr := stream.OfParallel(2, 1, 2, 3, 4).ToArray()
				gomega.Expect(arr).To(gomega.Equal([]interface{}{1, 2, 3, 4}))
			})
		})
		ginkgo.When("Construct a stream with array", func() {
			ginkgo.It("should construct a parallel stream from an interface array", func() {
				arr := stream.FromArrayParallel(2, []interface{}{1, 2, 3, 4}).ToArray()
				gomega.Expect(arr).To(gomega.Equal([]interface{}{1, 2, 3, 4}))
			})
		})
		ginkgo.When("Construct a stream with typed array", func() {
			ginkgo.It("should construct a parallel stream from a type array", func() {
				arr := stream.FromTypedArrayParallel(2, []int{1, 2, 3, 4}).ToArray()
				gomega.Expect(arr).To(gomega.Equal([]interface{}{1, 2, 3, 4}))
			})
		})
		ginkgo.When("Construct a stream with concatenation", func() {
			ginkgo.It("should construct a parallel stream from several streams", func() {
				s := stream.ConcatAsParallel(2, stream.Of(1, 2, 3, 4), stream.OfParallel(2, 5, 6, 7, 8))
				gomega.Expect(s.ToArray()).To(gomega.Equal([]interface{}{1, 2, 3, 4, 5, 6, 7, 8}))
				gomega.Expect(s.IsParallel()).To(gomega.BeTrue())
			})
		})
	})
	ginkgo.Context("Stream tranformation test", func() {
		ginkgo.When("Executing AsParallel", func() {
			ginkgo.It("should return itself", func() {
				original := stream.OfParallel(2, 2, 3, 4)
				s := original.AsParallel(2)
				gomega.Expect(s.IsParallel()).To(gomega.BeTrue())
				gomega.Expect(s).To(gomega.Equal(original))
			})
		})
		ginkgo.When("Executing AsSequence", func() {
			ginkgo.It("should return a sequential stream", func() {
				s := stream.OfParallel(2, 2, 3, 4).AsSequence()
				gomega.Expect(s.IsParallel()).To(gomega.BeFalse())
			})
		})
	})
	ginkgo.Context("Sigle operation test", func() {
		ginkgo.When("Executing AllMatch", func() {
			ginkgo.It("should return true", func() {
				result := stream.OfParallel(2, 1, 2, 3, 4).AllMatch(func(item interface{}) bool {
					return item.(int) > 0
				})
				gomega.Expect(result).To(gomega.BeTrue())
			})
			ginkgo.It("should return false", func() {
				result := stream.OfParallel(2, 1, 2, 3, 4).AllMatch(func(item interface{}) bool {
					return item.(int) > 1
				})
				gomega.Expect(result).To(gomega.BeFalse())
			})
		})
		ginkgo.When("Executing AnyMatch", func() {
			ginkgo.It("should return true", func() {
				result := stream.OfParallel(2, 1, 2, 3, 4).AnyMatch(func(item interface{}) bool {
					return item.(int) > 1
				})
				gomega.Expect(result).To(gomega.BeTrue())
			})
			ginkgo.It("should return false", func() {
				result := stream.OfParallel(2, 1, 2, 3, 4).AnyMatch(func(item interface{}) bool {
					return item.(int) > 4
				})
				gomega.Expect(result).To(gomega.BeFalse())
			})
		})
		ginkgo.When("Executing Count", func() {
			ginkgo.It("should return the length of the data stream", func() {
				count := stream.OfParallel(2, 1, 2, 3, 4).Count()
				gomega.Expect(count).To(gomega.Equal(4))
			})
		})
		ginkgo.When("Executing Distinct", func() {
			ginkgo.It("should return an array of unique values", func() {
				arr := stream.OfParallel(2, 1, 1, 2, 2).Distinct(func(item interface{}) string {
					return strconv.Itoa(item.(int))
				}).ToArray()
				gomega.Expect(arr).To(BagEquals([]interface{}{1, 2}))
			})
		})
		ginkgo.When("Executing Filter", func() {
			ginkgo.It("should return an array of values greater than 2", func() {
				arr := stream.OfParallel(2, 1, 2, 3, 4).Filter(func(item interface{}) bool {
					return item.(int) > 2
				}).ToArray()
				gomega.Expect(arr).To(BagEquals([]interface{}{3, 4}))
			})
		})
		ginkgo.When("Executing FilterOrdered", func() {
			ginkgo.It("should return an array of values greater than 2", func() {
				arr := stream.OfParallel(2, 1, 2, 3, 4).FilterOrdered(func(item interface{}) bool {
					return item.(int) > 2
				}).ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{3, 4}))
			})
		})
		ginkgo.When("Executing FindAny", func() {
			ginkgo.It("should find a value", func() {
				result := stream.OfParallel(2, 1, 2, 3, 4).FindAny()
				gomega.Expect(result.IsPresent()).To(gomega.BeTrue())
			})
			ginkgo.It("should find nothing", func() {
				result := stream.OfParallel(1).FindAny()
				gomega.Expect(result.IsPresent()).To(gomega.BeFalse())
			})
		})
		ginkgo.When("Executing FindFirst", func() {
			ginkgo.It("should find the first value", func() {
				result := stream.OfParallel(2, 1, 2, 3, 4).FindFirst()
				gomega.Expect(result).To(gomega.Equal(util.Of(1)))
			})
			ginkgo.It("should find nothing", func() {
				result := stream.OfParallel(1).FindFirst()
				gomega.Expect(result.IsPresent()).To(gomega.BeFalse())
			})
		})
		ginkgo.When("Executing FlatMap", func() {
			ginkgo.It("should return corresponding words in a row", func() {
				m := make(map[int][]string)
				m[1] = []string{"one", "un"}
				m[2] = []string{"two", "deux"}
				m[3] = []string{"three", "trois"}
				m[4] = []string{"four", "quatre"}
				arr := stream.OfParallel(2, 1, 2, 3, 4).FlatMap(func(item interface{}) []interface{} {
					var result []interface{}
					for _, strings := range m[item.(int)] {
						result = append(result, strings)
					}
					return result
				}).ToArray()
				gomega.Expect(arr).To(BagEquals([]interface{}{"one", "un", "two", "deux", "three", "trois", "four", "quatre"}))
			})
		})
		ginkgo.When("Executing FlatMapOrdered", func() {
			ginkgo.It("should return corresponding words in a row", func() {
				m := make(map[int][]string)
				m[1] = []string{"one", "un"}
				m[2] = []string{"two", "deux"}
				m[3] = []string{"three", "trois"}
				m[4] = []string{"four", "quatre"}
				arr := stream.OfParallel(2, 1, 2, 3, 4).FlatMapOrdered(func(item interface{}) []interface{} {
					var result []interface{}
					for _, strings := range m[item.(int)] {
						result = append(result, strings)
					}
					return result
				}).ToTypedArray(reflect.TypeOf("")).Interface().([]string)
				gomega.Expect(arr).To(gomega.Equal([]string{"one", "un", "two", "deux", "three", "trois", "four", "quatre"}))
			})
		})
		ginkgo.When("Executing ForEach", func() {
			ginkgo.It("should do some side effects", func() {
				var arr []interface{}
				stream.OfParallel(2, 1, 2, 3, 4).ForEach(func(item interface{}) {
					arr = append(arr, item.(int))
				})
				gomega.Expect(arr).To(BagEquals([]interface{}{1, 2, 3, 4}))
			})
		})
		ginkgo.When("Executing Limit", func() {
			ginkgo.It("should retain a single value", func() {
				arr := stream.OfParallel(2, 1, 2, 3, 4).Limit(1).ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{1}))
			})
			ginkgo.It("should retain all values", func() {
				arr := stream.OfParallel(2, 1, 2, 3, 4).Limit(10).ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{1, 2, 3, 4}))
			})
		})
		ginkgo.When("Executing Map", func() {
			ginkgo.It("should return corresponding English words", func() {
				m := make(map[int]string)
				m[1] = "one"
				m[2] = "two"
				m[3] = "three"
				m[4] = "four"
				arr := stream.OfParallel(2, 1, 2, 3, 4).Map(func(item interface{}) interface{} {
					return m[item.(int)]
				}).ToArray()
				gomega.Expect(arr).To(BagEquals([]interface{}{"one", "two", "three", "four"}))
			})
		})
		ginkgo.When("Executing MapOrdered", func() {
			ginkgo.It("should return corresponding English words", func() {
				m := make(map[int]string)
				m[1] = "one"
				m[2] = "two"
				m[3] = "three"
				m[4] = "four"
				arr := stream.OfParallel(2, 1, 2, 3, 4).MapOrdered(func(item interface{}) interface{} {
					return m[item.(int)]
				}).ToTypedArray(reflect.TypeOf("")).Interface().([]string)
				gomega.Expect(arr).To(gomega.Equal([]string{"one", "two", "three", "four"}))
			})
		})
		ginkgo.When("Executing Max", func() {
			ginkgo.It("should return the maximum value", func() {
				max := stream.OfParallel(2, 2, 4, 3, 1).Max(func(a, b interface{}) bool {
					return a.(int) < b.(int)
				})
				gomega.Expect(max).To(gomega.Equal(util.Of(4)))
			})
			ginkgo.It("should return empty", func() {
				max := stream.OfParallel(1).Max(func(a, b interface{}) bool {
					return a.(int) < b.(int)
				})
				gomega.Expect(max.IsPresent()).To(gomega.BeFalse())
			})
		})
		ginkgo.When("Executing Min", func() {
			ginkgo.It("should return the minimal value", func() {
				min := stream.OfParallel(2, 3, 4, 1, 2).Min(func(a, b interface{}) bool {
					return a.(int) < b.(int)
				})
				gomega.Expect(min).To(gomega.Equal(util.Of(1)))
			})
			ginkgo.It("should return empty", func() {
				min := stream.OfParallel(1).Min(func(a, b interface{}) bool {
					return a.(int) < b.(int)
				})
				gomega.Expect(min.IsPresent()).To(gomega.BeFalse())
			})
		})
		ginkgo.When("Executing NoneMatch", func() {
			ginkgo.It("should return true", func() {
				result := stream.OfParallel(2, 1, 2, 3, 4).NoneMatch(func(item interface{}) bool {
					return item.(int) > 4
				})
				gomega.Expect(result).To(gomega.BeTrue())
			})
			ginkgo.It("should return false", func() {
				result := stream.OfParallel(2, 1, 2, 3, 4).NoneMatch(func(item interface{}) bool {
					return item.(int) > 1
				})
				gomega.Expect(result).To(gomega.BeFalse())
			})
		})
		ginkgo.When("Executing Peek", func() {
			ginkgo.It("should add one to every item", func() {
				arr := stream.OfParallel(2, 1, 2, 3, 4).Map(func(item interface{}) interface{} {
					return item.(int) + 1
				}).ToArray()
				gomega.Expect(arr).To(BagEquals([]interface{}{2, 3, 4, 5}))
			})
		})
		ginkgo.When("Executing Reduce", func() {
			ginkgo.It("should reduce to 10", func() {
				result := stream.OfParallel(2, 1, 2, 3, 4).Reduce(0, func(acc, cur interface{}) interface{} {
					return acc.(int) + cur.(int)
				})
				gomega.Expect(result).To(gomega.Equal(10))
			})
			ginkgo.It("should reduce to initial value", func() {
				result := stream.OfParallel(1).Reduce(10, func(acc, cur interface{}) interface{} {
					return cur.(int)
				})
				gomega.Expect(result).To(gomega.Equal(10))
			})
		})
		ginkgo.When("Executing ReduceCombine", func() {
			ginkgo.It("should reduce to 10 of int64 type", func() {
				result := stream.OfParallel(2, 1, 2, 3, 4).ReduceCombine(int64(0), func(acc, cur interface{}) interface{} {
					return int64(cur.(int))
				}, func(acc, cur interface{}) interface{} {
					return acc.(int64) + cur.(int64)
				})
				gomega.Expect(result).To(gomega.Equal(int64(10)))
			})
			ginkgo.It("should reduce to initial value of int64 type", func() {
				result := stream.OfParallel(1).ReduceCombine(int64(0), func(acc, cur interface{}) interface{} {
					return int64(cur.(int))
				}, func(acc, cur interface{}) interface{} {
					return acc.(int64) + cur.(int64)
				})
				gomega.Expect(result).To(gomega.Equal(int64(0)))
			})
		})
		ginkgo.When("Executing ReduceOptional", func() {
			ginkgo.It("should reduce to 10", func() {
				result := stream.OfParallel(2, 1, 2, 3, 4).ReduceOptional(func(acc, cur interface{}) interface{} {
					return acc.(int) + cur.(int)
				})
				gomega.Expect(result).To(gomega.Equal(util.Of(10)))
			})
			ginkgo.It("should reduce to nothing", func() {
				result := stream.OfParallel(1).ReduceOptional(func(acc, cur interface{}) interface{} {
					return acc.(int) + cur.(int)
				})
				gomega.Expect(result.IsPresent()).To(gomega.BeFalse())
			})
		})
		ginkgo.When("Executing Reverse", func() {
			ginkgo.It("should reverse the original array", func() {
				arr := stream.OfParallel(2, 1, 2, 3, 4).Reverse().ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{4, 3, 2, 1}))
			})
		})
		ginkgo.When("Executing Skip", func() {
			ginkgo.It("should remain three values", func() {
				arr := stream.OfParallel(2, 1, 2, 3, 4).Skip(1).ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{2, 3, 4}))
			})
			ginkgo.It("should remain no values", func() {
				arr := stream.OfParallel(2, 1, 2, 3, 4).Skip(10).ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{}))
			})
		})
		ginkgo.When("Executing Sorted", func() {
			ginkgo.It("should sort values in ascending order", func() {
				arr := stream.OfParallel(2, 4, 2, 1, 3).Sorted(func(a, b interface{}) bool {
					return a.(int) < b.(int)
				}).ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{1, 2, 3, 4}))
			})
		})
	})

	ginkgo.Context("Combined operation test", func() {
		ginkgo.When("Executing combined operations", func() {
			ginkgo.It("should return correctly", func() {
				arr := stream.OfParallel(2, 1, 2, 3, 4).Filter(func(item interface{}) bool {
					return item.(int) > 2
				}).Map(func(item interface{}) interface{} {
					return item.(int) + 1
				}).Reverse().Sorted(func(a, b interface{}) bool {
					return a.(int) < b.(int)
				}).Skip(1).ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{5}))
			})
			ginkgo.It("should return correctly", func() {
				count := stream.OfParallel(2, 1, 2, 3, 4).FlatMap(func(item interface{}) []interface{} {
					return []interface{}{item.(int)}
				}).Limit(2).Count()
				gomega.Expect(count).To(gomega.Equal(2))
			})
		})
	})

	ginkgo.Context("Collector test", func() {
		ginkgo.When("Executing collecting operations", func() {
			ginkgo.It("should return an interface array", func() {
				arr := stream.OfParallel(2, 1, 2, 3, 4).ToArray()
				gomega.Expect(arr).To(gomega.Equal([]interface{}{1, 2, 3, 4}))
			})
			ginkgo.It("should return a map with interface key and map", func() {
				m := stream.OfParallel(2, 1, 2, 3, 4).ToMap(func(item interface{}) interface{} {
					return item
				}, func(item interface{}) interface{} {
					return item
				})
				gomega.Expect(m).To(gomega.Equal(map[interface{}]interface{}{1: 1, 2: 2, 3: 3, 4: 4}))
			})
			ginkgo.It("should return a typed array", func() {
				arr := stream.OfParallel(2, 1, 2, 3, 4).ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{1, 2, 3, 4}))
			})
			ginkgo.It("should return a typed map with interface key and map", func() {
				m := stream.OfParallel(2, 1, 2, 3, 4).ToTypedMap(reflect.TypeOf(map[int]int{}), func(item interface{}) interface{} {
					return item
				}, func(item interface{}) interface{} {
					return item
				}).Interface().(map[int]int)
				gomega.Expect(m).To(gomega.Equal(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}))
			})
		})
	})
})

type BagMatcher struct {
	Expected []interface{}
}

func BagEquals(expected []interface{}) types.GomegaMatcher {
	return &BagMatcher{Expected: expected}
}

func (matcher *BagMatcher) Match(actual interface{}) (success bool, err error) {
	return bagEquals(matcher.Expected, actual.([]interface{})), nil
}

func (matcher *BagMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("to contain same elements as %+v", actual)
}

func (matcher *BagMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("not to contain same elements as %+v", actual)
}

func bagEquals(a, b []interface{}) bool {
	l := list.Empty
	for _, i := range b {
		l = l.Add(i)
	}
	for _, i := range a {
		index := l.FindIndex(func(item interface{}) bool {
			return item == i
		})
		var err error
		if index >= 0 {
			l, err = l.Remove(uint(index))
			if err != nil {
				panic("Should succeed when removing")
			}
		} else {
			return false
		}
	}
	return true
}
