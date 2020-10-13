package fields

import (
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_KeysSpecImpl_empty(t *testing.T) {
	instance := KeysSpecImpl{}
	fields := []struct {
		nameAndDefault string
		getter         func() string
		field          *string
	}{
		{"timestamp", instance.GetTimestamp, &instance.Timestamp},
		{"message", instance.GetMessage, &instance.Message},
		{"logger", instance.GetLogger, &instance.Logger},
		{"error", instance.GetError, &instance.Error},
	}

	for _, field := range fields {
		t.Run(field.nameAndDefault, func(t *testing.T) {
			actual := field.getter()

			assert.ToBeEqual(t, field.nameAndDefault, actual)
		})
	}
}

func Test_KeysSpecImpl_configured(t *testing.T) {
	instance := KeysSpecImpl{
		Timestamp: "aTimestamp",
		Message:   "aMessage",
		Logger:    "aLogger",
		Error:     "anError",
	}
	fields := []struct {
		name     string
		getter   func() string
		expected string
	}{
		{"timestamp", instance.GetTimestamp, instance.Timestamp},
		{"message", instance.GetMessage, instance.Message},
		{"logger", instance.GetLogger, instance.Logger},
		{"error", instance.GetError, instance.Error},
	}

	for _, field := range fields {
		t.Run(field.name, func(t *testing.T) {
			actual := field.getter()

			assert.ToBeEqual(t, field.expected, actual)
		})
	}
}
