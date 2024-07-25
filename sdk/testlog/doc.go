// Package testlog provides a Logger which will be connected to testing.T.Log()
// of the go SDK.
//
// If you're looking for an instance to record all logged events see
// github.com/echocat/slf4g/testing/recording package.
//
// # Usage
//
// The easiest way to enable the slf4g framework in your tests is, simply:
//
//	import (
//		"testing"
//		"github.com/echocat/slf4g"
//		"github.com/echocat/slf4g/sdk/testlog"
//	)
//
//	func TestMyGreatStuff(t *testing.T) {
//		testlog.Hook(t)
//
//		log.Info("Yeah! This is a log!")
//	}
//
// ... that's it!
//
// See Hook(..) for more details.
package testlog
