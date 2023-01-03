package interceptor

import (
	"fmt"
	"math"
	"testing"

	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/testing/recording"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_NewFatal(t *testing.T) {
	actual := NewFatal()

	assert.ToBeEqual(t, 13, actual.ExitCode)
}

func Test_NewFatal_customized(t *testing.T) {
	actual := NewFatal(func(fatal *Fatal) {
		fatal.ExitCode = 66
	})

	assert.ToBeEqual(t, 66, actual.ExitCode)
}

func Test_Fatal_OnBeforeLog(t *testing.T) {
	instance := NewFatal()
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEventWithFields(level.Warn, fields.With("foo", "bar"))

	actual := instance.OnBeforeLog(givenEvent, givenLogger.Provider)

	assert.ToBeSame(t, givenEvent, actual)
}

func Test_Fatal_OnAfterLog_doNothing(t *testing.T) {
	old := fatalExit

	instance := NewFatal()
	givenLogger := recording.NewLogger()

	for _, l := range level.GetProvider().GetLevels() {
		if level.Fatal.CompareTo(l) > 0 {
			t.Run(fmt.Sprint(l), func(t *testing.T) {
				defer func() { fatalExit = old }()

				fatalExit = func(exitCode int) {
					t.Fatalf("Expected to never be called; but was called with <%d>", exitCode)
				}
				givenEvent := givenLogger.NewEventWithFields(l, fields.With("foo", "bar"))

				actual := instance.OnAfterLog(givenEvent, givenLogger.Provider)

				assert.ToBeEqual(t, true, actual)
			})
		}
	}
}

func Test_Fatal_OnAfterLog_exists(t *testing.T) {
	old := fatalExit

	instance := NewFatal()
	givenLogger := recording.NewLogger()

	for _, l := range level.GetProvider().GetLevels() {
		if level.Fatal.CompareTo(l) <= 0 {
			t.Run(fmt.Sprint(l), func(t *testing.T) {
				fatalExitWasCalled := false
				defer func() { fatalExit = old }()

				fatalExit = func(actualExitCode int) {
					fatalExitWasCalled = true
					assert.ToBeEqual(t, 13, actualExitCode)
				}
				givenEvent := givenLogger.NewEventWithFields(l, fields.With("foo", "bar"))

				actual := instance.OnAfterLog(givenEvent, givenLogger.Provider)

				assert.ToBeEqual(t, false, actual)
				assert.ToBeEqual(t, true, fatalExitWasCalled)
			})
		}
	}
}

func Test_Fatal_GetPriority(t *testing.T) {
	instance := NewFatal()

	actual := instance.GetPriority()

	assert.ToBeEqual(t, int16(math.MaxInt16), actual)
}
