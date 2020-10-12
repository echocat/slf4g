package log

import (
	"fmt"
	"sync"
	"testing"

	"github.com/echocat/slf4g/fields"

	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_GetLogger(t *testing.T) {
	givenLogger := newMockLogger("foo")
	givenProvider := newMockProvider("bar")
	defer setProvider(givenProvider)()

	givenProvider.provider = func(name string) Logger {
		assert.ToBeEqual(t, "foo", name)
		return givenLogger
	}

	actual := GetLogger("foo")

	assert.ToBeSame(t, givenLogger, actual)
}

func Test_GetRootLogger(t *testing.T) {
	givenLogger := newMockLogger("foo")
	defer setRootLogger(givenLogger)()

	actual := GetRootLogger()

	assert.ToBeSame(t, givenLogger, actual)
}

func Test_IsLevelEnabled(t *testing.T) {
	givenLogger := newMockLogger("foo")
	defer setRootLogger(givenLogger)()

	for _, l := range level.GetProvider().GetLevels() {
		t.Run(fmt.Sprintf("level-%d", l), func(t *testing.T) {
			givenLogger.level = l

			assert.ToBeEqual(t, false, IsLevelEnabled(l-1))
			assert.ToBeEqual(t, true, IsLevelEnabled(l))
			assert.ToBeEqual(t, true, IsLevelEnabled(l+1))
		})
	}
}

func Test_Trace(t *testing.T) {
	givenLogger := newMockLogger("foo")
	givenLogger.initLoggedEvents()
	givenLogger.level = 1
	defer setRootLogger(givenLogger)()
	messageKey := givenLogger.provider.GetFieldKeysSpec().GetMessage()

	Trace()
	Trace(1)
	Trace(1, 2, 3)

	assert.ToBeEqual(t, 3, len(*givenLogger.loggedEvents))
	assert.ToBeEqualUsing(t,
		NewEvent(givenLogger.provider, level.Trace, 2),
		(*givenLogger.loggedEvents)[0],
		IsEventEqual,
	)
	assert.ToBeEqualUsing(t,
		NewEvent(givenLogger.provider, level.Trace, 2).
			With(messageKey, 1),
		(*givenLogger.loggedEvents)[1],
		IsEventEqual,
	)
	assert.ToBeEqualUsing(t,
		NewEvent(givenLogger.provider, level.Trace, 2).
			With(messageKey, []interface{}{1, 2, 3}),
		(*givenLogger.loggedEvents)[2],
		IsEventEqual,
	)
}

func Test_Tracef(t *testing.T) {
	givenLogger := newMockLogger("foo")
	givenLogger.initLoggedEvents()
	givenLogger.level = 1
	defer setRootLogger(givenLogger)()
	messageKey := givenLogger.provider.GetFieldKeysSpec().GetMessage()

	Tracef("hello")
	Tracef("hello %d", 1)

	assert.ToBeEqual(t, 2, len(*givenLogger.loggedEvents))
	assert.ToBeEqualUsing(t,
		NewEvent(givenLogger.provider, level.Trace, 2).
			With(messageKey, fields.LazyFormat("hello")),
		(*givenLogger.loggedEvents)[0],
		IsEventEqual,
	)
	assert.ToBeEqualUsing(t,
		NewEvent(givenLogger.provider, level.Trace, 2).
			With(messageKey, fields.LazyFormat("hello %d", 1)),
		(*givenLogger.loggedEvents)[1],
		IsEventEqual,
	)
}

func Test_IsTraceEnabled(t *testing.T) {
	givenLogger := newMockLogger("foo")
	defer setRootLogger(givenLogger)()

	givenLogger.level = 1
	assert.ToBeEqual(t, true, IsTraceEnabled())
	givenLogger.level = level.Trace
	assert.ToBeEqual(t, true, IsTraceEnabled())
	givenLogger.level = 10000
	assert.ToBeEqual(t, false, IsTraceEnabled())
}

var providerVLock = new(sync.Mutex)

func setProvider(to Provider) func() {
	providerVLock.Lock()
	old := providerV
	providerV = to
	return func() {
		defer providerVLock.Unlock()
		providerV = old
	}
}

func setRootLogger(to *mockLogger) func() {
	to.provider.rootProvider = func() Logger {
		return to
	}
	return setProvider(to.provider)
}
