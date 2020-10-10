package level

import (
	"sync/atomic"
	"testing"

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

	assert.ToBeSame(t, defaultProviderV, actual)
	// Also now cached
	assert.ToBeSame(t, defaultProviderV, getCurrentProvider())
}

func Test_getProvider_returnsRegistered(t *testing.T) {
	defer resetGlobal()
	instance := &defaultProvider{"instance"}
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
	instance := &defaultProvider{"instance"}

	// Nothing cached
	assert.ToBeNil(t, getCurrentProvider())

	SetProvider(instance)

	actual := getProvider()

	assert.ToBeSame(t, instance, actual)

	// Cached
	assert.ToBeSame(t, instance, getCurrentProvider())
}

func Test_getProvider_failsIfMoreThenOneAreRegistered(t *testing.T) {
	defer resetGlobal()
	RegisterProvider(&defaultProvider{"instance1"})
	RegisterProvider(&defaultProvider{"instance2"})

	// Nothing cached
	assert.ToBeNil(t, getCurrentProvider())

	assert.Execution(t, func() {
		_ = getProvider()
	}).WillPanicWith("^There are more than provider registered; Got:.*")
}

func Test_RegisterProvider_registersOne(t *testing.T) {
	defer resetGlobal()
	instanceA := &defaultProvider{"instanceA"}
	instanceB := &defaultProvider{"instanceB"}
	SetProvider(defaultProviderV)

	assert.ToBeEqual(t, map[string]Provider{}, knownProviders)
	assert.ToBeNotNil(t, getCurrentProvider())

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
	instance1 := &defaultProvider{"instance"}
	instance2 := &defaultProvider{"instance"}

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
	instance := &defaultProvider{"instance"}

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
	instanceA := &defaultProvider{"instanceA"}
	instanceB := &otherProvider{defaultProvider{"instanceB"}}
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

	assert.ToBeEqual(t, []Provider{defaultProviderV}, GetAllProviders())
}

func Test_GetAllProviders_afterOneRegistered(t *testing.T) {
	defer resetGlobal()
	instance1 := &defaultProvider{"instance"}
	RegisterProvider(instance1)

	assert.ToBeEqual(t, []Provider{instance1, defaultProviderV}, GetAllProviders())
}

func resetGlobal() {
	SetProvider(nil)
	knownProvidersMutex.Lock()
	defer knownProvidersMutex.Unlock()
	knownProviders = map[string]Provider{}
}

func getCurrentProvider() Provider {
	if v := (*Provider)(atomic.LoadPointer(&providerPointer)); v != nil && *v != nil {
		return *v
	}
	return nil
}

type otherProvider struct {
	defaultProvider
}
