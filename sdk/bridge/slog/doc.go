//go:build go1.21

// Package sdk/slog provides methods to either hook into the SDK slog logger itself
// or create compatible instances.
//
// # Hooks
//
// The simples way is to simply anonymously import the hook package to configure
// the whole application to use the slf4g framework on any usage of the SDK
// based slog loggers.
//
//	import (
//	   _ "github.com/echocat/slf4g/hooks/sdkslog"
//	)
package sdk
