# STREAM-GO

![go-stream workflow](https://github.com/dynastywind/streamgo/actions/workflows/go.yml/badge.svg)

This is a GOLANG port of Java's Stream API. It replicates most of the major methods supported in *java.util.stream.Stream* class.

Although it has achieved great success in parallel computing, it is notorious in GOLANG that there is little community support to provide production-ready, easy-to-understand and state-of-art utility libraries (e.g. advanced data structures, stream APIs to process data fluently), when comparing to other traditional languages like Java and C++. This repository is a small effort to make gophers' life better.

# Components

Go-Stream is a lazy-evaluated stream API which only invoke real computation when you ask for the result. That means, you can build up your processing pipeline without computing them immediately. This brings in lots of potential benefits like converting a sequential stream to a parallel one. File structure is as below:

- stream: core functions

    - operation: common operation functions

- util:   utility functions

## Sequential Stream

A sequential stream takes in a sequence of data and process it one by one, in FIFO order. To initalize a sequential stream, simply type code like below:

```go
import "github.com/dynastywind/go-stream/stream"

s := stream.Of(1, 2)

```
Then you can start to build up your own pipeline to consume your data. A possible pipeline would be like:

```go
s.Filter(func(item interface{}) bool {
    return item > 1
}).Reverse()
```

Since our stream is lazily evaluated, computation will only take place when you trigger a completion function like:

```go
s.ToArray()
```

Because Go does not support generic types currently (up to *Go 1.16*), the way to generate a typed array is a little bit ugly (and may have performance issues since reflection is adopted here). To have a typed array after processing your data, use **ToTypedArray** function and restore type information like this:

```go
// Take int as an example. For other types, simply adjust by yourself.
s.ToTypedArray(reflect.TypeOf(1)).Interface().([]int)
```

## Parallel Stream

A parallel stream takes in a sequence of data and triggers several *go routines* to process it parallelly. Thanks to Go's great support for parallel computing, this is not a hard one to implement (also not as easy as I initially thought).

> **Warning:** Unless using ordered version functions (with performance decay, of course), data order is not guaranteed in parallel stream. So only choose those functions when you want to speed up your computation and order is not cared of (or you will sort it).

To initialize a parallel stream, do this:

```go
s := stream.OfParallel(2, "a", "b")
```

Or you can convert a sequential stream to its parallel counterpart.

```go
// Specifies how many go routines you want to use
s = s.AsParallel(2)
```

And convert a parallel stream to its sequential brother like:

```go
s = s.AsSequence()
```

# Others

This repository will be updated further once *Go 1.17* is officially out. Any thoughts that will make this tool better are welcomed.

## About Sorting

Right now different sorting algorithms are adopted in different streams' *Sorted* method.

Sequential stream sort: [Heap sort](https://en.wikipedia.org/wiki/Heapsort)

Parallel stream sort: [Merge sort](https://en.wikipedia.org/wiki/Merge_sort)

Support to specify a customized sorting algorithm to process your data is a good possible improvement.

# License

MIT
