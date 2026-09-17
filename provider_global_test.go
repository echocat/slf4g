package log

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_GetProvider_returnsFacade(t *testing.T) {
	defer resetGlobal()

	actual := GetProvider()

	assert.ToBeOfType(t, providerFacade(nil), actual)
	assert.ToBeEqual(t, providerFacade(getProvider), actual)
}

func Test_getProvider_returnsDefaultIfNoOtherIsRegistered(t *testing.T) {
	defer resetGlobal()

	// Nothing cached
	assert.ToBeNil(t, getCurrentProvider())

	actual := getProvider()

	assert.ToBeSame(t, fallbackProviderV, actual)
	// Also now cached
	assert.ToBeSame(t, fallbackProviderV, getCurrentProvider())
}

func Test_getProvider_returnsRegistered(t *testing.T) {
	defer resetGlobal()
	instance := newMockProvider("instance")
	RegisterProvider(instance)

	// Nothing cached
	assert.ToBeNil(t, getCurrentProvider())

	actual := getProvider()

	assert.ToBeSame(t, instance, actual)
	// Also now cached
	assert.ToBeSame(t, instance, getCurrentProvider())
}

func Test_getProvider_returnsCachedRegistered(t *testing.T) {
	defer resetGlobal()
	instance := newMockProvider("instance")

	// Nothing cached
	assert.ToBeNil(t, getCurrentProvider())

	SetProvider(instance)

	actual := getProvider()

	assert.ToBeSame(t, instance, actual)

	// Cached
	assert.ToBeSame(t, instance, getCurrentProvider())
}

func Test_getProvider_cachedDoesNotWaitForRegistryUpdates(t *testing.T) {
	defer resetGlobal()
	instance := newMockProvider("instance")
	SetProvider(instance)

	done := make(chan Provider, 1)
	knownProvidersMutex.Lock()
	go func() {
		done <- getProvider()
	}()

	var actual Provider
	select {
	case actual = <-done:
	case <-time.After(time.Second):
		knownProvidersMutex.Unlock()
		t.Fatal("cached provider lookup waited for the registry lock")
	}
	knownProvidersMutex.Unlock()

	assert.ToBeSame(t, instance, actual)
}

func Test_SetProvider_waitsForRegistryUpdates(t *testing.T) {
	defer resetGlobal()
	instance := newMockProvider("instance")
	started := make(chan struct{})
	done := make(chan Provider, 1)

	knownProvidersMutex.Lock()
	go func() {
		close(started)
		done <- SetProvider(instance)
	}()
	<-started

	select {
	case <-done:
		knownProvidersMutex.Unlock()
		t.Fatal("provider update completed while the registry lock was held")
	case <-time.After(10 * time.Millisecond):
	}
	knownProvidersMutex.Unlock()

	var previous Provider
	select {
	case previous = <-done:
	case <-time.After(time.Second):
		t.Fatal("provider update did not complete after releasing the registry lock")
	}
	assert.ToBeNil(t, previous)
	assert.ToBeSame(t, instance, getProvider())
}

func Test_getProvider_failsIfMoreThenOneAreRegistered(t *testing.T) {
	defer resetGlobal()
	RegisterProvider(newMockProvider("instance1"))
	RegisterProvider(newMockProvider("instance2"))

	// Nothing cached
	assert.ToBeNil(t, getCurrentProvider())

	assert.Execution(t, func() {
		_ = getProvider()
	}).WillPanicWith("^There are more than provider registered; Got:.*")
}

func Test_RegisterProvider_registersOne(t *testing.T) {
	defer resetGlobal()
	instanceA := newMockProvider("instanceA")
	instanceB := newMockProvider("instanceB")

	assert.ToBeEqual(t, map[string]Provider{}, knownProviders)
	assert.ToBeNil(t, getCurrentProvider())

	actualA := RegisterProvider(instanceA)
	assert.ToBeSame(t, instanceA, actualA)
	assert.ToBeEqual(t, map[string]Provider{"instanceA": instanceA}, knownProviders)
	assert.ToBeNil(t, getCurrentProvider())

	actualB := RegisterProvider(instanceB)
	assert.ToBeSame(t, instanceB, actualB)
	assert.ToBeEqual(t, map[string]Provider{"instanceA": instanceA, "instanceB": instanceB}, knownProviders)
	assert.ToBeNil(t, getCurrentProvider())
}

