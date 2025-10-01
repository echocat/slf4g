package color_test

import (
	"bytes"
	"fmt"

	"github.com/echocat/slf4g/native/color"
)

func ExampleDetectSupportForWriter_detection() {
	var buf bytes.Buffer // This does not support colors at all.

	prepared, supported, err := color.DetectSupportForWriter(&buf)
	if err != nil {
		panic(err)
	}

	msg := []byte("Hello, world!")
	if supported.IsSupported() {
		msg = colorize(msg)
	}
	_, _ = prepared.Write(msg)

	fmt.Println(buf.String())

	// Output:
	// Hello, world!
}

func colorize(_ []byte) []byte {
	panic("should never be called.")
}
