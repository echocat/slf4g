package log

import (
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_NewLoggerCache(t *testing.T) {
	givenRootLogger := newMockLogger("aRootOne")
	givenRootFactory := func() Logger { return givenRootLogger }
	givenFactory := func(name string) Logger { return newMockLogger(name) }

	actual := NewLoggerCache(givenRootFactory, givenFactory)

	assert.ToBeOfType(t, &loggerCache{}, actual)
	assert.ToBeSame(t, givenRootLogger, actual.(*loggerCache).root)
	assert.ToBeSame(t, givenFactory, actual.(*loggerCache).factory)
	assert.ToBeNotNil(t, actual.(*loggerCache).loggers)
}

func Test_NewLoggerCache_panicsIfRootFactoryReturnsNil(t *testing.T) {
	givenRootFactory := func() Logger { return nil }
	givenFactory := func(name string) Logger { return newMockLogger(name) }

	assert.Execution(t, func() {
		NewLoggerCache(givenRootFactory, givenFactory)
	}).WillPanicWith("Root factory returned a nil root logger.")
}

func Test_loggerCache_GetRootLogger(t *testing.T) {
	givenRootLogger := newMockLogger("aRootOne")
	instance := &loggerCache{root: givenRootLogger}

	actual1 := instance.GetRootLogger()
	actual2 := instance.GetRootLogger()

	assert.ToBeSame(t, givenRootLogger, actual1)
	assert.ToBeSame(t, givenRootLogger, actual2)
}

func Test_loggerCache_GetLogger(t *testing.T) {
	instance := &loggerCache{
		loggers: make(map[string]Logger),
		factory: func(name string) Logger { return newMockLogger(name) },
	}

	actualFoo1 := instance.GetLogger("foo")
	actualFoo2 := instance.GetLogger("foo")
	actualBar1 := instance.GetLogger("bar")
	actualBar2 := instance.GetLogger("bar")

	assert.ToBeNotNil(t, actualFoo1)
	assert.ToBeNotNil(t, actualFoo2)
	assert.ToBeNotNil(t, actualBar1)
	assert.ToBeNotNil(t, actualBar2)

	assert.ToBeEqual(t, "foo", actualFoo1.GetName())
	assert.ToBeEqual(t, "foo", actualFoo2.GetName())
	assert.ToBeEqual(t, "bar", actualBar1.GetName())
	assert.ToBeEqual(t, "bar", actualBar2.GetName())

	assert.ToBeSame(t, actualFoo1, actualFoo2)
	assert.ToBeSame(t, actualBar1, actualBar2)

	assert.ToBeNotSame(t, actualFoo1, actualBar1)
	assert.ToBeNotSame(t, actualFoo2, actualBar2)
}

func Test_loggerCache_GetLogger_returnsRootIfFactoryReturnsNil(t *testing.T) {
	givenRootLogger := newMockLogger("root")
	instance := &loggerCache{
		root:    givenRootLogger,
		loggers: make(map[string]Logger),
		factory: func(name string) Logger { return nil },
	}

	actual1 := instance.GetLogger("foo")
	actual2 := instance.GetLogger("bar")

	assert.ToBeNotNil(t, actual1)
	assert.ToBeNotNil(t, actual2)

	assert.ToBeEqual(t, "root", actual1.GetName())
	assert.ToBeEqual(t, "root", actual2.GetName())

	assert.ToBeSame(t, givenRootLogger, actual1)
	assert.ToBeSame(t, givenRootLogger, actual2)
}
