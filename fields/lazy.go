package fields

import (
	"fmt"

	"github.com/echocat/slf4g/internal/support"
)

// Lazy is a value which CAN be initialized on usage.
//
// This is very useful in the context of Fields where sometimes the evaluating
// of values could be cost intensive, but maybe you either might log stuff on a
// level which might not be always enabled or the operation might be happening
// on an extra routine/thread.
type Lazy interface {
	// Get is the method which will be called at the moment where the value
	// should be consumed.
	Get() any
}

// LazyFunc wraps Lazy into a single function pointer.
func LazyFunc(provider func() any) Lazy {
	return lazyFunc(provider)
}

type lazyFunc func() any

func (instance lazyFunc) Get() any {
	if instance == nil {
		return nil
	}
	return instance()
}

// LazyFormat returns a value which will be executed the fmt.Sprintf action at
// the moment when it will be consumed or in other words: Lazy.Get() is called.
func LazyFormat(format string, args ...any) Lazy {
	return &lazyFormat{format, append([]any(nil), args...)}
}

type lazyFormat struct {
	format string
	args   []any
}

func (instance *lazyFormat) Get() any {
	return instance.String()
}

func (instance *lazyFormat) String() string {
	targetArgs := make([]any, len(instance.args))
	for i, arg := range instance.args {
		targetArgs[i] = resolveLazy(arg)
	}
	return fmt.Sprintf(instance.format, targetArgs...)
}

func resolveLazy(value any) any {
	if value == nil {
		return nil
	}
	if lazy, ok := value.(Lazy); ok {
		if support.IsNil(lazy) {
			return nil
		}
		return lazy.Get()
	}
	return value
}
