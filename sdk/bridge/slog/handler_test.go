//go:build go1.21

package sdk

import (
	"context"
	"fmt"
	sdk "log/slog"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/testing/recording"
)

func TestNewHandler(t *testing.T) {
	aLogger := recording.NewCoreLogger()
	anotherHandler := &Handler{}

	actual := NewHandler(aLogger, func(v *Handler) {
		v.parent = anotherHandler
	}, func(v *Handler) {
		v.fieldKeyPath = newFieldKeyPath(nil, "foo")
	})

	assert.ToBeNotNil(t, actual)
	assert.ToBeSame(t, aLogger, actual.Delegate)
	assert.ToBeEqual(t, "foo.", actual.fieldKeyPath.prefix())
	assert.ToBeSame(t, anotherHandler, actual.parent)
}

func TestHandler_Enabled(t *testing.T) {
	aLogger := recording.NewCoreLogger()
	aLogger.SetLevel(level.Warn)

	instance := &Handler{
		Delegate: aLogger,
	}

	cases := []struct {
		givenLevel sdk.Level
		expected   bool
	}{
		{LevelDebug, false},
		{LevelInfo, false},
		{LevelWarn - 1, false},
		{LevelWarn, true},
		{LevelWarn + 1, true},
		{LevelFatal, true},
		{LevelFatal + 1, true},
	}

	for _, c := range cases {
		t.Run(c.givenLevel.String(), func(t *testing.T) {
			actual := instance.Enabled(context.TODO(), c.givenLevel)
			assert.ToBeEqual(t, c.expected, actual)
		})
	}
}

func TestHandler_Handle(t *testing.T) {

	aTime, err := time.Parse(time.RFC3339, "2025-10-01T15:30:15Z")
	assert.ToBeNoError(t, err)

	cases := []struct {
		name          string
		givenLevel    sdk.Level
		expectedLevel level.Level
		expectedError string
	}{
		{"regular", LevelInfo, level.Info, ""},
		{"custom", LevelWarn + 1, level.Warn, ""},
		{"above fatal", 666, level.Fatal, ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			baseLogger := recording.NewCoreLogger()
			helperWasCalled := false

			instance := &Handler{
				Delegate: &delegateCoreLoggerWithHelper{baseLogger, func() {
					helperWasCalled = true
				}},
				DetectSkipFrames: func(skip uint16) uint16 {
					return skip + 1
				},
			}

			actualErr := instance.Handle(context.TODO(), sdk.Record{
				Time:    aTime,
				Message: "aMessage",
				Level:   c.givenLevel,
				PC:      666,
			})
			assert.ToBeEqual(t, true, helperWasCalled)
			if expectedErr := c.expectedError; expectedErr == "" {
				assert.ToBeEqual(t, 1, baseLogger.Len())

				actual := baseLogger.Get(0)
				assert.ToBeEqual(t, c.expectedLevel, actual.GetLevel())
				_, withProgramCounter := actual.(interface{ GetProgramCounter() uintptr })
				assert.ToBeEqual(t, false, withProgramCounter)

				actualAsMap, err := fields.AsMap(actual)
				assert.ToBeNoError(t, err)

				assert.ToBeEqual(t, map[string]any{
					"logger":    baseLogger,
					"timestamp": aTime,
					"message":   "aMessage",
				}, actualAsMap)
			} else {
				assert.ToBeMatching(t, expectedErr, actualErr)
			}
		})
	}
}

func TestHandler_Handle_preservesProgramCounter(t *testing.T) {
	baseLogger := recording.NewCoreLogger()
	instance := &Handler{Delegate: baseLogger}

	actualErr := instance.Handle(context.Background(), sdk.Record{Level: LevelInfo, PC: 666})

	assert.ToBeNoError(t, actualErr)
	actual := baseLogger.Get(0)
	actualPC, ok := actual.(interface{ GetProgramCounter() uintptr })
	assert.ToBeEqual(t, true, ok)
	assert.ToBeEqual(t, uintptr(666), actualPC.GetProgramCounter())
}

