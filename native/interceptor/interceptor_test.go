package interceptor

import (
	"math"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/testing/recording"

	"github.com/echocat/slf4g/internal/test/assert"

	log "github.com/echocat/slf4g"
)

func Test_Interceptors_Add(t *testing.T) {
	given := []Interceptor{
		noopInterceptorButSorted(66),
		noopInterceptorButSorted(-3),
		noopInterceptorButSorted(0),
		noopInterceptorButSorted(1),
		noopInterceptorButSorted(-300),
		noopInterceptorButSorted(444),
		noopInterceptorButSorted(3345),
		noopInterceptorButSorted(-30),
	}
	expected := []Interceptor{
		noopInterceptorButSorted(-300),
		noopInterceptorButSorted(-30),
		noopInterceptorButSorted(-3),
		noopInterceptorButSorted(0),
		noopInterceptorButSorted(1),
		noopInterceptorButSorted(66),
		noopInterceptorButSorted(444),
		noopInterceptorButSorted(3345),
	}
	instance := &Interceptors{}

	for _, g := range given {
		actual := instance.Add(g)
		assert.ToBeSame(t, instance, actual)
	}

	assert.ToBeEqual(t, expected, []Interceptor(*instance))
}

func Test_Interceptors_Add_doesNotModifyExistingSnapshot(t *testing.T) {
	instance := make(Interceptors, 2, 3)
	instance[0] = noopInterceptorButSorted(1)
	instance[1] = noopInterceptorButSorted(2)
	snapshot := instance

	instance.Add(noopInterceptorButSorted(0))

	assert.ToBeEqual(t, []Interceptor{
		noopInterceptorButSorted(1),
		noopInterceptorButSorted(2),
	}, []Interceptor(snapshot))
	assert.ToBeEqual(t, []Interceptor{
		noopInterceptorButSorted(0),
		noopInterceptorButSorted(1),
		noopInterceptorButSorted(2),
	}, []Interceptor(instance))
}

func Test_Interceptors_Add_concurrentlyWithSnapshot(t *testing.T) {
	const count = 100
	var instance Interceptors
	start := make(chan struct{})
	readDone := make(chan struct{})
	var additions sync.WaitGroup
	additions.Add(count)

	for i := 0; i < count; i++ {
		priority := int16(i)
		go func() {
			defer additions.Done()
			<-start
			instance.Add(noopInterceptorButSorted(priority))
		}()
	}

	var reader sync.WaitGroup
	reader.Add(1)
	go func() {
		defer reader.Done()
		<-start
		for {
			select {
			case <-readDone:
				return
			default:
				if snapshot := instance.Snapshot(); !sort.IsSorted(snapshot) {
					t.Errorf("observed unsorted interceptor snapshot: %v", snapshot)
					return
				}
			}
		}
	}()

	close(start)
	additions.Wait()
	close(readDone)
	reader.Wait()

	actual := instance.Snapshot()
	assert.ToBeEqual(t, count, len(actual))
	for i, interceptor := range actual {
		assert.ToBeEqual(t, int16(i), interceptor.GetPriority())
	}
}

func Test_Interceptors_Add_allowsAddToDifferentCollectionWhileResolvingPriority(t *testing.T) {
	var nested Interceptors
	nestedDone := make(chan struct{})
	var once sync.Once
	reentrant := interceptorWithPriorityFunc(func() int16 {
		once.Do(func() {
			go func() {
				nested.Add(Noop())
				close(nestedDone)
			}()
		})
		select {
		case <-nestedDone:
		case <-time.After(time.Second):
			t.Fatal("nested Add blocked while resolving interceptor priority")
		}
		return 0
	})
	instance := Interceptors{noopInterceptorButSorted(1)}

	instance.Add(reentrant)

	assert.ToBeEqual(t, 1, len(nested.Snapshot()))
}

func Test_Interceptors_OnBeforeLog(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEventWithFields(level.Warn, fields.With("foo", "bar"))

	instance := make(Interceptors, 10)
	for i := range len(instance) {
		instance[i] = newOnBeforeLogCalledInterceptor()
	}

	actual := instance.OnBeforeLog(givenEvent, givenLogger.Provider)

	assert.ToBeSame(t, givenEvent, actual)

	for _, it := range instance {
		assert.ToBeEqual(t, onBeforeLogCalledInterceptor(1), *(it.(*onBeforeLogCalledInterceptor)))
	}
}

