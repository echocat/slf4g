[![PkgGoDev](https://pkg.go.dev/badge/github.com/echocat/slf4g/native)](https://pkg.go.dev/github.com/echocat/slf4g/native)

# slf4g Native

This is the native/reference implementation of [Simple Log Framework for Golang (slf4g)](..).

## Usage

1. Import the required dependencies to your current project (in best with a [Go Modules project](https://blog.golang.org/using-go-modules))
    ```bash
    $ go get -u github.com/echocat/slf4g
    $ go get -u github.com/echocat/slf4g/native
    ```
2. Configure your application to use the selected logger implementation:

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

3. In each package create a logger variable, in best case you create a file named `common.go` or `package.go` which will contain it:

    ```go
    package foo
   
    import (
    	"github.com/echocat/slf4g"
    )
   
    var logger = log.GetLoggerForCurrentPackage()
    ```

4. Now you're ready to go. In every file of this package you can do stuff, like:

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

## Customization

Set the log level globally to Debug

```go
native.DefaultProvider.Level = level.Debug
```

Configure the text formatter to be used.

```go
formatter.Default = formatter.NewText(func (v *formatter.Text) {
	// ... which never colorizes something.
	v.ColorMode = color.ModeNever

	// ... and just prints hours, minutes and seconds
	v.TimeLayout = "150405"
})
```

Configures a writer consumer that writes everything to stdout (instead of stderr; which is the default)

```go
consumer.Default = consumer.NewWriter(os.Stdout)
```

Add an interceptor which will exit the application if someone logs something on level.Fatal or above. This is disabled by default.

```go
interceptor.Default.Add(interceptor.NewFatal())
```

Change the location.Discovery to log everything detail instead of simplified (which is the default).

```go
location.DefaultDiscovery = location.NewCallerDiscovery(func (t *location.CallerDiscovery) {
	t.ReportingDetail = location.CallerReportingDetailDetailed
})
```

## Flags or similar

You can use the package [facade/value](facade/value) to easily configure the logger using flag libraries like the SDK implementation or other compatible ones.

```go
pv := value.NewProvider(native.DefaultProvider)

flag.Var(pv.Consumer.Formatter, "log.format", "Configures the log format.")
flag.Var(pv.Level, "log.level", "Configures the log level.")

flag.Parse()
```

Now you can call you program with:

```bash
$ <myExecutable> -log.format=json -log.level=debug ...
```

## API

How the whole API works in general please refer the [documentation of slf4g](..) directly.
