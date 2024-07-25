// Package value provides a value facade for native.Provider to be able
// to easy be configured using flag libraries like the SDK implementation
// or other compatible ones. It implements encoding.TextMarshaler and
// encoding.TextUnmarshaler, too.
//
// Example:
//
//	pv := value.NewProvider(native.DefaultProvider)
//
//	flag.Var(pv.Consumer.Formatter, "log.format", "Configures the log format.")
//	flag.Var(pv.Level, "log.level", "Configures the log level.")
//
//	flag.Parse()
//
// Now you can call:
//
//	$ <myExecutable> -log.format=json -log.level=debug ...
package value
