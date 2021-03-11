package value_test

import (
	"flag"

	"github.com/echocat/slf4g/native"
	"github.com/echocat/slf4g/native/facade/value"
)

func ExampleNewProvider() {
	pv := value.NewProvider(native.DefaultProvider)

	flag.Var(pv.Consumer.Formatter, "log.format", "Configures the log format.")
	flag.Var(pv.Level, "log.level", "Configures the log level.")

	flag.Parse()

	// Now you can call:
	// $ <myExecutable> -log.format=json -log.level=debug ...
}
