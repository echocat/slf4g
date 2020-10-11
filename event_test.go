package log

import (
	"testing"

	"github.com/echocat/slf4g/fields"

	"github.com/echocat/slf4g/internal/test/assert"

	"github.com/echocat/slf4g/level"
)

func Test_NewEvent_withoutFields(t *testing.T) {
	givenProvider := &testProvider{"test"}
	givenLevel := level.Error
	givenCallDepth := 66

	actual := NewEvent(givenProvider, givenLevel, givenCallDepth)

	assert.ToBeOfType(t, &eventImpl{}, actual)
	assert.ToBeSame(t, givenProvider, actual.(*eventImpl).provider)
	assert.ToBeEqual(t, givenLevel, actual.GetLevel())
	assert.ToBeEqual(t, givenCallDepth, actual.GetCallDepth())
	assert.ToBeEqual(t, fields.Empty(), actual.(*eventImpl).fields)
	assert.ToBeNil(t, actual.GetContext())
}

func Test_NewEvent_witOneFields(t *testing.T) {
	givenProvider := &testProvider{"test"}
	givenLevel := level.Error
	givenCallDepth := 66

	givenFields1 := fields.With("a", "1")

	actual := NewEvent(givenProvider, givenLevel, givenCallDepth, givenFields1)

	assert.ToBeOfType(t, &eventImpl{}, actual)
	assert.ToBeSame(t, givenProvider, actual.(*eventImpl).provider)
	assert.ToBeEqual(t, givenLevel, actual.GetLevel())
	assert.ToBeEqual(t, givenCallDepth, actual.GetCallDepth())
	assert.ToBeSame(t, givenFields1, actual.(*eventImpl).fields)
	assert.ToBeNil(t, actual.GetContext())
}

func Test_NewEvent_wit3Fields(t *testing.T) {
	givenProvider := &testProvider{"test"}
	givenLevel := level.Error
	givenCallDepth := 66

	givenFields1 := fields.With("a", 1)
	givenFields2 := fields.With("a", 2).With("b", 2)
	givenFields3 := fields.With("a", 3).With("c", 3)

	actual := NewEvent(givenProvider, givenLevel, givenCallDepth, givenFields1, givenFields2, givenFields3)

	assert.ToBeOfType(t, &eventImpl{}, actual)
	assert.ToBeSame(t, givenProvider, actual.(*eventImpl).provider)
	assert.ToBeEqual(t, givenLevel, actual.GetLevel())
	assert.ToBeEqual(t, givenCallDepth, actual.GetCallDepth())
	equal, equalErr := fields.IsEqual(fields.With("a", 3).With("b", 2).With("c", 3), actual.(*eventImpl).fields)
	assert.ToBeNil(t, equalErr)
	assert.ToBeEqual(t, true, equal)
	assert.ToBeNil(t, actual.GetContext())
}

//type testEvent struct {
//}
//
//func (instance *testEvent) GetLevel() level.Level {
//	panic("not implemented in tests")
//}
//
//func (instance *testEvent) GetCallDepth() int {
//	panic("not implemented in tests")
//}
//
//func (instance *testEvent) GetContext() interface{} {
//	panic("not implemented in tests")
//}
//
//func (instance *testEvent) ForEach(func(key string, value interface{}) error) error {
//	panic("not implemented in tests")
//}
//
//func (instance *testEvent) Get(string) (interface{}, bool) {
//	panic("not implemented in tests")
//}
//
//func (instance *testEvent) Len() int {
//	panic("not implemented in tests")
//}
//
//func (instance *testEvent) With(string, interface{}) Event {
//	panic("not implemented in tests")
//}
//
//func (instance *testEvent) Withf(string, string, ...interface{}) Event {
//	panic("not implemented in tests")
//}
//
//func (instance *testEvent) WithError(error) Event {
//	panic("not implemented in tests")
//}
//
//func (instance *testEvent) WithAll(map[string]interface{}) Event {
//	panic("not implemented in tests")
//}
//
//func (instance *testEvent) Without(...string) Event {
//	panic("not implemented in tests")
//}
//
//func (instance *testEvent) WithCallDepth(int) Event {
//	panic("not implemented in tests")
//}
//
//func (instance *testEvent) WithContext(interface{}) Event {
//	panic("not implemented in tests")
//}
