package testlog

import (
	"sync"
	"testing"
	_ "unsafe"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
	tlevel "github.com/echocat/slf4g/sdk/testlog/level"
)

// NewProvider creates a new instance of Provider which is ready to use.
//
// tb should hold an instance of either *testing.T, *testing.B or *testing.F.
//
// customizer can be used to change the behavior of the Provider:
//   - Level
//   - FailAtLevel
//   - FailNowAtLevel
//   - TimeFormat
//   - LevelFormatter
//   - Name
//   - AllLevels
//   - FieldKeysSpec
func NewProvider(tb testing.TB, customizer ...func(*Provider)) *Provider {
	result := &Provider{tb: tb, startedNs: runtimeNano()}

	for _, c := range customizer {
		c(result)
	}

	return result
}

const (
	// DefaultLevel specifies the default level.Level of an instance of Provider
	// which be used if no other level was defined.
	DefaultLevel = level.Debug

	// NeverFailLevel is used for Provider.FailAtLevel and indicates that regardless
	// at which level each log.Event is logged, this event will never lead to a fail of
	// the tests.
	NeverFailLevel = level.Level(65535)

	// NoopTimeFormat tells the Provider to not print any timestamp in the log messages.
	NoopTimeFormat = "<noop time format>"

	// SinceTestStartedMcsTimeFormat tells the Provider to print only the microseconds since
	// the test started (is based on Hook() and/or NewProvider()).
	SinceTestStartedMcsTimeFormat = "<since test started>"
)

var (
	// DefaultFailAtLevel is used if FailAtLevel was not used.
	DefaultFailAtLevel = level.Error

	// DefaultFailNowAtLevel is used if FailNowAtLevel was not used.
	DefaultFailNowAtLevel = level.Fatal

	// DefaultTimeFormat is used if TimeFormat was not used.
	DefaultTimeFormat = SinceTestStartedMcsTimeFormat
)

// Provider is an implementation of log.Provider which ensures that everything is
// logged using testing.TB#Log(). Use NewProvider(..) to get a new instance.
type Provider struct {
	tb        testing.TB
	startedNs int64

	name           string
	level          level.Level
	allLevels      level.Levels
	fieldKeysSpec  fields.KeysSpec
	failAtLevel    level.Level
	failNowAtLevel level.Level
	timeFormat     string
	levelFormatter tlevel.Formatter

	coreRootLogger *coreLogger
	rootLogger     log.Logger
	initRootLogger sync.Once

	// For testing only
	interceptLogDepth func(string, uint16)
	interceptFail     func()
	interceptFailNow  func()
}

//go:linkname runtimeNano runtime.nanotime
func runtimeNano() int64

func (instance *Provider) initIfRequired() {
	instance.initRootLogger.Do(func() {
		instance.coreRootLogger = &coreLogger{instance, RootLoggerName, 0}
		instance.rootLogger = log.NewLogger(instance.coreRootLogger)
	})
}

// GetRootLogger implements log.Provider#GetRootLogger()
func (instance *Provider) GetRootLogger() log.Logger {
	instance.initIfRequired()
	return instance.rootLogger
}

// GetLogger implements log.Provider#GetLogger()
func (instance *Provider) GetLogger(name string) log.Logger {
	if name == RootLoggerName {
		return instance.GetRootLogger()
	}

	instance.initIfRequired()
	return log.NewLogger(&coreLogger{instance, name, 0})
}

// GetName implements log.Provider#GetName()
func (instance *Provider) GetName() string {
	if v := instance.name; v != "" {
		return v
	}
	return instance.tb.Name()
}

// GetAllLevels implements log.Provider#GetAllLevels()
func (instance *Provider) GetAllLevels() level.Levels {
	if v := instance.allLevels; v != nil {
		return v
	}
	return level.GetProvider().GetLevels()
}