func Test_Interceptors_OnBeforeLog_withNilEvent(t *testing.T) {
	givenLogger := recording.NewLogger()

	instance := make(Interceptors, 10)
	for i := range len(instance) {
		instance[i] = newOnBeforeLogCalledInterceptor()
	}

	actual := instance.OnBeforeLog(nil, givenLogger.Provider)

	assert.ToBeNil(t, actual)

	for _, it := range instance {
		assert.ToBeEqual(t, onBeforeLogCalledInterceptor(0), *(it.(*onBeforeLogCalledInterceptor)))
	}
}

func Test_Interceptors_OnBeforeLog_oneReturnsNil(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEventWithFields(level.Warn, fields.With("foo", "bar"))

	instance := make(Interceptors, 10)
	for i := range len(instance) {
		if i == 2 {
			instance[2] = OnBeforeLogFunc(func(actualEvent log.Event, _ log.Provider) log.Event {
				assert.ToBeSame(t, givenEvent, actualEvent)
				return nil
			})
		} else {
			instance[i] = newOnBeforeLogCalledInterceptor()
		}
	}

	actual := instance.OnBeforeLog(givenEvent, givenLogger.Provider)

	assert.ToBeNil(t, actual)

	for i, it := range instance {
		if i < 2 {
			assert.ToBeEqual(t, onBeforeLogCalledInterceptor(1), *(it.(*onBeforeLogCalledInterceptor)))
		} else if i > 2 {
			assert.ToBeEqual(t, onBeforeLogCalledInterceptor(0), *(it.(*onBeforeLogCalledInterceptor)))
		}
	}
}

func Test_Interceptors_OnAfterLog(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEventWithFields(level.Warn, fields.With("foo", "bar"))

	instance := make(Interceptors, 10)
	for i := range len(instance) {
		instance[i] = newOnAfterLogCalledInterceptor(t, givenEvent, true)
	}

	actual := instance.OnAfterLog(givenEvent, givenLogger.Provider)

	assert.ToBeEqual(t, true, actual)

	for _, it := range instance {
		assert.ToBeEqual(t, 1, it.(*onAfterLogCalledInterceptor).calledAmount)
	}
}

func Test_Interceptors_OnAfterLog_withStop(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEventWithFields(level.Warn, fields.With("foo", "bar"))

	instance := make(Interceptors, 10)
	for i := range len(instance) {
		instance[i] = newOnAfterLogCalledInterceptor(t, givenEvent, i < len(instance)-2)
	}

	actual := instance.OnAfterLog(givenEvent, givenLogger.Provider)

	assert.ToBeEqual(t, false, actual)

	for i, it := range instance {
		if i <= len(instance)-2 {
			assert.ToBeEqual(t, 1, it.(*onAfterLogCalledInterceptor).calledAmount)
		} else {
			assert.ToBeEqual(t, 0, it.(*onAfterLogCalledInterceptor).calledAmount)
		}
	}
}

func Test_Interceptors_OnAfterLog_withNilEvent(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEventWithFields(level.Warn, fields.With("foo", "bar"))

	instance := make(Interceptors, 10)
	for i := range len(instance) {
		instance[i] = newOnAfterLogCalledInterceptor(t, givenEvent, true)
	}

	actual := instance.OnAfterLog(nil, givenLogger.Provider)

	assert.ToBeEqual(t, false, actual)

	for _, it := range instance {
		assert.ToBeEqual(t, 0, it.(*onAfterLogCalledInterceptor).calledAmount)
	}
}

func Test_Interceptors_Len(t *testing.T) {
	instance := Interceptors{
		noopInterceptorButSorted(66),
		noopInterceptorButSorted(-3),
		noopInterceptorButSorted(0),
		noopInterceptorButSorted(1),
		noopInterceptorButSorted(-300),
		noopInterceptorButSorted(444),
		noopInterceptorButSorted(3345),
		noopInterceptorButSorted(-30),
	}

	actual := instance.Len()

	assert.ToBeEqual(t, 8, actual)
}

