package color_test

import (
	"os"

	"github.com/echocat/slf4g/native/color"
)

func ExampleDetectSupportForWriter_detection() {
	prepared, supported, err := color.DetectSupportForWriter(os.Stderr)
	if err != nil {
		panic(err)
	}

	msg := []byte("Hello, world!")
	if supported.IsSupported() {
		msg = colorize(msg)
	}
	_, _ = prepared.Write(msg)
}

func colorize(in []byte) []byte {
	result := in
	// Great stuff happens to colorize...
	return result
}
