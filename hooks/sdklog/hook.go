// Package hook_sdklog is an automatic hook for usage together with Go's SDK
// logging ([log.Logger]).
//
// Importing this package anonymously will configure the whole application to
// use the slf4g framework on any usage of the SDK based loggers.
//
//	import (
//	   _ "github.com/echocat/slf4g/hooks/sdklog"
//	)
package hook_sdklog

import std "github.com/echocat/slf4g/sdk/bridge"

func init() {
	std.Configure()
}