func TestHandler_Handle_respectsCustomizedDefaultDetectSkipFramesWithProgramCounter(t *testing.T) {
	original := DefaultDetectSkipFrames
	defer func() { DefaultDetectSkipFrames = original }()
	detectorCalled := false
	DefaultDetectSkipFrames = func(skip uint16) uint16 {
		detectorCalled = true
		return skip
	}
	baseLogger := recording.NewCoreLogger()
	instance := &Handler{Delegate: baseLogger}

	actualErr := instance.Handle(context.Background(), sdk.Record{Level: LevelInfo, PC: 666})

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, true, detectorCalled)
	actual := baseLogger.Get(0)
	_, withProgramCounter := actual.(interface{ GetProgramCounter() uintptr })
	assert.ToBeEqual(t, false, withProgramCounter)
}

func TestHandler_Handle_respectsRejectedProgramCounterEvent(t *testing.T) {
	baseLogger := recording.NewCoreLogger()
	instance := &Handler{Delegate: &programCounterRejectingCoreLogger{baseLogger}}

	actualErr := instance.Handle(context.Background(), sdk.Record{Level: LevelInfo, PC: 666})

	assert.ToBeNoError(t, actualErr)
	actual := baseLogger.Get(0)
	_, ok := actual.(interface{ GetProgramCounter() uintptr })
	assert.ToBeEqual(t, false, ok)
}

func TestHandler_eventOfRecord(t *testing.T) {
	logger := recording.NewCoreLogger()
	aTime, err := time.Parse(time.RFC3339, "2025-10-01T15:30:15Z")
	assert.ToBeNoError(t, err)

	cases := []struct {
		name           string
		instance       *Handler
		attrs          attrs
		expectedLevel  level.Level
		expectedFields map[string]any
	}{
		{
			"simple",
			&Handler{},
			nil,
			level.Fatal,
			map[string]any{
				"timestamp": aTime,
				"message":   "aMessage",
			},
		},
		{
			"recordAttrs",
			&Handler{},
			attrs{
				sdk.Int("foo", 1),
				sdk.Int("bar", 2),
			},
			level.Fatal,
			map[string]any{
				"timestamp": aTime,
				"message":   "aMessage",
				"foo":       int64(1),
				"bar":       int64(2),
			},
		},
		{
			"groupedRecordAttrs",
			(&Handler{}).WithGroup("outer").WithGroup("inner").(*Handler),
			attrs{sdk.Int("foo", 1)},
			level.Fatal,
			map[string]any{
				"timestamp":       aTime,
				"message":         "aMessage",
				"outer.inner.foo": int64(1),
			},
		},
		{
			"recordLogValuer",
			&Handler{},
			attrs{
				sdk.Any("secret", redactingLogValuer{"not-for-the-log"}),
			},
			level.Fatal,
			map[string]any{
				"timestamp": aTime,
				"message":   "aMessage",
				"secret":    "[REDACTED]",
			},
		},
		{
			"handlerWithAttrs",
			&Handler{attrs: attrs{
				sdk.Int("foo", 1),
				sdk.Int("bar", 2),
			}},
			nil,
			level.Fatal,
			map[string]any{
				"timestamp": aTime,
				"message":   "aMessage",
				"foo":       int64(1),
				"bar":       int64(2),
			},
		},
		{
			"handlerWithLogValuer",
			&Handler{attrs: attrs{
				sdk.Any("secret", redactingLogValuer{"not-for-the-log"}),
			}},
			nil,
			level.Fatal,
			map[string]any{
				"timestamp": aTime,
				"message":   "aMessage",
				"secret":    "[REDACTED]",
			},
		},
		{
			"handlerWithAttrs_recordWithAttrs",
			&Handler{attrs: attrs{
				sdk.Int("foo", 11),
				sdk.Int("xyz", 13),
			}},
			attrs{
				sdk.Int("bar", 22),
				sdk.Int("xyz", 23),
			},
			level.Fatal,
			map[string]any{
				"timestamp": aTime,
				"message":   "aMessage",
				"foo":       int64(11),
				"bar":       int64(22),
				"xyz":       int64(23),
			},
		},
		{
			"parentWithoutAttrs_handlerWithAttrs_recordWithAttrs",
			&Handler{
				attrs: attrs{
					sdk.Int("foo", 11),
					sdk.Int("xyz", 13),
				},
				parent: &Handler{},
			},
			attrs{
				sdk.Int("bar", 22),
				sdk.Int("xyz", 23),
			},
			level.Fatal,
			map[string]any{
				"timestamp": aTime,
				"message":   "aMessage",
				"foo":       int64(11),
				"bar":       int64(22),
				"xyz":       int64(23),
			},
		},
		{
			"parentWitAttrs_handlerWithAttrs_recordWithAttrs",
			&Handler{
				attrs: attrs{
					sdk.Int("foo", 11),
					sdk.Int("xyz", 13),
				},
				parent: &Handler{attrs: attrs{
					sdk.Int("foo", 31),
					sdk.Int("abc", 34),
				}},
			},
			attrs{
				sdk.Int("bar", 22),
				sdk.Int("xyz", 23),
			},
			level.Fatal,
			map[string]any{
				"timestamp": aTime,
				"message":   "aMessage",
				"foo":       int64(11),
				"bar":       int64(22),
				"xyz":       int64(23),
				"abc":       int64(34),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			aRecord := sdk.Record{
				Time:    aTime,
				Message: "aMessage",
				Level:   LevelFatal,
				PC:      666,
			}
			aRecord.AddAttrs(c.attrs...)

			actual, actualErr := c.instance.eventOfRecord(logger, aRecord)
			assert.ToBeNoError(t, actualErr)
			assert.ToBeEqual(t, c.expectedLevel, actual.GetLevel())

			actualAsMap, err := fields.AsMap(actual)
			assert.ToBeNoError(t, err)

			assert.ToBeEqual(t, c.expectedFields, actualAsMap)
		})
	}
}

