package stream

type OperationTag string

const (
	DISTINCT         OperationTag = "DISTINCT"
	FILTER           OperationTag = "FILTER"
	FILTER_ORDERED   OperationTag = "FILTER_ORDERED"
	FLAT_MAP         OperationTag = "FLATMAP"
	FLAT_MAP_ORDERED OperationTag = "FLAT_MAP_ORDERED"
	LIMIT            OperationTag = "LIMIT"
	MAP              OperationTag = "MAP"
	MAP_ORDERED      OperationTag = "MAP_ORDERED"
	PEEK             OperationTag = "PEEK"
	REVERSE          OperationTag = "REVERSE"
	SKIP             OperationTag = "SKIP"
	SORTED           OperationTag = "SORTED"
)

type OperationDescriptor struct {
	tag    OperationTag
	params []interface{}
}
