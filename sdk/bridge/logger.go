package sdk

import (
	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/level"
)

// Logger is an interface which describes instances which are compatible with
// the SDK Logger instance.
type Logger interface {
	// Print calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Print.
	Print(...interface{})
	// Printf calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Printf.
	Printf(string, ...interface{})
	// Println calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Println.
	Println(...interface{})

	// Fatal is equivalent to l.Print() and can be followed by a call to
	// os.Exit(1).
	Fatal(...interface{})
	// Fatalf is equivalent to l.Printf() and can be followed by a call to
	// os.Exit(1).
	Fatalf(string, ...interface{})
	// Fatalln is equivalent to l.Println() and can be followed by a call to
	// os.Exit(1).
	Fatalln(...interface{})

	// Panic is equivalent to l.Print() and can followed by a call to panic().
	Panic(...interface{})
	// Panicf is equivalent to l.Printf() and can followed by a call to panic().
	Panicf(string, ...interface{})
	// Panicln is equivalent to l.Println() and can followed by a call to
	// panic().
	Panicln(...interface{})
}

// NewLogger creates a new instance of an SDK compatible Logger which forwards
// all it's events to the provided log.CoreLogger.
//
// printLevel defines the level which is used to log to the given log.CoreLogger
// on every Logger.Print(), Logger.Printf() and Logger.Println() event.
func NewLogger(target log.CoreLogger, printLevel level.Level, customizer ...func(*LoggerImpl)) Logger {
	result := &LoggerImpl{
		Delegate:   target,
		PrintLevel: printLevel,
	}

	for _, c := range customizer {
		c(result)
	}

	return result
}
