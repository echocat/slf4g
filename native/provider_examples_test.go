package native_test

import (
	"os"

	"github.com/echocat/slf4g/native/location"

	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/native"
	"github.com/echocat/slf4g/native/color"
	"github.com/echocat/slf4g/native/consumer"
	"github.com/echocat/slf4g/native/formatter"
	"github.com/echocat/slf4g/native/interceptor"
)

func Example_customization() {
	// Set the log level globally to Debug
	native.DefaultProvider.Level = level.Debug

	// Configure the console formatter to be used.
	formatter.Default = formatter.NewConsole(func(v *formatter.Console) {
		// ... which never colorizes something.
		v.ColorMode = color.ModeNever

		// ... and just prints hours, minutes and seconds
		v.TimeLayout = "150405"
	})

	// Configures a writer consumer that writes everything to stdout (instead
	// of stderr; which is the default)
	consumer.Default = consumer.NewWriter(os.Stdout)

	// Add an interceptor which will exit the application if someone logs
	// something on level.Fatal or above. This is disabled by default.
	interceptor.Default.Add(interceptor.NewFatal())

	// Change the location.Discovery to log everything detail instead of
	// simplified (which is the default).
	location.DefaultDiscovery = location.NewCallerAwareDiscovery(func(t *location.CallerAwareDiscovery) {
		t.Detail = location.CallerAwareDetailDetailed
	})
}
