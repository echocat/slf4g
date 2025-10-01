package color_test

import (
	"os"

	"github.com/echocat/slf4g/native/color"
)

func ExampleDetectSupportForWriter_detection() {
	// For this test we force the output to support colors.
	// Usually, you just use os.Stdout or os.Stderr.
	output := color.ForcedSupportedWriteFunc(os.Stdout.Write)

	prepared, supported, err := color.DetectSupportForWriter(output)
	if err != nil {
		panic(err)
	}

	msg := "Hello, world!"
	if supported.IsSupported() {
		msg = "colored(" + msg + ")"
	}
	_, _ = prepared.Write([]byte(msg))

	// Output:
	// colored(Hello, world!)
}