func TestHandler_eventOfRecord_withErrorInLevel(t *testing.T) {
	aLogger := recording.NewCoreLogger()
	aLevelMapper := &testingLevelMapper{func(v sdk.Level) (level.Level, error) {
		return 0, fmt.Errorf("illegal level: %d", v)
	}, nil}

	instance := &Handler{
		LevelMapper: aLevelMapper,
	}

	actual, actualErr := instance.eventOfRecord(aLogger, sdk.Record{Level: sdk.Level(666)})
	assert.ToBeMatching(t, "illegal level: 666", actualErr)
	assert.ToBeNil(t, actual)
}

func TestHandler_levelOfRecord(t *testing.T) {
	aLevelMapper := &testingLevelMapper{func(v sdk.Level) (level.Level, error) {
		switch v {
		case sdk.Level(666):
			return 666, nil
		default:
			return 0, fmt.Errorf("illegal level: %d", v)
		}
	}, nil}

	instance := &Handler{
		LevelMapper: aLevelMapper,
	}

	cases := []struct {
		name          string
		givenLevel    sdk.Level
		expectedLevel level.Level
		expectedErr   string
	}{
		{"666", 666, 666, ""},
		{"illegal", 123, 0, "illegal level: 123"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, actualErr := instance.levelOfRecord(sdk.Record{Level: c.givenLevel})
			if expectedErr := c.expectedErr; expectedErr == "" {
				assert.ToBeNoError(t, actualErr)
				assert.ToBeEqual(t, c.expectedLevel, actual)
			} else {
				assert.ToBeMatching(t, expectedErr, actualErr)
				assert.ToBeEqual(t, c.expectedLevel, actual)
			}
		})
	}
}