func Test_Interceptors_GetPriority(t *testing.T) {
	instance := Interceptors{
		noopInterceptorButSorted(-3),
		noopInterceptorButSorted(0),
		noopInterceptorButSorted(1),
		noopInterceptorButSorted(-300),
		noopInterceptorButSorted(444),
		noopInterceptorButSorted(3345),
		noopInterceptorButSorted(-30),
	}

	actual := instance.GetPriority()

	assert.ToBeEqual(t, int16(-300), actual)
}

func Test_OnBeforeLogFunc_OnBeforeLog(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEventWithFields(level.Warn, fields.With("foo", "bar"))
	expectedEvent := givenLogger.NewEventWithFields(level.Warn, fields.With("foo", "bar2"))

	instance := OnBeforeLogFunc(func(actualEvent log.Event, actualProvider log.Provider) (intercepted log.Event) {
		assert.ToBeSame(t, givenEvent, actualEvent)
		assert.ToBeSame(t, givenLogger.Provider, actualProvider)
		return expectedEvent
	})

	actual := instance.OnBeforeLog(givenEvent, givenLogger.Provider)

	assert.ToBeSame(t, expectedEvent, actual)
}

func Test_OnBeforeLogFunc_OnAfterLog(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEventWithFields(level.Warn, fields.With("foo", "bar"))

	instance := OnBeforeLogFunc(func(event log.Event, _ log.Provider) log.Event {
		assert.Fail(t, "Should not be called; but was.")
		return event
	})

	actual := instance.OnAfterLog(givenEvent, givenLogger.Provider)

	assert.ToBeEqual(t, true, actual)
}

func Test_OnBeforeLogFunc_GetPriority(t *testing.T) {
	instance := OnBeforeLogFunc(func(event log.Event, _ log.Provider) log.Event {
		assert.Fail(t, "Should not be called; but was.")
		return event
	})

	actual := instance.GetPriority()

	assert.ToBeEqual(t, int16(0), actual)
}

func Test_OnAfterLogFunc_OnBeforeLog(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEventWithFields(level.Warn, fields.With("foo", "bar"))

	instance := OnAfterLogFunc(func(actualEvent log.Event, actualProvider log.Provider) bool {
		assert.Fail(t, "Should not be called; but was.")
		return false
	})

	actual := instance.OnBeforeLog(givenEvent, givenLogger.Provider)

	assert.ToBeSame(t, givenEvent, actual)
}

func Test_OnAfterLogFunc_OnAfterLog(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEventWithFields(level.Warn, fields.With("foo", "bar"))

	instance := OnAfterLogFunc(func(actualEvent log.Event, actualProvider log.Provider) bool {
		assert.ToBeSame(t, givenEvent, actualEvent)
		assert.ToBeSame(t, givenLogger.Provider, actualProvider)
		return true
	})

	actual := instance.OnAfterLog(givenEvent, givenLogger.Provider)

	assert.ToBeEqual(t, true, actual)
}

func Test_OnAfterLogFunc_GetPriority(t *testing.T) {
	instance := OnAfterLogFunc(func(event log.Event, _ log.Provider) bool {
		assert.Fail(t, "Should not be called; but was.")
		return false
	})

	actual := instance.GetPriority()

	assert.ToBeEqual(t, int16(0), actual)
}

func Test_NewFacade(t *testing.T) {
	provider := func() Interceptor { return nil }

	actual := NewFacade(provider)

	assert.ToBeSame(t, facade(provider), actual)
}

func Test_facade_OnBeforeLog(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEventWithFields(level.Warn, fields.With("foo", "bar"))
	expectedEvent := givenLogger.NewEventWithFields(level.Warn, fields.With("foo", "bar2"))
	givenDelegate := OnBeforeLogFunc(func(actualEvent log.Event, actualProvider log.Provider) log.Event {
		assert.ToBeSame(t, givenEvent, actualEvent)
		assert.ToBeSame(t, givenLogger.Provider, actualProvider)
		return expectedEvent
	})
	instance := facade(func() Interceptor {
		return givenDelegate
	})

	actual := instance.OnBeforeLog(givenEvent, givenLogger.Provider)

	assert.ToBeSame(t, expectedEvent, actual)
}

