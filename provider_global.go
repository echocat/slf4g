package log

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

var (
	providerPointer     unsafe.Pointer
	providerV           = NewProviderFacade(getProvider)
	knownProviders      = map[string]Provider{}
	knownProvidersMutex sync.RWMutex
)

// GetProvider returns the actual Provider.
//
// The Provider returned by this method is guarded by a facade
// (see NewProviderFacade()) which ensures that usages of this Provider will
// always call the global configured Provider depending on whether the Provider
// was configured before calling this method or afterwards.
func GetProvider() Provider {
	return providerV
}

func getProvider() Provider {
	for {
		v := (*Provider)(atomic.LoadPointer(&providerPointer))
		if v != nil && *v != nil {
			return *v
		}

		p := exactOneProvider()

		if atomic.CompareAndSwapPointer(&providerPointer, unsafe.Pointer(v), unsafe.Pointer(&p)) {
			return p
		}
	}
}

func exactOneProvider() Provider {
	knownProvidersMutex.RLock()
	defer knownProvidersMutex.RUnlock()

	if len(knownProviders) > 1 {
		asStrings := make([]string, len(knownProviders))
		var i int
		for n, p := range knownProviders {
			asStrings[i] = fmt.Sprintf("%d# %v(%s)", i+1, reflect.TypeOf(p), n)
			i++
		}
		panic("There are more than provider registered; Got:\n\t" + strings.Join(asStrings, "\n\t"))
	}

	for _, p := range knownProviders {
		return p
	}

	// Everything failed, using the fallbackProvider now...
	return fallbackProviderV
}

// SetProvider forces the given Provider as the actual one which will be
// returned when calling GetProvider(). This will lead to that each Provider
// registered with RegisterProvider() will be ignored.
//
// This methods accepts also <nil>. In this case regular discovery mechanism
// will be enabled again when calling GetProvider().
//
// This methods always returns the previous set value (which can be <nil>, too).
func SetProvider(p Provider) Provider {
	for {
		oldP := (*Provider)(atomic.LoadPointer(&providerPointer))
		if atomic.CompareAndSwapPointer(&providerPointer, unsafe.Pointer(oldP), unsafe.Pointer(&p)) {
			if oldP != nil {
				return *oldP
			}
			return nil
		}
	}
}

// RegisterProvider registers the given provider as a usable one. If
// GetProvider() is called each Provider which was registered with this method
// will be then taken into consideration to be returned.
//
// It is not possible to register more than one instance of a Provider with the
// same name.
func RegisterProvider(p Provider) Provider {
	if p == nil {
		panic("Provided Provider is nil")
	}

	knownProvidersMutex.Lock()
	defer SetProvider(nil)
	defer knownProvidersMutex.Unlock()

	if existing := knownProviders[p.GetName()]; existing == p {
		return p
	} else if existing != nil {
		panic(fmt.Sprintf("Multiple try of registering of Provider with the same name: %s\n"+
			"\tAlready existing type: %v; new type: %v", p.GetName(), reflect.TypeOf(existing), reflect.TypeOf(p)))
	}

	knownProviders[p.GetName()] = p

	return p
}

// UnregisterProvider is doing the exact opposite of RegisterProvider().
func UnregisterProvider(name string) Provider {
	knownProvidersMutex.Lock()
	defer SetProvider(nil)
	defer knownProvidersMutex.Unlock()

	existing := knownProviders[name]

	delete(knownProviders, name)

	return existing
}

// GetAllProviders returns all knows providers registered with
// RegisterProvider().
func GetAllProviders() []Provider {
	knownProvidersMutex.RLock()
	defer knownProvidersMutex.RUnlock()

	result := make([]Provider, len(knownProviders))

	var i int
	for _, p := range knownProviders {
		result[i] = p
		i++
	}

	return result
}