// GetFieldKeysSpec implements log.Provider#GetFieldKeysSpec()
func (instance *Provider) GetFieldKeysSpec() fields.KeysSpec {
	if v := instance.fieldKeysSpec; v != nil {
		return v
	}
	return &fields.KeysSpecImpl{}
}

// GetLevel returns the current level.Level where this log.Provider is set to.
func (instance *Provider) GetLevel() level.Level {
	if v := instance.level; v != 0 {
		return v
	}
	return DefaultLevel
}

// SetLevel changes the current level.Level of this log.Provider. If set to
// 0 it will force this Provider to use DefaultLevel.
func (instance *Provider) SetLevel(v level.Level) {
	instance.level = v
}

func (instance *Provider) getFailAtLevel() level.Level {
	if v := instance.failAtLevel; v != 0 {
		return v
	}
	return DefaultFailAtLevel
}

func (instance *Provider) getFailNowAtLevel() level.Level {
	if v := instance.failNowAtLevel; v != 0 {
		return v
	}
	return DefaultFailNowAtLevel
}

func (instance *Provider) getTimeFormat() string {
	if v := instance.timeFormat; v != "" {
		return v
	}
	return DefaultTimeFormat
}

func (instance *Provider) getLevelFormatter() tlevel.Formatter {
	if v := instance.levelFormatter; v != nil {
		return v
	}
	return tlevel.DefaultFormatter
}

// Level specifies the level of the Provider which will be also inherited
// by all of its loggers. By default, the Provider will use DefaultLevel.
func Level(v level.Level) func(*Provider) {
	return func(provider *Provider) {
		provider.level = v
	}
}

// FailAtLevel defines a level.Level at which log.Event will lead to a test failure
// after each code of the test has passed (in contrast to FailNowAtLevel which will
// fail immediately) if they're logged with this a log.Logger handled by the
// Provider. If set to NeverFailLevel nothing happens. By default, the Provider
// will use DefaultFailAtLevel.
func FailAtLevel(v level.Level) func(*Provider) {
	return func(provider *Provider) {
		provider.failAtLevel = v
	}
}

// FailNowAtLevel defines a level.Level at which log.Event will lead to a test
// fails immediately (in contrast to FailAtLevel which allows to test to finish)
// if they're logged with this a log.Logger handled by the Provider. If set to
// NeverFailLevel nothing happens. By default, the Provider will use
// DefaultFailNowAtLevel.
func FailNowAtLevel(v level.Level) func(*Provider) {
	return func(provider *Provider) {
		provider.failNowAtLevel = v
	}
}

// TimeFormat defines how each entry will be formatted on print. If NoopTimeFormat
// is used, nothing will be printed. If SinceTestStartedMcsTimeFormat will be used
// no time is printed, but the microseconds since the test started. By default, the
// Provider will use DefaultTimeFormat. See time.Layout for more details.
func TimeFormat(v string) func(*Provider) {
	return func(provider *Provider) {
		provider.timeFormat = v
	}
}

// LevelFormatter formats the levels on printing. By default, the Provider will use
// level.DefaultFormatter.
func LevelFormatter(v tlevel.Formatter) func(*Provider) {
	return func(provider *Provider) {
		provider.levelFormatter = v
	}
}

// Name specifies the name of the Provider. By default, the Provider will use
// testing.TB#Name().
func Name(v string) func(*Provider) {
	return func(provider *Provider) {
		provider.name = v
	}
}

// AllLevels specifies the levels which are supported by the Provider and all
// of its loggers. By default, the Provider will use level.GetProvider()#GetLevels().
func AllLevels(v level.Levels) func(*Provider) {
	return func(provider *Provider) {
		provider.allLevels = v
	}
}

// FieldKeysSpec specifies the spec of the fields are supported by the Provider
// and all of its loggers. By default, the Provider will use the default instance
// of fields.KeysSpecImpl.
func FieldKeysSpec(v fields.KeysSpec) func(*Provider) {
	return func(provider *Provider) {
		provider.fieldKeysSpec = v
	}
}