func Test_RegisterProvider_failsOnNil(t *testing.T) {
	defer resetGlobal()
	assert.Execution(t, func() {
		RegisterProvider(nil)
	}).WillPanicWith("^Provided Provider is nil$")
}

func Test_RegisterProvider_failsOnDoubleRegistration(t *testing.T) {
	defer resetGlobal()
	instance1 := newMockProvider("instance")
	instance2 := newMockProvider("instance")

	assert.ToBeEqual(t, map[string]Provider{}, knownProviders)
	assert.ToBeNil(t, getCurrentProvider())

	actual := RegisterProvider(instance1)
	assert.ToBeSame(t, instance1, actual)
	assert.ToBeEqual(t, map[string]Provider{"instance": instance1}, knownProviders)
	assert.ToBeNil(t, getCurrentProvider())

	assert.Execution(t, func() {
		RegisterProvider(instance2)
	}).WillPanicWith("^Multiple try of registering of Provider with the same name: instance\n\tAlready existing type: ")
}

func Test_RegisterProvider_ignoresDoubleRegistrationOfSameInstance(t *testing.T) {
	defer resetGlobal()
	instance := newMockProvider("instance")

	assert.ToBeEqual(t, map[string]Provider{}, knownProviders)
	assert.ToBeNil(t, getCurrentProvider())

	actual1 := RegisterProvider(instance)
	assert.ToBeSame(t, instance, actual1)
	assert.ToBeEqual(t, map[string]Provider{"instance": instance}, knownProviders)
	assert.ToBeNil(t, getCurrentProvider())

	actual2 := RegisterProvider(instance)
	assert.ToBeSame(t, instance, actual2)
	assert.ToBeEqual(t, map[string]Provider{"instance": instance}, knownProviders)
	assert.ToBeNil(t, getCurrentProvider())
}

func Test_UnregisterProvider(t *testing.T) {
	defer resetGlobal()
	instanceA := newMockProvider("instanceA")
	instanceB := newOtherMockProvider("instanceB")
	RegisterProvider(instanceA)
	RegisterProvider(instanceB)
	SetProvider(instanceA)

	actualA := UnregisterProvider("instanceA")
	assert.ToBeSame(t, instanceA, actualA)
	assert.ToBeEqual(t, map[string]Provider{"instanceB": instanceB}, knownProviders)
	assert.ToBeNil(t, getCurrentProvider())

	actualB := UnregisterProvider("instanceB")
	assert.ToBeSame(t, instanceB, actualB)
	assert.ToBeEqual(t, map[string]Provider{}, knownProviders)
	assert.ToBeNil(t, getCurrentProvider())
}

func Test_GetAllProviders_nonRegistered(t *testing.T) {
	defer resetGlobal()

	assert.ToBeEqual(t, []Provider{}, GetAllProviders())
}

func Test_GetAllProviders_afterOneRegistered(t *testing.T) {
	defer resetGlobal()
	instance1 := newMockProvider("instance")
	RegisterProvider(instance1)

	assert.ToBeEqual(t, []Provider{instance1}, GetAllProviders())
}

func resetGlobal() {
	knownProvidersMutex.Lock()
	defer knownProvidersMutex.Unlock()
	knownProviders = map[string]Provider{}
	storeProvider(nil)
}

func getCurrentProvider() Provider {
	if v := (*Provider)(atomic.LoadPointer(&providerPointer)); v != nil && *v != nil {
		return *v
	}
	return nil
}

func newOtherMockProvider(name string) Provider {
	return &otherMockProvider{newMockProvider(name)}
}

type otherMockProvider struct {
	*mockProvider
}
