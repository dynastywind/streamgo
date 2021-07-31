package stream_test

import (
	"reflect"
	"strconv"

	"github.com/dynastywind/go-stream/stream"
	"github.com/dynastywind/go-stream/util"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Test if every function works well", func() {
	ginkgo.Context("Stream construction test", func() {
		ginkgo.When("Construct stream with items", func() {
			ginkgo.It("should a construct a sequential stream", func() {
				arr := stream.Of(1, 2, 3, 4).ToArray()
				gomega.Expect(arr).To(gomega.Equal([]interface{}{1, 2, 3, 4}))
			})
		})
		ginkgo.When("Construct a stream with array", func() {
			ginkgo.It("should construct a sequential stream from an interface array", func() {
				arr := stream.FromArray([]interface{}{1, 2, 3, 4}).ToArray()
				gomega.Expect(arr).To(gomega.Equal([]interface{}{1, 2, 3, 4}))
			})
		})
		ginkgo.When("Construct a stream with typed array", func() {
			ginkgo.It("should construct a sequential stream from a typed array", func() {
				arr := stream.FromTypedArray([]int{1, 2, 3, 4}).ToArray()
				gomega.Expect(arr).To(gomega.Equal([]interface{}{1, 2, 3, 4}))
			})
		})
		ginkgo.When("Construct a stream with concatenation", func() {
			ginkgo.It("should construct a sequential stream from several streams", func() {
				s := stream.Concat(stream.Of(1, 2, 3, 4), stream.OfParallel(2, 5, 6, 7, 8))
				gomega.Expect(s.ToArray()).To(gomega.Equal([]interface{}{1, 2, 3, 4, 5, 6, 7, 8}))
				gomega.Expect(s.IsParallel()).To(gomega.BeFalse())
			})
		})
	})
	ginkgo.Context("Stream tranformation test", func() {
		ginkgo.When("Executing AsParallel", func() {
			ginkgo.It("should return a parallel stream", func() {
				s := stream.Of(1, 2, 3, 4).AsParallel(2)
				gomega.Expect(s.IsParallel()).To(gomega.BeTrue())
			})
		})
		ginkgo.When("Executing AsSequence", func() {
			ginkgo.It("should return itself", func() {
				original := stream.Of(1, 2, 3, 4)
				s := original.AsSequence()
				gomega.Expect(s.IsParallel()).To(gomega.BeFalse())
				gomega.Expect(s).To(gomega.Equal(original))
			})
		})
	})
	ginkgo.Context("Sigle operation test", func() {
		ginkgo.When("Executing AllMatch", func() {
			ginkgo.It("should return true", func() {
				result := stream.Of(1, 2, 3, 4).AllMatch(func(item interface{}) bool {
					return item.(int) > 0
				})
				gomega.Expect(result).To(gomega.BeTrue())
			})
			ginkgo.It("should return false", func() {
				result := stream.Of(1, 2, 3, 4).AllMatch(func(item interface{}) bool {
					return item.(int) > 1
				})
				gomega.Expect(result).To(gomega.BeFalse())
			})
		})
		ginkgo.When("Executing AnyMatch", func() {
			ginkgo.It("should return true", func() {
				result := stream.Of(1, 2, 3, 4).AnyMatch(func(item interface{}) bool {
					return item.(int) > 1
				})
				gomega.Expect(result).To(gomega.BeTrue())
			})
			ginkgo.It("should return false", func() {
				result := stream.Of(1, 2, 3, 4).AnyMatch(func(item interface{}) bool {
					return item.(int) > 4
				})
				gomega.Expect(result).To(gomega.BeFalse())
			})
		})
		ginkgo.When("Executing Count", func() {
			ginkgo.It("should return the length of the data stream", func() {
				count := stream.Of(1, 2, 3, 4).Count()
				gomega.Expect(count).To(gomega.Equal(4))
			})
		})
		ginkgo.When("Executing Distinct", func() {
			ginkgo.It("should return an array of unique values", func() {
				arr := stream.Of(1, 1, 2, 2).Distinct(func(item interface{}) string {
					return strconv.Itoa(item.(int))
				}).ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{1, 2}))
			})
		})
		ginkgo.When("Executing Filter", func() {
			ginkgo.It("should return an array of values greater than 2", func() {
				arr := stream.Of(1, 2, 3, 4).Filter(func(item interface{}) bool {
					return item.(int) > 2
				}).ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{3, 4}))
			})
		})
		ginkgo.When("Executing FilterOrdered", func() {
			ginkgo.It("should return an array of values greater than 2", func() {
				arr := stream.Of(1, 2, 3, 4).FilterOrdered(func(item interface{}) bool {
					return item.(int) > 2
				}).ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{3, 4}))
			})
		})
		ginkgo.When("Executing FindAny", func() {
			ginkgo.It("should find a value", func() {
				result := stream.Of(1, 2, 3, 4).FindAny()
				gomega.Expect(result.IsPresent()).To(gomega.BeTrue())
			})
			ginkgo.It("should find nothing", func() {
				result := stream.Of().FindAny()
				gomega.Expect(result.IsPresent()).To(gomega.BeFalse())
			})
		})
		ginkgo.When("Executing FindFirst", func() {
			ginkgo.It("should find the first value", func() {
				result := stream.Of(1, 2, 3, 4).FindFirst()
				gomega.Expect(result).To(gomega.Equal(util.Of(1)))
			})
			ginkgo.It("should find nothing", func() {
				result := stream.Of().FindFirst()
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
				arr := stream.Of(1, 2, 3, 4).FlatMap(func(item interface{}) []interface{} {
					var result []interface{}
					for _, strings := range m[item.(int)] {
						result = append(result, strings)
					}
					return result
				}).ToTypedArray(reflect.TypeOf("")).Interface().([]string)
				gomega.Expect(arr).To(gomega.Equal([]string{"one", "un", "two", "deux", "three", "trois", "four", "quatre"}))
			})
		})
		ginkgo.When("Executing FlatMapOrdered", func() {
			ginkgo.It("should return corresponding words in a row", func() {
				m := make(map[int][]string)
				m[1] = []string{"one", "un"}
				m[2] = []string{"two", "deux"}
				m[3] = []string{"three", "trois"}
				m[4] = []string{"four", "quatre"}
				arr := stream.Of(1, 2, 3, 4).FlatMapOrdered(func(item interface{}) []interface{} {
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
				var arr []int
				stream.Of(1, 2, 3, 4).ForEach(func(item interface{}) {
					arr = append(arr, item.(int))
				})
				gomega.Expect(arr).To(gomega.Equal([]int{1, 2, 3, 4}))
			})
		})
		ginkgo.When("Executing Limit", func() {
			ginkgo.It("should retain a single value", func() {
				arr := stream.Of(1, 2, 3, 4).Limit(1).ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{1}))
			})
			ginkgo.It("should retain all values", func() {
				arr := stream.Of(1, 2, 3, 4).Limit(10).ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
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
				arr := stream.Of(1, 2, 3, 4).Map(func(item interface{}) interface{} {
					return m[item.(int)]
				}).ToTypedArray(reflect.TypeOf("")).Interface().([]string)
				gomega.Expect(arr).To(gomega.Equal([]string{"one", "two", "three", "four"}))
			})
		})
		ginkgo.When("Executing MapOrdered", func() {
			ginkgo.It("should return corresponding English words", func() {
				m := make(map[int]string)
				m[1] = "one"
				m[2] = "two"
				m[3] = "three"
				m[4] = "four"
				arr := stream.Of(1, 2, 3, 4).MapOrdered(func(item interface{}) interface{} {
					return m[item.(int)]
				}).ToTypedArray(reflect.TypeOf("")).Interface().([]string)
				gomega.Expect(arr).To(gomega.Equal([]string{"one", "two", "three", "four"}))
			})
		})
		ginkgo.When("Executing Max", func() {
			ginkgo.It("should return the maximum value", func() {
				max := stream.Of(2, 4, 3, 1).Max(func(a, b interface{}) bool {
					return a.(int) < b.(int)
				})
				gomega.Expect(max).To(gomega.Equal(util.Of(4)))
			})
			ginkgo.It("should return empty", func() {
				max := stream.Of().Max(func(a, b interface{}) bool {
					return a.(int) < b.(int)
				})
				gomega.Expect(max.IsPresent()).To(gomega.BeFalse())
			})
		})
		ginkgo.When("Executing Min", func() {
			ginkgo.It("should return the minimal value", func() {
				min := stream.Of(3, 4, 1, 2).Min(func(a, b interface{}) bool {
					return a.(int) < b.(int)
				})
				gomega.Expect(min).To(gomega.Equal(util.Of(1)))
			})
			ginkgo.It("should return empty", func() {
				min := stream.Of().Min(func(a, b interface{}) bool {
					return a.(int) < b.(int)
				})
				gomega.Expect(min.IsPresent()).To(gomega.BeFalse())
			})
		})
		ginkgo.When("Executing NoneMatch", func() {
			ginkgo.It("should return true", func() {
				result := stream.Of(1, 2, 3, 4).NoneMatch(func(item interface{}) bool {
					return item.(int) > 4
				})
				gomega.Expect(result).To(gomega.BeTrue())
			})
			ginkgo.It("should return false", func() {
				result := stream.Of(1, 2, 3, 4).NoneMatch(func(item interface{}) bool {
					return item.(int) > 1
				})
				gomega.Expect(result).To(gomega.BeFalse())
			})
		})
		ginkgo.When("Executing Peek", func() {
			ginkgo.It("should add one to every item", func() {
				arr := stream.Of(1, 2, 3, 4).Map(func(item interface{}) interface{} {
					return item.(int) + 1
				}).ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{2, 3, 4, 5}))
			})
		})
		ginkgo.When("Executing Reduce", func() {
			ginkgo.It("should reduce to 10", func() {
				result := stream.Of(1, 2, 3, 4).Reduce(0, func(acc, cur interface{}) interface{} {
					return acc.(int) + cur.(int)
				})
				gomega.Expect(result).To(gomega.Equal(10))
			})
			ginkgo.It("should reduce to initial value", func() {
				result := stream.Of().Reduce(10, func(acc, cur interface{}) interface{} {
					return cur.(int)
				})
				gomega.Expect(result).To(gomega.Equal(10))
			})
		})
		ginkgo.When("Executing ReduceCombine", func() {
			ginkgo.It("should reduce to 10 of int64 type", func() {
				result := stream.Of(1, 2, 3, 4).ReduceCombine(int64(0), func(acc, cur interface{}) interface{} {
					return int64(cur.(int))
				}, func(acc, cur interface{}) interface{} {
					return acc.(int64) + cur.(int64)
				})
				gomega.Expect(result).To(gomega.Equal(int64(10)))
			})
			ginkgo.It("should reduce to initial value of int64 type", func() {
				result := stream.Of().ReduceCombine(int64(0), func(acc, cur interface{}) interface{} {
					return int64(cur.(int))
				}, func(acc, cur interface{}) interface{} {
					return acc.(int64) + cur.(int64)
				})
				gomega.Expect(result).To(gomega.Equal(int64(0)))
			})
		})
		ginkgo.When("Executing ReduceOptional", func() {
			ginkgo.It("should reduce to 10", func() {
				result := stream.Of(1, 2, 3, 4).ReduceOptional(func(acc, cur interface{}) interface{} {
					return acc.(int) + cur.(int)
				})
				gomega.Expect(result).To(gomega.Equal(util.Of(10)))
			})
			ginkgo.It("should reduce to nothing", func() {
				result := stream.Of().ReduceOptional(func(acc, cur interface{}) interface{} {
					return acc.(int) + cur.(int)
				})
				gomega.Expect(result.IsPresent()).To(gomega.BeFalse())
			})
		})
		ginkgo.When("Executing Reverse", func() {
			ginkgo.It("should reverse the original array", func() {
				arr := stream.Of(1, 2, 3, 4).Reverse().ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{4, 3, 2, 1}))
			})
		})
		ginkgo.When("Executing Skip", func() {
			ginkgo.It("should remain three values", func() {
				arr := stream.Of(1, 2, 3, 4).Skip(1).ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{2, 3, 4}))
			})
			ginkgo.It("should remain no values", func() {
				arr := stream.Of(1, 2, 3, 4).Skip(10).ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{}))
			})
		})
		ginkgo.When("Executing Sorted", func() {
			ginkgo.It("should sort values in ascending order", func() {
				arr := stream.Of(4, 2, 1, 3).Sorted(func(a, b interface{}) bool {
					return a.(int) < b.(int)
				}).ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{1, 2, 3, 4}))
			})
		})
	})

	ginkgo.Context("Combined operation test", func() {
		ginkgo.When("Executing combined operations", func() {
			ginkgo.It("should return correctly", func() {
				arr := stream.Of(1, 2, 3, 4).Filter(func(item interface{}) bool {
					return item.(int) > 2
				}).Map(func(item interface{}) interface{} {
					return item.(int) + 1
				}).Reverse().Sorted(func(a, b interface{}) bool {
					return a.(int) < b.(int)
				}).Skip(1).ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{5}))
			})
			ginkgo.It("should return correctly", func() {
				count := stream.Of(1, 2, 3, 4).FlatMap(func(item interface{}) []interface{} {
					return []interface{}{item.(int)}
				}).Limit(2).Count()
				gomega.Expect(count).To(gomega.Equal(2))
			})
		})
	})

	ginkgo.Context("Collector test", func() {
		ginkgo.When("Executing collecting operations", func() {
			ginkgo.It("should return an interface array", func() {
				arr := stream.Of(1, 2, 3, 4).ToArray()
				gomega.Expect(arr).To(gomega.Equal([]interface{}{1, 2, 3, 4}))
			})
			ginkgo.It("should return a map with interface key and map", func() {
				m := stream.Of(1, 2, 3, 4).ToMap(func(item interface{}) interface{} {
					return item
				}, func(item interface{}) interface{} {
					return item
				})
				gomega.Expect(m).To(gomega.Equal(map[interface{}]interface{}{1: 1, 2: 2, 3: 3, 4: 4}))
			})
			ginkgo.It("should return a typed array", func() {
				arr := stream.Of(1, 2, 3, 4).ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
				gomega.Expect(arr).To(gomega.Equal([]int{1, 2, 3, 4}))
			})
			ginkgo.It("should return a typed map with interface key and map", func() {
				m := stream.Of(1, 2, 3, 4).ToTypedMap(reflect.TypeOf(map[int]int{}), func(item interface{}) interface{} {
					return item
				}, func(item interface{}) interface{} {
					return item
				}).Interface().(map[int]int)
				gomega.Expect(m).To(gomega.Equal(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}))
			})
		})
	})
})
