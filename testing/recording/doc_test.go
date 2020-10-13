package recording_test

import (
	"fmt"

	"github.com/echocat/slf4g/testing/recording"

	log "github.com/echocat/slf4g"
)

func Example() {
	// 1. Create at first a new provider
	provider := recording.NewProvider()
	// 2. Hook it into the global registry AND ensure it will be always executed
	//    at the end of this method
	defer provider.HookGlobally()()

	// 3. Log something with the global logging methods
	log.Info("foo")
	log.Warn("bar")

	// 4. Print every recorded event to stdout
	fmt.Println(provider.GetAll())

	// At the end: Everything will be reset.
}
