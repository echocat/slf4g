// This is the reference implementation of a logger of the slf4g framework
// (https://github.com/echocat/slf4g).
//
// Usage
//
// For the most common cases it is fully enough to anonymously import this
// package in your main.go; nothing more is needed.
//
// github.com/foo/bar/main/main.go:
//    package main
//
//    import (
//       "github.com/foo/bar"
//       _ "github.com/echocat/slf4g/native"
//    )
//
//    func main() {
//        bar.SayHello()
//    }
//
// github.com/foo/bar/bar.go:
//    package bar
//
//    import (
//       "github.com/echocat/slf4g"
//    )
//
//    func SayHello() {
//        log.Info("Hello, world!")
//    }
//
// See more useful stuff in the examples sections.
package native
