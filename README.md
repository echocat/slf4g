[![PkgGoDev](https://pkg.go.dev/badge/github.com/echocat/slf4g)](https://pkg.go.dev/github.com/echocat/slf4g)
[![Continuous Integration](https://github.com/echocat/slf4g/workflows/Continuous%20Integration/badge.svg)](https://github.com/echocat/slf4g/actions?query=workflow%3A%22Continuous+Integration%22)
[![Coverage Status](https://coveralls.io/repos/github/echocat/slf4g/badge.svg?branch=main)](https://coveralls.io/github/echocat/slf4g?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/echocat/slf4g)](https://goreportcard.com/report/github.com/echocat/slf4g)

# Simple Log Framework for Golang (slf4g)

## TOC

* [Principles](#principles)
* [Motivation](#motivation)
* [Getting started](#getting-started)
* [Implementations](#implementations)
* [Bridges](#bridges) (and their [hooks](#hooks))
* [Contributing](#contributing)
* [License](#license)

## Principles

### KISS for users

If you want to log something it should be easy, intuitive and straight forward. This should work either for very small applications (just one file) or big applications, too.

You should not care about the implementation, you want just to use it.

### KISS for implementators

If you want to implement a new logger you should not be required to educate your users how to use it.

You just want to write a logger for a user-case. Someone else should take care of the public API.

### Separation

A logging framework should not be both, by design: API and implementation.

You want to have one API and the possibility to use whatever implementation you want.

### Interoperable

Regardless how used or implemented you want that just everything sticks together, and you're not ending up over and over again in writing new wrappers or in worst-case see different styled log messages in the console.

Every library should just work transparently with one logger.

## Motivation

I've tried out many logging frameworks for Golang. They're coming with different promises, like:

1. There are ones which tries to be "blazing fast"; focussing on be fast and non-blocking to be able to log as much log events as possible.

2. Other ones are trying to be as minimalistic as possible. Just using a very few amount of code to work.

3. ...

...but overall they're just violating the [Principles listed above](#principles) and we're ending up in just mess.

Slf4g is born out of feeling this pain every day again and be simply annoyed. It is inspired by [Simple Logging Facade for Java (SLF4J)](http://www.slf4j.org/), which was born out of the same pains; but obviously in Java before. Since [SLF4J](http://www.slf4j.org/) exists, and it is now broadly used in Java, nobody does experience this issues any longer.

## Getting started

It is very easy to use slf4g (as the naming is promising ☺️):

1. Import the API to your current project (in best with a [Go Modules project](https://blog.golang.org/using-go-modules))
    ```bash
    $ go get -u github.com/echocat/slf4g
    ```

2. Select one of the implementation of [slf4g](https://github.com/echocat/slf4g) and import it too (see [Implementations](#implementations)). Example:

    ```bash
    $ go get -u github.com/echocat/slf4g/native
    ```

   > ℹ️ If you do not pick one implementation the fallback logger is used. It works, too, but is obviously less powerful and is not customizable. It is comparable with the [SDK based logger](https://pkg.go.dev/log).

3. Configure your application to use the selected logger implementation:

   This should be only done in `main/main.go`:

    ```go
    package main
   
    import (
    	_ "github.com/echocat/slf4g/native"
    )
   
    func main() {
    	// do your stuff...
    }
    ```

4. In each package create a logger variable, in best case you create a file named `common.go` or `package.go` which will contain it:

    ```go
    package foo
   
    import (
    	"github.com/echocat/slf4g"
    )
   
    var logger = log.GetLoggerForCurrentPackage()
    ```

5. Now you're ready to go. In every file of this package you can do stuff, like:

    ```go
    package foo
   
    func MyFunction() {
    	logger.Info("Hello, world!")

    	if !loaded {
    		logger.With("field", 123).
    		       Warn("That's not great.")
    	}

    	if err := doSomething(); err != nil {
    		logger.WithError(err).
    		       Error("Doh!")
    	}
    }
    ```

   For sure, you're able to simply do stuff like that (although to ensure interoperability this is not recommended):

    ```go
    package foo
    
    import (
    	"github.com/echocat/slf4g"
    )
    
    func MyFunction() {
    	log.Info("Hello, world!")
    
    	if !loaded {
    		log.With("field", 123).
    		    Warn("That's not great.")
    	}

    	if err := doSomething(); err != nil {
    		log.WithError(err).
    		    Error("Doh!")
    	}
    }
    ```

Done. Enjoy!

## Implementations

1. [native](native): Reference implementation of [slf4g](https://github.com/echocat/slf4g), best use for most applications.

2. [testlog](sdk/testlog): Ensure that everything which is logged within test by [slf4g](https://github.com/echocat/slf4g) appears correctly within tests.

3. [recording](testing/recording): Will record everything which is logged by [slf4g](https://github.com/echocat/slf4g) and can then be asserted inside test cases.

## Bridges

There are several bridges available to use [slf4g](https://github.com/echocat/slf4g) in other frameworks:

1. [sdk/bridge](sdk/bridge) to implement the [Go's SDK log interface](https://pkg.go.dev/log).
2. [sdk/bridge/slog](sdk/bridge/slog) to implement the [Go's SDK slog interface](https://pkg.go.dev/log/slog).
3. [github.com/echocat/slf4g-logr](https://github.com/echocat/slf4g-logr) to implement [github.com/go-logr/logr](https://github.com/go-logr/logr).
4. [github.com/echocat/slf4g-logrus](https://github.com/echocat/slf4g-logrus) to implement [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus).
5. [github.com/echocat/slf4g-klog](https://github.com/echocat/slf4g-klog) to implement [k8s.io/klog/v2](https://github.com/kubernetes/klog).

### Hooks

These are automatically registering itself by simply calling an anonymous import (instead of explicitly importing - see [Bridges](#bridges) above), like:

```go
package main

import (
	// For hook into SDK's log package
	_ "github.com/echocat/slf4g/hooks/sdklog"

	// For hook into SDK's log/slog package
	_ "github.com/echocat/slf4g/hooks/sdkslog"

	// For hook into github.com/sirupsen/logrus
	_ "github.com/echocat/slf4g-logrus/logrus2slf4g/hook"

	// For hook into Kubernetes' k8s.io/klog/v2
	_ "github.com/echocat/slf4g-klog/bridge/hook"
)
```

## Contributing

**slf4g** is an open source project by [echocat](https://echocat.org). So if you want to make this project even better, you can contribute to this project on [Github](https://github.com/echocat/slf4g) by [fork us](https://github.com/echocat/slf4g/fork).

If you commit code to this project, you have to accept that this code will be released under the [license](#license) of this project.

## License

See the [LICENSE](LICENSE) file.
