//go:build go1.21

// Package hook_sdkslog is an automatic hook for usage together with Go's SDK
// [log/slog] logging ([log/slog.Logger]).
//
// Importing this package anonymously will configure the whole application to
// use the slf4g framework on any usage of the SDK based [log/slog] loggers.
//
//	import (
//	   _ "github.com/echocat/slf4g/hooks/sdkslog"
//	)
package hook_sdkslog

import std "github.com/echocat/slf4g/sdk/bridge/slog"

func init() {
	std.Configure()
}
