// Package log is the Simple Logging Facade for Go provides an easy access who
// everyone who wants to log something and do not want to care how it is logged
// and gives others the possibility to implement their own loggers in easy way.
//
// # Usage
//
// There are 2 common ways to use this framework.
//
// 1. By getting a logger for your current package. This is the most common way
// and is quite clean for one package. If the package contains too many logic
// it might worth it to use the 2nd approach (see below).
//
//	package foo
//
//	import "github.com/echocat/slf4g"
//
//	var logger = log.GetLoggerForCurrentPackage()
//
//	func sayHello() {
//	   // will log with logger="github.com/user/foo"
//	   logger.Info("Hello, world!")
//	}
//
// 2. By getting a logger for the object the logger is for. This is the most
// clean approach and will give you later the maximum flexibility and control.
//
//	package foo
//
//	import "github.com/echocat/slf4g"
//
//	var logger = log.GetLogger(myType{})
//
//	type myType struct {
//	   ...
//	}
//
//	func (mt myType) sayHello() {
//	   // will log with logger="github.com/user/foo.myType"
//	   logger.Info("Hello, world!")
//	}
//
// 3. By using the global packages methods which is quite equal to how the SDK
// base logger works. This is only recommend for small application and not for
// libraries you like to export.
//
//	package foo
//
//	import "github.com/echocat/slf4g"
//
//	func sayHello() {
//	   // will log with logger=ROOT
//	   log.Info("Hello, world!")
//	}
package log
