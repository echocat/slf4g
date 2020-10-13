// Package sdk/bridge provides methods to either hook into the SDK logger itself
// or create compatible instances.
//
// Hooks
//
// The simples way is to simply anonymously import the hook package to configure
// the whole application to use the slf4g framework on any usage of the SDK
// based loggers.
//
//    import (
//       _ "github.com/echocat/slf4g/sdk/bridge/hook"
//    )
//
// For manual hooks please see the examples.
package sdk