func Test_facade_OnAfterLog(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEventWithFields(level.Warn, fields.With("foo", "bar"))
	givenDelegate := OnAfterLogFunc(func(actualEvent log.Event, actualProvider log.Provider) bool {
		assert.ToBeSame(t, givenEvent, actualEvent)
		assert.ToBeSame(t, givenLogger.Provider, actualProvider)
		return true
	})
	instance := facade(func() Interceptor {
		return givenDelegate
	})

	actual := instance.OnAfterLog(givenEvent, givenLogger.Provider)

	assert.ToBeEqual(t, true, actual)
}

func Test_facade_GetPriority(t *testing.T) {
	givenDelegate := noopInterceptorButSorted(666)
	instance := facade(func() Interceptor {
		return givenDelegate
	})

	actual := instance.GetPriority()

	assert.ToBeEqual(t, int16(666), actual)
}

func Test_Noop(t *testing.T) {
	actual := Noop()

	assert.ToBeSame(t, noopV, actual)
}

func Test_noop_OnBeforeLog(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEventWithFields(level.Warn, fields.With("foo", "bar"))
	instance := &noop{}

	actual := instance.OnBeforeLog(givenEvent, givenLogger.Provider)

	assert.ToBeSame(t, givenEvent, actual)
}

func Test_noop_OnAfterLog(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEventWithFields(level.Warn, fields.With("foo", "bar"))
	instance := &noop{}

	actual := instance.OnAfterLog(givenEvent, givenLogger.Provider)

	assert.ToBeEqual(t, true, actual)
}

func Test_noop_GetPriority(t *testing.T) {
	instance := &noop{}

	actual := instance.GetPriority()

	assert.ToBeEqual(t, int16(math.MaxInt16), actual)
}

type noopInterceptorButSorted int16

type interceptorWithPriorityFunc func() int16

func (instance interceptorWithPriorityFunc) OnBeforeLog(event log.Event, _ log.Provider) log.Event {
	return event
}

func (instance interceptorWithPriorityFunc) OnAfterLog(log.Event, log.Provider) bool {
	return true
}

func (instance interceptorWithPriorityFunc) GetPriority() int16 {
	return instance()
}

func (instance noopInterceptorButSorted) OnBeforeLog(log.Event, log.Provider) log.Event {
	panic("not implemented")
}

func (instance noopInterceptorButSorted) OnAfterLog(log.Event, log.Provider) bool {
	panic("not implemented")
}

func (instance noopInterceptorButSorted) GetPriority() int16 {
	return int16(instance)
}

func newOnBeforeLogCalledInterceptor() *onBeforeLogCalledInterceptor {
	result := onBeforeLogCalledInterceptor(0)
	return &result
}

type onBeforeLogCalledInterceptor int

func (instance *onBeforeLogCalledInterceptor) OnBeforeLog(event log.Event, _ log.Provider) log.Event {
	*instance++
	return event
}

func (instance *onBeforeLogCalledInterceptor) OnAfterLog(log.Event, log.Provider) bool {
	panic("not implemented")
}

func (instance *onBeforeLogCalledInterceptor) GetPriority() int16 {
	panic("not implemented")
}

func newOnAfterLogCalledInterceptor(t *testing.T, expectedEvent log.Event, canContinue bool) *onAfterLogCalledInterceptor {
	return &onAfterLogCalledInterceptor{
		calledAmount:  0,
		expectedEvent: expectedEvent,
		t:             t,
		canContinue:   canContinue,
	}
}

type onAfterLogCalledInterceptor struct {
	calledAmount  int
	expectedEvent log.Event
	t             *testing.T
	canContinue   bool
}

func (instance onAfterLogCalledInterceptor) OnBeforeLog(log.Event, log.Provider) log.Event {
	panic("not implemented")
}

func (instance *onAfterLogCalledInterceptor) OnAfterLog(givenEvent log.Event, _ log.Provider) (canContinue bool) {
	instance.calledAmount++
	assert.ToBeSame(instance.t, instance.expectedEvent, givenEvent)
	return instance.canContinue
}

func (instance onAfterLogCalledInterceptor) GetPriority() int16 {
	panic("not implemented")
}
