package fields

import "fmt"

// Lazy is a value which CAN be initialized on usage.
//
// This is very useful in the context of Fields where sometimes the evaluating
// of values could be cost intensive but maybe you either might log stuff on a
// level which might not be always enabled or the operation might be happening
// on an extra routine/thread.
type Lazy interface {
	// Get is the method which will be called in the moment where the value
	// should be consumed.
	Get() interface{}
}

// LazyFunc wraps Lazy into a single function pointer.
type LazyFunc func() interface{}

func (instance LazyFunc) Get() interface{} {
	return instance()
}

// LazyFormat returns a value which will be execute the fmt.Sprintf action in
// the moment when it will be consumed or in other words: Lazy.Get() is called.
func LazyFormat(format string, args ...interface{}) Lazy {
	return &lazyFormat{format, args}
}

type lazyFormat struct {
	format string
	args   []interface{}
}

func (instance *lazyFormat) Get() interface{} {
	return instance.String()
}

func (instance *lazyFormat) String() string {
	return fmt.Sprintf(instance.format, instance.args...)
}
