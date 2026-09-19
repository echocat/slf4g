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
	Print(...any)
	// Printf calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Printf.
	Printf(string, ...any)
	// Println calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Println.
	Println(...any)

	// Fatal is equivalent to l.Print() and can be followed by a call to
	// os.Exit(1).
	Fatal(...any)
	// Fatalf is equivalent to l.Printf() and can be followed by a call to
	// os.Exit(1).
	Fatalf(string, ...any)
	// Fatalln is equivalent to l.Println() and can be followed by a call to
	// os.Exit(1).
	Fatalln(...any)

	// Panic is equivalent to l.Print() and can followed by a call to panic().
	Panic(...any)
	// Panicf is equivalent to l.Printf() and can followed by a call to panic().
	Panicf(string, ...any)
	// Panicln is equivalent to l.Println() and can followed by a call to
	// panic().
	Panicln(...any)
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
