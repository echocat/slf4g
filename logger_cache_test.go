package log

import (
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

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

func Test_loggerCache_GetLogger_supportsReentrantFactory(t *testing.T) {
	var instance LoggerCache
	instance = NewLoggerCache(
		func() Logger { return newMockLogger("root") },
		func(name string) Logger {
			if name == "outer" {
				assert.ToBeEqual(t, "inner", instance.GetLogger("inner").GetName())
			}
			return newMockLogger(name)
		},
	)

	done := make(chan Logger)
	go func() {
		done <- instance.GetLogger("outer")
	}()

	select {
	case actual := <-done:
		assert.ToBeEqual(t, "outer", actual.GetName())
	case <-time.After(time.Second):
		t.Fatal("reentrant logger factory deadlocked")
	}
}

func Test_loggerCache_GetLogger_createsConcurrentNameOnce(t *testing.T) {
	var factoryCalls int32
	factoryEntered := make(chan struct{})
	releaseFactory := make(chan struct{})
	instance := NewLoggerCache(
		func() Logger { return newMockLogger("root") },
		func(name string) Logger {
			if atomic.AddInt32(&factoryCalls, 1) == 1 {
				close(factoryEntered)
			}
			<-releaseFactory
			return newMockLogger(name)
		},
	)

	start := make(chan struct{})
	results := make(chan Logger, 32)
	var wait sync.WaitGroup
	for i := 0; i < cap(results); i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			<-start
			results <- instance.GetLogger("same")
		}()
	}
	close(start)
	<-factoryEntered
	close(releaseFactory)
	wait.Wait()
	close(results)

	var expected Logger
	for actual := range results {
		if expected == nil {
			expected = actual
		} else {
			assert.ToBeSame(t, expected, actual)
		}
	}
	assert.ToBeEqual(t, int32(1), atomic.LoadInt32(&factoryCalls))
}

func Test_loggerCache_GetLogger_recoversAfterFactoryPanic(t *testing.T) {
	var factoryCalls int
	instance := NewLoggerCache(
		func() Logger { return newMockLogger("root") },
		func(name string) Logger {
			factoryCalls++
			if factoryCalls == 1 {
				panic("expected")
			}
			return newMockLogger(name)
		},
	)

	assert.Execution(t, func() {
		instance.GetLogger("name")
	}).WillPanicWith("expected")
	assert.ToBeEqual(t, "name", instance.GetLogger("name").GetName())
	assert.ToBeEqual(t, 2, factoryCalls)
}

func Test_loggerCache_GetNames(t *testing.T) {
	instance := &loggerCache{
		root:    newMockLogger("root"),
		loggers: make(map[string]Logger),
		factory: func(name string) Logger { return newMockLogger(name) },
	}

	actualBefore := instance.GetNames()
	assert.ToBeEqual(t, []string{}, actualBefore)

	assert.ToBeNotNil(t, instance.GetRootLogger())
	assert.ToBeNotNil(t, instance.GetLogger("foo"))
	assert.ToBeNotNil(t, instance.GetLogger("bar"))

	actualAfter := instance.GetNames()
	sort.Strings(actualAfter)
	assert.ToBeEqual(t, []string{"bar", "foo"}, actualAfter)
}
