package globalprovider

import (
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func TestPush_restoresPreviousValue(t *testing.T) {
	defer Set(nil)
	previous := &struct{ name string }{"previous"}
	installed := &struct{ name string }{"installed"}
	Set(previous)

	cleanup := Push(installed)
	assert.ToBeSame(t, installed, Get())

	cleanup()
	assert.ToBeSame(t, previous, Get())
}

func TestPush_skipsInactiveScopes(t *testing.T) {
	defer Set(nil)
	previous := &struct{ name string }{"previous"}
	first := &struct{ name string }{"first"}
	second := &struct{ name string }{"second"}
	third := &struct{ name string }{"third"}
	Set(previous)
	cleanupFirst := Push(first)
	cleanupSecond := Push(second)
	cleanupThird := Push(third)

	cleanupSecond()
	cleanupFirst()
	assert.ToBeSame(t, third, Get())

	cleanupThird()
	assert.ToBeSame(t, previous, Get())
}

func TestPush_doesNotOverwriteLaterSet(t *testing.T) {
	defer Set(nil)
	previous := &struct{ name string }{"previous"}
	installed := &struct{ name string }{"installed"}
	later := &struct{ name string }{"later"}
	Set(previous)
	cleanup := Push(installed)

	Set(later)
	cleanup()

	assert.ToBeSame(t, later, Get())
}

func TestPush_distinguishesLaterInstallationOfSameValue(t *testing.T) {
	defer Set(nil)
	previous := &struct{ name string }{"previous"}
	installed := &struct{ name string }{"installed"}
	Set(previous)
	cleanup := Push(installed)

	Set(installed)
	cleanup()

	assert.ToBeSame(t, installed, Get())
}

func TestPush_distinguishesScopedInstallationsOfSameValue(t *testing.T) {
	defer Set(nil)
	previous := &struct{ name string }{"previous"}
	installed := &struct{ name string }{"installed"}
	Set(previous)
	cleanupFirst := Push(installed)
	cleanupSecond := Push(installed)

	cleanupFirst()
	assert.ToBeSame(t, installed, Get())

	cleanupSecond()
	assert.ToBeSame(t, previous, Get())
}

func TestPush_cleanupIsIdempotent(t *testing.T) {
	defer Set(nil)
	previous := &struct{ name string }{"previous"}
	installed := &struct{ name string }{"installed"}
	Set(previous)
	cleanup := Push(installed)

	cleanup()
	cleanup()

	assert.ToBeSame(t, previous, Get())
}
