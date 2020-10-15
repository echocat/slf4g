package color

import (
	"fmt"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_Mode_ShouldColorize(t *testing.T) {
	cases := []struct {
		instance        Mode
		expectedNone    bool
		expectedNative  bool
		expectedAssumed bool
	}{
		{ModeAuto, false, true, true},
		{ModeAlways, true, true, true},
		{ModeNever, false, false, false},
	}
	for _, c := range cases {
		t.Run(c.instance.String(), func(t *testing.T) {
			assert.ToBeEqual(t, c.expectedNone, c.instance.ShouldColorize(SupportedNone))
			assert.ToBeEqual(t, c.expectedNative, c.instance.ShouldColorize(SupportedNative))
			assert.ToBeEqual(t, c.expectedAssumed, c.instance.ShouldColorize(SupportedAssumed))
		})
	}
}

func Test_Mode_MarshalText(t *testing.T) {
	cases := []struct {
		expected string
		instance Mode
	}{
		{"auto", ModeAuto},
		{"always", ModeAlways},
		{"never", ModeNever},
	}
	for _, c := range cases {
		t.Run(c.expected, func(t *testing.T) {
			actual, actualErr := c.instance.MarshalText()
			assert.ToBeNil(t, actualErr)
			assert.ToBeEqual(t, c.expected, string(actual))
		})
	}
}

func Test_Mode_MarshalText_errors(t *testing.T) {
	instance := Mode(66)

	actual, actualErr := instance.MarshalText()
	assert.ToBeEqual(t, fmt.Errorf("%w: 66", ErrIllegalMode), actualErr)
	assert.ToBeNil(t, actual)
}

func Test_Mode_UnmarshalText(t *testing.T) {
	cases := []struct {
		given    string
		expected Mode
	}{
		{"auto", ModeAuto},
		{"automatic", ModeAuto},
		{"detect", ModeAuto},
		{"always", ModeAlways},
		{"on", ModeAlways},
		{"true", ModeAlways},
		{"1", ModeAlways},
		{"never", ModeNever},
		{"off", ModeNever},
		{"false", ModeNever},
		{"0", ModeNever},
	}
	for _, c := range cases {
		t.Run(c.given, func(t *testing.T) {
			var instance Mode

			actualErr := instance.UnmarshalText([]byte(c.given))

			assert.ToBeNil(t, actualErr)
			assert.ToBeEqual(t, c.expected, instance)
		})
	}
}

func Test_Mode_UnmarshalText_errors(t *testing.T) {
	var instance Mode

	actualErr := instance.UnmarshalText([]byte("foo"))

	assert.ToBeEqual(t, fmt.Errorf("%w: foo", ErrIllegalMode), actualErr)
}

func Test_Mode_String(t *testing.T) {
	cases := []struct {
		expected string
		instance Mode
	}{
		{"auto", ModeAuto},
		{"always", ModeAlways},
		{"never", ModeNever},
		{"illegal-color-mode-66", Mode(66)},
	}
	for _, c := range cases {
		t.Run(c.expected, func(t *testing.T) {
			actual := c.instance.String()

			assert.ToBeEqual(t, c.expected, actual)
		})
	}
}

func Test_Mode_Set(t *testing.T) {
	cases := []struct {
		given    string
		expected Mode
	}{
		{"auto", ModeAuto},
		{"automatic", ModeAuto},
		{"detect", ModeAuto},
		{"always", ModeAlways},
		{"on", ModeAlways},
		{"true", ModeAlways},
		{"1", ModeAlways},
		{"never", ModeNever},
		{"off", ModeNever},
		{"false", ModeNever},
		{"0", ModeNever},
	}
	for _, c := range cases {
		t.Run(c.given, func(t *testing.T) {
			var instance Mode

			actualErr := instance.Set(c.given)

			assert.ToBeNil(t, actualErr)
			assert.ToBeEqual(t, c.expected, instance)
		})
	}
}

func Test_Mode_Set_errors(t *testing.T) {
	var instance Mode

	actualErr := instance.Set("foo")

	assert.ToBeEqual(t, fmt.Errorf("%w: foo", ErrIllegalMode), actualErr)
}

func Test_Modes_Strings(t *testing.T) {
	assert.ToBeEqual(t, []string{"always", "auto"}, Modes{ModeAlways, ModeAuto}.Strings())
	assert.ToBeEqual(t, []string{"always"}, Modes{ModeAlways}.Strings())
	assert.ToBeEqual(t, []string{}, Modes{}.Strings())
}

func Test_Modes_String(t *testing.T) {
	assert.ToBeEqual(t, "always,auto", Modes{ModeAlways, ModeAuto}.String())
	assert.ToBeEqual(t, "always", Modes{ModeAlways}.String())
	assert.ToBeEqual(t, "", Modes{}.String())
}

func Test_AllModes(t *testing.T) {
	actual := AllModes()

	assert.ToBeEqual(t, Modes{ModeAuto, ModeAlways, ModeNever}, actual)
}
