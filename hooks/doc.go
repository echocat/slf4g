// Package hooks contains a collection of several automatic hooks, like:
//
//  1. [github.com/echocat/slf4g/hooks/sdklog] which configures the SDK's classic
//     [log] package to use slf4g as its logging backend.
//  2. [github.com/echocat/slf4g/hooks/sdkslog] which configures the SDK's modern
//     [log/slog] package to use slf4g as its logging backend.
//
// Simply import those packages as follows:
//
//	import (
//		// For hook into SDK's log package
//		_ "github.com/echocat/slf4g/hooks/sdklog"
//
//		// For hook into SDK's log/slog package
//		_ "github.com/echocat/slf4g/hooks/sdkslog"
//	)
package hooks
