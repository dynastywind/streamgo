package stream

import (
	"reflect"

	"github.com/dynastywind/go-stream/util"
)

type Stream interface {
	// AsParallel returns a parallel stream
	//
	// @param	routines	Number of go routines to use
	// @return	A parallel stream with the same pipeline as original stream
	AsParallel(routines int) Stream

	// AsSequence returns a sequential stream
	//
	// @return	A sequential stream with the same pipeline as original stream
	AsSequence() Stream

	// AllMatch returns true if all items in the data stream match a prediction
	//
	// @param predict	Prediction function
	// @return			True if all data items matches condition, false otherwise
	AllMatch(predict func(interface{}) bool) bool

	// AnyMatch returns true if any items in the data stream match a prediction
	//
	// @param predict	Prediction function
	// @return			True if any data items matches condition, false otherwise
	AnyMatch(predict func(interface{}) bool) bool

	// Count returns total number of items in data stream
	//
	// @return	Number of items in data stream
	Count() int

	// Distinct returns a data stream containing unique items
	//
	// @param hash	Function to generate data item's identity
	// @return		A data stream with unique items
	Distinct(hash func(interface{}) string) Stream

	// Filter returns a new stream containing only items matching filter condition
	// This method does not guarantee the processing order
	//
	// @param filter	Function to detect if a data item meets the condition
	// @return			A stream containing filtered data items
	Filter(filter func(interface{}) bool) Stream

	// FilterOrdered does the same thing as Filter but keeps the original order
	//
	// @param filter	Function to detect if a data item meets the condition
	// @return			A stream containing filtered data items
	FilterOrdered(filter func(interface{}) bool) Stream

	// FindAny randomly returns an item in the data stream if exists
	//
	// @return	Any item in data stream
	FindAny() *util.Optional

	// FindFirst find the first item in data stream if exists
	//
	// @return	First item in data stream
	FindFirst() *util.Optional

	// FlatMap applies a mapping function, which will generate a list of new items, onto every item in data stream, and then flatten the results into one list
	// This method does not guarantee the order of original items in the data stream
	//
	// @param mapper	Function to map an item into a list of new items of any type
	// @return			A stream after applyting flatmap operation
	FlatMap(mapper func(interface{}) []interface{}) Stream

	// FlatMapOrdered does the same thing as FlatMap, besides that it keeps the order of original items in the data stream
	//
	// @param mapper	Function to map an item into a list of new items of any type
	// @return			A stream after applyting flatmap operation
	FlatMapOrdered(mapper func(interface{}) []interface{}) Stream

	// ForEach applies a consumer function onto each item in data stream
	//
	// @param	Function to be applied onto data items
	ForEach(consumer func(interface{}))

	// IsParallel returns true if this stream is a parallel one
	//
	// @return	True if this stream is a parallel one, false otherwise
	IsParallel() bool

	// Limit returns a stream with only limited number of data items
	//
	// @param limit	Number of items to keep
	// @return		A stream after applying limit operation
	Limit(limit int) Stream

	// Map applies a function onto every item in data stream and returns another stream of data items with, maybe, different type
	// This method does not guarantee the processing order
	//
	// @param mapper	Function to be transform a data type into another one
	// @return			A stream after applying map operation
	Map(mapper func(interface{}) interface{}) Stream

	// MapOrdered does the same thing as Map function but keeps the order of the original data items
	//
	// @param mapper	Function to be transform a data type into another one
	// @return			A stream after applying map operation
	MapOrdered(mapper func(interface{}) interface{}) Stream

	// Max returns the maximum value in the data stream
	//
	// @param less	Function to judge which value is smaller
	// @return		Maximum value in the data stream
	Max(less func(interface{}, interface{}) bool) *util.Optional

	// Min returns the minimal value in the data stream
	//
	// @param less	Function to judge which value is smaller
	// @return		Minimal value in the data stream
	Min(less func(interface{}, interface{}) bool) *util.Optional

	// NoneMatch returns true if none of the items in the data stream matches the given condition
	//
	// @param predict	Prediction function
	// @return			True if none of items matches the given condition, false otherwise
	NoneMatch(predict func(interface{}) bool) bool

	// Peek applies a function onto each items in data stream and returns a new stream
	//
	// @param peeker	Function to be applied onto each data item
	// @return			A stream with items in original data stream
	Peek(peeker func(interface{})) Stream

	// Reduce returns a single value after accumulatively merge every data item in the stream
	//
	// @param init		Initial value to be accumulated
	// @param reducer	Function to merge elements
	// @return			A merged result
	Reduce(init interface{}, reducer func(interface{}, interface{}) interface{}) interface{}

	// ReduceCombine returns a single value after accumulatively merge and combine every data item in the stream
	//
	// @param init		Initial value to be accumulated
	// @param reducer	Function to merge elements
	// @param combiner	Function to combine merged result and other value
	// @return			A merged result
	ReduceCombine(init interface{}, reducer func(interface{}, interface{}) interface{}, combiner func(interface{}, interface{}) interface{}) interface{}

	// ReduceOptional returns a single value after accumulatively merge every data item in the stream if any item exists, empty otherwise
	//
	// @param reducer	Function to merge elements
	// @return			A merged result or empty
	ReduceOptional(reducer func(interface{}, interface{}) interface{}) *util.Optional

	// Reverse completely reverse the current order of data items in this stream and return a new stream
	//
	// @return	A stream with items' order reversed in original stream
	Reverse() Stream

	// Skip throws the first N items away in the data stream and returns the new stream
	// If N is greater than current stream length, an empty stream is returned
	//
	// @param skip	Number of items to be skipped
	// @return		A stream with first several elements given up
	Skip(skip int) Stream

	// Sorted sorts the data stream in ascending order
	//
	// @param less	Function to judge which value is smaller
	// @return		A stream with data items sorted in ascending order
	Sorted(less func(interface{}, interface{}) bool) Stream

	// ToArray collects data from this stream into an array
	//
	// @return	An array whose data is generated from this stream
	ToArray() []interface{}

	// ToMap collects data from this stream and transform to a map
	//
	// @param keyMapper		Function to map data item to map key
	// @param valueMapper	Function to map data item to map value
	// @return				A map whose data is generated from this stream
	ToMap(keyMapper func(interface{}) interface{}, valueMapper func(interface{}) interface{}) map[interface{}]interface{}

	// ToTypedArray does the same thing as ToArray method but will transform the result into a typed one via reflection
	//
	// @param t	Type of array element
	// @return	Typed array containing stream processing result
	ToTypedArray(t reflect.Type) reflect.Value

	// ToTypedMap does the same thing as ToMap method but will transform the result into a typed key-value pair via reflection
	//
	// @param t 			Type of target map
	// @param keyMapper		Function to map data item to map key
	// @param valueMapper	Function to map data item to map value
	// @return				Typed map containing stream processing result
	ToTypedMap(t reflect.Type, keyMapper func(interface{}) interface{}, valueMapper func(interface{}) interface{}) reflect.Value
}
