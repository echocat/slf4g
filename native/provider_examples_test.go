package native_test

import (
	"os"

	log "github.com/echocat/slf4g"
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

	var minMessageWidth int16 = 20

	// Configure the text formatter to be used.
	formatter.Default = formatter.NewText(func(v *formatter.Text) {
		// ... which never colorizes something.
		v.ColorMode = color.ModeNever

		// ... with a minimal message width of fixed 20
		v.MinMessageWidth = &minMessageWidth

		// ... print nothing for the time
		//     (we have to create reproducible output for this example 😉)
		v.TimeLayout = " "
	})

	// Configures a writer consumer that writes everything to stdout (instead
	// of stderr; which is the default)
	consumer.Default = consumer.NewWriter(os.Stdout)

	// Add an interceptor which will exit the application if someone logs
	// something on level.Fatal or above. This is disabled by default.
	interceptor.Default.Add(interceptor.NewFatal())

	// Change the location.Discovery to log everything detail instead of
	// simplified (which is the default).
	location.DefaultDiscovery = location.NewCallerDiscovery(func(t *location.CallerDiscovery) {
		t.ReportingDetail = location.CallerReportingDetailDetailed
	})

	log.Info("Hello, world!")

	// Output:
	// [ INFO] Hello, world!        location=github.com/echocat/slf4g/native_test.Example_customization:50
}
