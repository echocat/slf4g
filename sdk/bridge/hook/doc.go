// Package hook is an automatic hook for usage together with Go's SDK logging
// ([log.Logger]).
//
// Importing this package anonymously will configure the whole application to
// use the slf4g framework on any usage of the SDK based loggers.
//
//	import (
//	   _ "github.com/echocat/slf4g/sdk/bridge/hook"
//	)
package hook