func TestHandler_WithAttrs(t *testing.T) {
	aCoreLogger := recording.NewCoreLogger()
	aLevelMapper := NewLevelMapperFacade(nil)
	aDetectSkipFrames := func(uint16) uint16 { panic("should never be called") }
	someAttrs := attrs{
		sdk.Int("foo", -1),
		sdk.Int("foo.foo", 1),
		sdk.Int("foo.xyz", 123),
	}

	instance := &Handler{
		Delegate:         aCoreLogger,
		LevelMapper:      aLevelMapper,
		DetectSkipFrames: aDetectSkipFrames,
		parent:           nil,
		fieldKeyPath:     newFieldKeyPath(nil, "foo"),
		attrs:            someAttrs,
	}

	actual := instance.WithAttrs(attrs{
		sdk.Int("bar", 2),
		sdk.Int("xyz", 666),
	})
	assert.ToBeNotSame(t, instance, actual)
	assert.ToBeOfType(t, (*Handler)(nil), actual)
	actualC := actual.(*Handler)
	assert.ToBeSame(t, instance.Delegate, actualC.Delegate)
	assert.ToBeSame(t, instance.LevelMapper, actualC.LevelMapper)
	assert.ToBeSame(t, instance.DetectSkipFrames, actualC.DetectSkipFrames)
	assert.ToBeSame(t, instance, actualC.parent)
	assert.ToBeSame(t, instance.fieldKeyPath, actualC.fieldKeyPath)
	assert.ToBeEqual(t, "foo.", actualC.fieldKeyPath.prefix())
	assert.ToBeEqual(t, attrs{
		sdk.Int("foo.bar", 2),
		sdk.Int("foo.xyz", 666),
	}, actualC.attrs)
}

func TestHandler_fields_deepLineage(t *testing.T) {
	const depth = 20_000
	var instance sdk.Handler = &Handler{}
	for i := 0; i < depth; i++ {
		instance = instance.WithAttrs([]sdk.Attr{sdk.Int(strconv.Itoa(i), i)})
	}

	actual := instance.(*Handler).fields()

	assert.ToBeEqual(t, depth, actual.Len())
	actualFirst, actualFirstExists := actual.Get("0")
	assert.ToBeEqual(t, int64(0), actualFirst)
	assert.ToBeEqual(t, true, actualFirstExists)
	actualLast, actualLastExists := actual.Get(strconv.Itoa(depth - 1))
	assert.ToBeEqual(t, int64(depth-1), actualLast)
	assert.ToBeEqual(t, true, actualLastExists)
}

func TestHandler_fields_preservesLineageOrder(t *testing.T) {
	instance := (&Handler{}).
		WithAttrs([]sdk.Attr{sdk.Int("first", 1), sdk.Int("replaced", 1)}).
		WithAttrs([]sdk.Attr{sdk.Int("replaced", 2), sdk.Int("last", 3)}).
		WithGroup("group").
		WithAttrs([]sdk.Attr{sdk.Int("child", 4)}).(*Handler)

	var actualKeys []string
	actualErr := instance.fields().ForEach(func(key string, _ any) error {
		actualKeys = append(actualKeys, key)
		return nil
	})

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, []string{"group.child", "first", "replaced", "last"}, actualKeys)
	actualReplaced, actualReplacedExists := instance.fields().Get("replaced")
	assert.ToBeEqual(t, int64(2), actualReplaced)
	assert.ToBeEqual(t, true, actualReplacedExists)
}

func TestHandler_WithGroup(t *testing.T) {
	aCoreLogger := recording.NewCoreLogger()
	aLevelMapper := NewLevelMapperFacade(nil)
	aDetectSkipFrames := func(uint16) uint16 { panic("should never be called") }
	someAttrs := attrs{sdk.Int("foo", 1)}

	instance := &Handler{
		Delegate:         aCoreLogger,
		LevelMapper:      aLevelMapper,
		DetectSkipFrames: aDetectSkipFrames,
		parent:           nil,
		fieldKeyPath:     newFieldKeyPath(nil, "foo"),
		attrs:            someAttrs,
	}

	actual := instance.WithGroup("bar")
	assert.ToBeNotSame(t, instance, actual)
	assert.ToBeOfType(t, (*Handler)(nil), actual)
	actualC := actual.(*Handler)
	assert.ToBeSame(t, instance.Delegate, actualC.Delegate)
	assert.ToBeSame(t, instance.LevelMapper, actualC.LevelMapper)
	assert.ToBeSame(t, instance.DetectSkipFrames, actualC.DetectSkipFrames)
	assert.ToBeSame(t, instance, actualC.parent)
	assert.ToBeEqual(t, "foo.bar.", actualC.fieldKeyPath.prefix())
	assert.ToBeSame(t, instance.fieldKeyPath, actualC.fieldKeyPath.parent)
	assert.ToBeEqual(t, "bar", actualC.fieldKeyPath.key)
	assert.ToBeEqual(t, attrs(nil), actualC.attrs)
}

