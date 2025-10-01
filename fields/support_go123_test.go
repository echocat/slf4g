//go:build go1.23
// +build go1.23

package fields

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func ExampleIter() {
	fields := With("bar", 2).With("foo", 1)

	for f := range Iter(fields) {
		fmt.Println(f)
	}

	// Output:
	// foo=1
	// bar=2
}

func ExampleIter_withErrors() {
	fields := With("bar", 2).With("foo", 1)

	for f := range Iter(fields) {
		fmt.Println(f)
	}

	// Output:
	// foo=1
	// bar=2
}

func TestIter_regular(t *testing.T) {
	fe := ForEachFunc(func(consumer func(key string, value interface{}) error) error {
		if err := consumer("foo", 1); err != nil {
			return err
		}
		return errors.New("expected")
	})

	for f, err := range Iter(fe) {
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(f)
		}
	}

	// Output:
	// foo=1
	// expected
}

func TestIter_err(t *testing.T) {
	testError := errors.New("expected")
	fe := ForEachFunc(func(consumer func(key string, value interface{}) error) error {
		assert.ToBeNoError(t, consumer("foo", 1))
		assert.ToBeNoError(t, consumer("bar", 2))
		return testError
	})

	instance := Iter(fe)
	i := 0
	instance(func(f Field, err error) bool {
		switch i {
		case 0:
			assert.ToBeNoError(t, err)
			assert.ToBeEqual(t, f.Key(), "foo")
			assert.ToBeEqual(t, f.Value(), 1)
		case 1:
			assert.ToBeNoError(t, err)
			assert.ToBeEqual(t, f.Key(), "bar")
			assert.ToBeEqual(t, f.Value(), 2)
		case 2:
			assert.ToBeEqual(t, testError, err)
			assert.ToBeNil(t, f)
		default:
			assert.Failf(t, "unexpected iteration: %d", i)
		}
		i++
		return true
	})
	assert.ToBeEqual(t, 3, i)
}

func TestCollect(t *testing.T) {
	actual := Collect(slices.Values([]Field{
		NewField("foo", 1),
		NewField("bar", 2),
	}))

	i := 0
	assert.ToBeNoError(t, actual.ForEach(func(key string, value interface{}) error {
		i++
		if i > 2 {
			assert.Failf(t, "expected to be called with maximal of 2 elements; but was %d times", i)
		}
		switch key {
		case "foo":
			assert.ToBeEqual(t, value, 1)
		case "bar":
			assert.ToBeEqual(t, value, 2)
		default:
			assert.Failf(t, "unexpected key: %s", key)
		}
		return nil
	}))
}

func TestCollectKeyValue(t *testing.T) {
	actual := CollectKeyValue(maps.All(map[string]interface{}{
		"foo": 1,
		"bar": 2,
	}))

	i := 0
	assert.ToBeNoError(t, actual.ForEach(func(key string, value interface{}) error {
		i++
		if i > 2 {
			assert.Failf(t, "expected to be called with maximal of 2 elements; but was %d times", i)
		}
		switch key {
		case "foo":
			assert.ToBeEqual(t, value, 1)
		case "bar":
			assert.ToBeEqual(t, value, 2)
		default:
			assert.Failf(t, "unexpected key: %s", key)
		}
		return nil
	}))
}

func TestCollectErr(t *testing.T) {
	testError := errors.New("expected")
	actual, actualErr := CollectErr(func(yield func(Field, error) bool) {
		assert.ToBeEqual(t, true, yield(NewField("foo", 1), nil))
		assert.ToBeEqual(t, false, yield(nil, testError))
	})
	assert.ToBeEqual(t, testError, actualErr)
	assert.ToBeNil(t, actual)
}
