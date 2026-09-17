package log

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/echocat/slf4g/internal/globalprovider"
)

var (
	providerV      = NewProviderFacade(getProvider)
	knownProviders = map[string]Provider{}
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
	return globalprovider.Resolve(func() interface{} {
		return exactOneProvider()
	}).(Provider)
}

func exactOneProvider() Provider {
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
	previous := globalprovider.Set(p)
	if previous != nil {
		return previous.(Provider)
	}
	return nil
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
	name := p.GetName()

	globalprovider.Invalidate(func() {
		if existing := knownProviders[name]; existing == p {
			return
		} else if existing != nil {
			panic(fmt.Sprintf("Multiple try of registering of Provider with the same name: %s\n"+
				"\tAlready existing type: %v; new type: %v", name, reflect.TypeOf(existing), reflect.TypeOf(p)))
		}

		knownProviders[name] = p
	})
	return p
}

// UnregisterProvider is doing the exact opposite of RegisterProvider().
func UnregisterProvider(name string) Provider {
	var existing Provider
	globalprovider.Invalidate(func() {
		existing = knownProviders[name]
		delete(knownProviders, name)
	})

	return existing
}

// GetAllProviders returns all knows providers registered with
// RegisterProvider().
func GetAllProviders() []Provider {
	var result []Provider
	globalprovider.Read(func() {
		result = make([]Provider, len(knownProviders))

		var i int
		for _, p := range knownProviders {
			result[i] = p
			i++
		}
	})

	return result
}

func loadProvider() Provider {
	if value := globalprovider.Get(); value != nil {
		return value.(Provider)
	}
	return nil
}