func TestHandler_WithGroup_empty(t *testing.T) {
	instance := &Handler{fieldKeyPath: newFieldKeyPath(nil, "foo")}

	actual := instance.WithGroup("")

	assert.ToBeSame(t, instance, actual)
}

func TestHandler_WithGroup_deepPath(t *testing.T) {
	const depth = 20_000
	var instance sdk.Handler = &Handler{}
	for i := 0; i < depth; i++ {
		instance = instance.WithGroup("g")
	}

	grouped := instance.(*Handler)
	expectedPrefix := strings.Repeat("g.", depth)
	assert.ToBeEqual(t, expectedPrefix, grouped.fieldKeyPath.prefix())
	assert.ToBeEqual(t, depth*2, grouped.fieldKeyPath.length)
	assert.ToBeEqual(t, "g", grouped.fieldKeyPath.key)

	withAttr := grouped.WithAttrs([]sdk.Attr{sdk.Int("value", 1)}).(*Handler)
	actualValue, actualExists := withAttr.fields().Get(expectedPrefix + "value")
	assert.ToBeEqual(t, int64(1), actualValue)
	assert.ToBeEqual(t, true, actualExists)
}

func TestHandler_WithGroup_preservesPathsAndBranches(t *testing.T) {
	base := (&Handler{}).WithAttrs([]sdk.Attr{sdk.Int("root", 1)})
	outer := base.WithGroup("outer")
	assert.ToBeSame(t, outer, outer.WithGroup(""))
	deep := outer.
		WithAttrs([]sdk.Attr{sdk.Int("middle", 2)}).
		WithGroup("inner").
		WithAttrs([]sdk.Attr{sdk.Int("deep", 3)}).(*Handler)
	sibling := outer.
		WithGroup("sibling").
		WithAttrs([]sdk.Attr{sdk.Int("value", 4)}).(*Handler)

	var actualDeepKeys []string
	actualDeepErr := deep.fields().ForEach(func(key string, _ any) error {
		actualDeepKeys = append(actualDeepKeys, key)
		return nil
	})
	assert.ToBeNoError(t, actualDeepErr)
	assert.ToBeEqual(t, []string{"outer.inner.deep", "outer.middle", "root"}, actualDeepKeys)
	actualSibling, actualSiblingErr := fields.AsMap(sibling.fields())
	assert.ToBeNoError(t, actualSiblingErr)
	assert.ToBeEqual(t, map[string]any{"outer.sibling.value": int64(4), "root": int64(1)}, actualSibling)
}

func TestHandler_WithGroup_preservesDotsInNames(t *testing.T) {
	dotted := (&Handler{}).
		WithGroup("a.b").
		WithAttrs([]sdk.Attr{sdk.Int("value", 1)}).(*Handler)
	nested := (&Handler{}).
		WithGroup("a").
		WithGroup("b").
		WithAttrs([]sdk.Attr{sdk.Int("value", 1)}).(*Handler)

	actualDotted, actualDottedErr := fields.AsMap(dotted.fields())
	actualNested, actualNestedErr := fields.AsMap(nested.fields())
	assert.ToBeNoError(t, actualDottedErr)
	assert.ToBeNoError(t, actualNestedErr)
	assert.ToBeEqual(t, actualDotted, actualNested)
}

func TestFieldKeyPath_prefix_concurrent(t *testing.T) {
	path := newFieldKeyPath(newFieldKeyPath(nil, "outer"), "inner")
	results := make(chan string, 32)
	var waitGroup sync.WaitGroup
	for i := 0; i < cap(results); i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			results <- path.prefix()
		}()
	}
	waitGroup.Wait()
	close(results)

	for actual := range results {
		assert.ToBeEqual(t, "outer.inner.", actual)
	}
}

