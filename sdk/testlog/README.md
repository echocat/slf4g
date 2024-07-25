[![PkgGoDev](https://pkg.go.dev/badge/github.com/echocat/slf4g/sdk/testlog)](https://pkg.go.dev/github.com/echocat/slf4g/sdk/testlog)

# slf4g testing Logger implementation

Provides a Logger which will be connected to [`testing.T.Log()`](https://pkg.go.dev/testing#T.Log) of the go SDK.

If you're looking for an instance to record all logged events see [`github.com/echocat/slf4g/testing/recording`](../../testing/recording) package.

## Usage

The easiest way to enable the slf4g framework in your tests is, simply:

```golang
package foo

import (
	"testing"
	"github.com/echocat/slf4g"
	"github.com/echocat/slf4g/sdk/testlog"
)

func TestMyGreatStuff(t *testing.T) {
	testlog.Hook(t)

	log.Info("Yeah! This is a log!")
}
```

... that's it!

See [`Hook(..)`](hook.go) for more details.
