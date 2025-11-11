//go:build go1.21

package sdk

import (
	"log/slog"
	"testing"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/testing/recording"
)

func TestNew(t *testing.T) {
	aLogger := recording.NewCoreLogger()
	anotherHandler := &Handler{}

	actual := New(aLogger, func(v *Handler) {
		v.parent = anotherHandler
	}, func(v *Handler) {
		v.fieldKeyPrefix = "foo"
	})

	assert.ToBeNotNil(t, actual)
	actualHandler := actual.Handler()
	assert.ToBeNotNil(t, actualHandler)
	assert.ToBeOfType(t, (*Handler)(nil), actualHandler)
	cActualHandler := actualHandler.(*Handler)
	assert.ToBeSame(t, aLogger, cActualHandler.Delegate)
	assert.ToBeEqual(t, "foo", cActualHandler.fieldKeyPrefix)
	assert.ToBeSame(t, anotherHandler, cActualHandler.parent)
}

func TestConfigure(t *testing.T) {
	def := slog.Default()
	defer slog.SetDefault(def)

	rootLogger := log.GetRootLogger()
	anotherHandler := &Handler{}

	Configure(func(v *Handler) {
		v.parent = anotherHandler
	}, func(v *Handler) {
		v.fieldKeyPrefix = "foo"
	})

	actual := slog.Default()
	assert.ToBeNotNil(t, actual)
	actualHandler := actual.Handler()
	assert.ToBeNotNil(t, actualHandler)
	assert.ToBeOfType(t, (*Handler)(nil), actualHandler)
	cActualHandler := actualHandler.(*Handler)
	assert.ToBeOfType(t, rootLogger, cActualHandler.Delegate)
	assert.ToBeEqual(t, "foo", cActualHandler.fieldKeyPrefix)
	assert.ToBeSame(t, anotherHandler, cActualHandler.parent)
}

func TestConfigureWith(t *testing.T) {
	def := slog.Default()
	defer slog.SetDefault(def)

	aLogger := recording.NewCoreLogger()
	anotherHandler := &Handler{}

	ConfigureWith(aLogger, func(v *Handler) {
		v.parent = anotherHandler
	}, func(v *Handler) {
		v.fieldKeyPrefix = "foo"
	})

	actual := slog.Default()
	assert.ToBeNotNil(t, actual)
	actualHandler := actual.Handler()
	assert.ToBeNotNil(t, actualHandler)
	assert.ToBeOfType(t, (*Handler)(nil), actualHandler)
	cActualHandler := actualHandler.(*Handler)
	assert.ToBeSame(t, aLogger, cActualHandler.Delegate)
	assert.ToBeEqual(t, "foo", cActualHandler.fieldKeyPrefix)
	assert.ToBeSame(t, anotherHandler, cActualHandler.parent)
}
