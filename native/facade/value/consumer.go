package value

import (
	"fmt"
	"os"
	"reflect"

	"github.com/echocat/slf4g/native/consumer"
	"github.com/echocat/slf4g/native/formatter"
)

// ConsumerTarget defines an object that receives the consumer.Consumer managed
// by the Consumer value facade.
type ConsumerTarget interface {
	consumer.MutableAware
}

// Consumer is a value facade for transparent setting of consumer.Consumer
// for the slf4g/native implementation. This is quite handy for usage
// with flags package of the SDK or similar flag libraries. This might
// be usable, too in contexts where serialization might be required.
type Consumer struct {
	// Level is the corresponding formatter.Formatter facade.
	Formatter Formatter
}

// NewProvider create a new instance of Provider with the given target ProviderTarget instance.
func NewConsumer(target ConsumerTarget, customizer ...func(*Consumer)) Consumer {
	if c := target.GetConsumer(); c == nil {
		if d := consumer.Default; d != nil {
			target.SetConsumer(d)
		} else {
			target.SetConsumer(consumer.NewWriter(os.Stderr))
		}
	}

	fa, ok := target.GetConsumer().(formatter.MutableAware)
	if !ok {
		panic(fmt.Errorf("%v does not implement %v", reflect.TypeOf(target.GetConsumer()), reflect.TypeOf((formatter.MutableAware)(nil))))
	}

	result := Consumer{
		Formatter: NewFormatter(fa),
	}

	for _, c := range customizer {
		c(&result)
	}

	return result
}
