package value

import "fmt"

type Lazy interface {
	Value
	Get() Value
}

func Format(format string, args ...interface{}) Lazy {
	return &lazyFormat{format, args}
}

type lazyFormat struct {
	format string
	args   []interface{}
}

func (instance *lazyFormat) Get() Value {
	return fmt.Sprintf(instance.format, instance.args...)
}
