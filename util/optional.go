package util

import (
	"fmt"
)

type Optional struct {
	data interface{}
}

func Of(data interface{}) *Optional {
	if data == nil {
		panic("Data should not be nil")
	}
	return &Optional{
		data: data,
	}
}

func OfNillable(data interface{}) *Optional {
	if data == nil {
		return &Optional{}
	}
	return &Optional{
		data: data,
	}
}

func OfEmpty() *Optional {
	return &Optional{}
}

func (opt *Optional) IsPresent() bool {
	return opt.data != nil
}

func (opt *Optional) Get() interface{} {
	if !opt.IsPresent() {
		panic("No such element")
	}
	return opt.data
}

func (opt *Optional) OrElse(or interface{}) interface{} {
	if !opt.IsPresent() {
		return or
	}
	return opt.Get()
}

func (opt *Optional) OrElseGet(getter func() interface{}) interface{} {
	if !opt.IsPresent() {
		return getter()
	}
	return opt.Get()
}

func (opt *Optional) OrElseError(err func() error) (interface{}, error) {
	if !opt.IsPresent() {
		return nil, err()
	}
	return opt.data, nil
}

func (opt *Optional) Filter(filter func(data interface{}) bool) *Optional {
	if !opt.IsPresent() || filter(opt.Get()) {
		return opt
	}
	return OfEmpty()
}

func (opt *Optional) Map(mapper func(data interface{}) interface{}) *Optional {
	if !opt.IsPresent() {
		return opt
	}
	return OfNillable(mapper(opt.Get()))
}

func (opt *Optional) FlatMap(mapper func(data interface{}) *Optional) *Optional {
	if !opt.IsPresent() {
		return opt
	}
	return mapper(opt.Get())
}

func (opt *Optional) IfPresent(consumer func(data interface{})) {
	if opt.IsPresent() {
		consumer(opt.Get())
	}
}

func (opt *Optional) String() string {
	if !opt.IsPresent() {
		return ""
	}
	return fmt.Sprint(opt.data)
}