func TestHandler_getDelegate(t *testing.T) {
	aLogger := recording.NewCoreLogger()
	cases := []struct {
		name             string
		givenDelegate    log.CoreLogger
		expectedDelegate log.CoreLogger
	}{
		{"provided", aLogger, aLogger},
		{"nil", nil, nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			instance := &Handler{Delegate: c.givenDelegate}

			actual := instance.getDelegate()
			if expected := c.expectedDelegate; expected != nil {
				assert.ToBeSame(t, expected, actual)
			} else {
				assert.ToBeOfType(t, log.GetRootLogger(), actual)
			}
		})
	}
}

func TestHandler_getLevelMapper(t *testing.T) {
	aLevelMapper := NewLevelMapperFacade(nil)
	cases := []struct {
		name                string
		givenLevelMapper    LevelMapper
		expectedLevelMapper LevelMapper
	}{
		{"provided", aLevelMapper, aLevelMapper},
		{"nil", nil, DefaultLevelMapper},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			instance := &Handler{LevelMapper: c.givenLevelMapper}

			actual := instance.getLevelMapper()
			assert.ToBeEqual(t, c.expectedLevelMapper, actual)
		})
	}
}

func TestHandler_getDetectSkipFrames(t *testing.T) {
	var aDetectSkipFrames DetectSkipFrames = func(uint16) uint16 {
		panic("this should never be called")
	}
	cases := []struct {
		name                     string
		givenDetectSkipFrames    DetectSkipFrames
		expectedDetectSkipFrames DetectSkipFrames
	}{
		{"provided", aDetectSkipFrames, aDetectSkipFrames},
		{"nil", nil, DefaultDetectSkipFrames},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			instance := &Handler{DetectSkipFrames: c.givenDetectSkipFrames}

			actual := instance.getDetectSkipFrames()
			assert.ToBeEqual(t, c.expectedDetectSkipFrames, actual)
		})
	}
}

func Test_helperOf(t *testing.T) {
	var called *bool
	cases := []struct {
		name               string
		givenCoreLogger    log.CoreLogger
		expectedToBeCalled bool
	}{
		{"with", someCoreLoggerWithHelper(func() { result := true; called = &result }), true},
		{"without", someCoreLogger{}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			initial := false
			called = &initial

			helperOf(c.givenCoreLogger)()

			assert.ToBeEqual(t, c.expectedToBeCalled, *called)
		})
	}
}

type someCoreLogger struct{}

func (instance someCoreLogger) Log(log.Event, uint16) {
	panic("should never be called")
}

func (instance someCoreLogger) IsLevelEnabled(level.Level) bool {
	panic("should never be called")
}

func (instance someCoreLogger) GetName() string {
	panic("should never be called")
}

func (instance someCoreLogger) NewEvent(level.Level, map[string]any) log.Event {
	panic("should never be called")
}

func (instance someCoreLogger) Accepts(log.Event) bool {
	panic("should never be called")
}

func (instance someCoreLogger) GetProvider() log.Provider {
	panic("should never be called")
}

type someCoreLoggerWithHelper func()

func (instance someCoreLoggerWithHelper) Helper() func() {
	return instance
}

func (instance someCoreLoggerWithHelper) Log(log.Event, uint16) {
	panic("should never be called")
}

func (instance someCoreLoggerWithHelper) IsLevelEnabled(level.Level) bool {
	panic("should never be called")
}

func (instance someCoreLoggerWithHelper) GetName() string {
	panic("should never be called")
}

func (instance someCoreLoggerWithHelper) NewEvent(level.Level, map[string]any) log.Event {
	panic("should never be called")
}

func (instance someCoreLoggerWithHelper) Accepts(log.Event) bool {
	panic("should never be called")
}

func (instance someCoreLoggerWithHelper) GetProvider() log.Provider {
	panic("should never be called")
}

type delegateCoreLoggerWithHelper struct {
	log.CoreLogger
	helper func()
}

type programCounterRejectingCoreLogger struct {
	*recording.CoreLogger
}

func (instance *programCounterRejectingCoreLogger) Accepts(event log.Event) bool {
	_, withProgramCounter := event.(interface{ GetProgramCounter() uintptr })
	return !withProgramCounter && instance.CoreLogger.Accepts(event)
}

func (instance delegateCoreLoggerWithHelper) Helper() func() {
	return instance.helper
}
