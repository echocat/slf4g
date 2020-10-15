package color

import (
	"fmt"
	"os"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_DetectSupportForWriter_assuming(t *testing.T) {
	old := SupportAssumptionDetections
	defer func() {
		SupportAssumptionDetections = old
	}()

	var expectedToBeSupported bool
	SupportAssumptionDetections = []SupportAssumptionDetection{func() bool {
		return expectedToBeSupported
	}}

	expectedToBeSupported = true
	actualWriter1, actualSupported1, actualErr1 := DetectSupportForWriter(os.Stdout)
	assert.ToBeNil(t, actualErr1)
	assert.ToBeSame(t, os.Stdout, actualWriter1)
	assert.ToBeEqual(t, SupportedAssumed, actualSupported1)

	expectedToBeSupported = false
	actualWriter2, actualSupported2, actualErr2 := DetectSupportForWriter(os.Stdout)
	assert.ToBeNil(t, actualErr2)
	assert.ToBeSame(t, os.Stdout, actualWriter2)
	assert.ToBeEqual(t, SupportedNone, actualSupported2)
}

func Test_Supported_IsSupported(t *testing.T) {
	assert.ToBeEqual(t, false, SupportedNone.IsSupported())
	assert.ToBeEqual(t, true, SupportedNative.IsSupported())
	assert.ToBeEqual(t, true, SupportedAssumed.IsSupported())
}

func Test_Supported_MarshalText(t *testing.T) {
	cases := []struct {
		expected string
		instance Supported
	}{
		{"none", SupportedNone},
		{"native", SupportedNative},
		{"assumed", SupportedAssumed},
	}
	for _, c := range cases {
		t.Run(c.expected, func(t *testing.T) {
			actual, actualErr := c.instance.MarshalText()
			assert.ToBeNil(t, actualErr)
			assert.ToBeEqual(t, c.expected, string(actual))
		})
	}
}

func Test_Supported_MarshalText_errors(t *testing.T) {
	instance := Supported(66)

	actual, actualErr := instance.MarshalText()
	assert.ToBeEqual(t, fmt.Errorf("%w: 66", ErrIllegalSupport), actualErr)
	assert.ToBeNil(t, actual)
}

func Test_Supported_UnmarshalText(t *testing.T) {
	cases := []struct {
		given    string
		expected Supported
	}{
		{"none", SupportedNone},
		{"no", SupportedNone},
		{"0", SupportedNone},
		{"never", SupportedNone},
		{"off", SupportedNone},
		{"false", SupportedNone},
		{"native", SupportedNative},
		{"assumed", SupportedAssumed},
	}
	for _, c := range cases {
		t.Run(c.given, func(t *testing.T) {
			var instance Supported

			actualErr := instance.UnmarshalText([]byte(c.given))

			assert.ToBeNil(t, actualErr)
			assert.ToBeEqual(t, c.expected, instance)
		})
	}
}

func Test_Supported_UnmarshalText_errors(t *testing.T) {
	var instance Supported

	actualErr := instance.UnmarshalText([]byte("foo"))

	assert.ToBeEqual(t, fmt.Errorf("%w: foo", ErrIllegalSupport), actualErr)
}

func Test_Supported_String(t *testing.T) {
	cases := []struct {
		expected string
		instance Supported
	}{
		{"none", SupportedNone},
		{"native", SupportedNative},
		{"assumed", SupportedAssumed},
		{"illegal-color-support-66", Supported(66)},
	}
	for _, c := range cases {
		t.Run(c.expected, func(t *testing.T) {
			actual := c.instance.String()

			assert.ToBeEqual(t, c.expected, actual)
		})
	}
}

func Test_Supported_Set(t *testing.T) {
	cases := []struct {
		given    string
		expected Supported
	}{
		{"none", SupportedNone},
		{"no", SupportedNone},
		{"0", SupportedNone},
		{"never", SupportedNone},
		{"off", SupportedNone},
		{"false", SupportedNone},
		{"native", SupportedNative},
		{"assumed", SupportedAssumed},
	}
	for _, c := range cases {
		t.Run(c.given, func(t *testing.T) {
			var instance Supported

			actualErr := instance.Set(c.given)

			assert.ToBeNil(t, actualErr)
			assert.ToBeEqual(t, c.expected, instance)
		})
	}
}

func Test_Supported_Set_errors(t *testing.T) {
	var instance Supported

	actualErr := instance.Set("foo")

	assert.ToBeEqual(t, fmt.Errorf("%w: foo", ErrIllegalSupport), actualErr)
}

func Test_Supports_Strings(t *testing.T) {
	assert.ToBeEqual(t, []string{"native", "none"}, Supports{SupportedNative, SupportedNone}.Strings())
	assert.ToBeEqual(t, []string{"native"}, Supports{SupportedNative}.Strings())
	assert.ToBeEqual(t, []string{}, Supports{}.Strings())
}

func Test_Supports_String(t *testing.T) {
	assert.ToBeEqual(t, "native,none", Supports{SupportedNative, SupportedNone}.String())
	assert.ToBeEqual(t, "native", Supports{SupportedNative}.String())
	assert.ToBeEqual(t, "", Supports{}.String())
}

func Test_AllSupports(t *testing.T) {
	actual := AllSupports()

	assert.ToBeEqual(t, Supports{SupportedNone, SupportedNative, SupportedAssumed}, actual)
}
