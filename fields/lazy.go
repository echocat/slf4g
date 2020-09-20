package fields

import "fmt"

type Lazy interface {
	Get() interface{}
}

type LazyFunc func() interface{}

func (instance LazyFunc) Get() interface{} {
	return instance()
}

func Format(format string, args ...interface{}) Lazy {
	return &lazyFormat{format, args}
}

type lazyFormat struct {
	format string
	args   []interface{}
}

func (instance *lazyFormat) Get() interface{} {
	return fmt.Sprintf(instance.format, instance.args...)
}

func (instance *lazyFormat) String() string {
	return instance.Get().(string)
}
