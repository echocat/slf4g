package log

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var (
	globalProviderFacade = NewProviderFacade(getProvider)
	globalProvider       Provider
	allProviders         = map[string]Provider{}
	providerMutex        sync.RWMutex
)

func GetProvider() Provider {
	// We're using this facade to deal with concurrency issues where someone already
	// addresses the reference to the global available provider but afterwards
	// the real provider is initiated.
	return globalProviderFacade
}

func getProvider() Provider {
	providerMutex.RLock()
	rLocked := true
	defer func() {
		if rLocked {
			providerMutex.RUnlock()
		}
	}()

	if p := globalProvider; p != nil {
		return p
	}

	providerMutex.RUnlock()
	rLocked = false
	providerMutex.Lock()
	defer providerMutex.Unlock()

	if p := globalProvider; p != nil {
		return p
	}

	if len(allProviders) > 1 {
		asStrings := make([]string, len(allProviders))
		var i int
		for n, p := range allProviders {
			asStrings[i] = fmt.Sprintf("%d# %v(%s)", i+1, reflect.TypeOf(p), n)
			i++
		}
		panic("There are more than provider registered; Got:\n\t" + strings.Join(asStrings, "\n\t"))
	}

	for _, p := range allProviders {
		globalProvider = p
		return p
	}

	// Everything failed, using a simple provider now...
	globalProvider = simpleProviderV
	return simpleProviderV
}

func SetProvider(p Provider) {
	providerMutex.Lock()
	defer providerMutex.Unlock()

	globalProvider = p
}

func RegisterProvider(p Provider) Provider {
	if p == nil {
		panic("Provided provider is nil")
	}

	providerMutex.Lock()
	defer providerMutex.Unlock()

	if existing := allProviders[p.GetName()]; existing == p {
		return p
	} else if existing != nil {
		panic(fmt.Sprintf("Multiple try of registring provider of provider with the same name: %s\n"+
			"\tAlready existing type: %v; new type: %v", p.GetName(), reflect.TypeOf(existing), reflect.TypeOf(p)))
	}

	allProviders[p.GetName()] = p

	return p
}

func UnregisterProvider(name string) Provider {
	providerMutex.Lock()
	defer providerMutex.Unlock()

	existing := allProviders[name]

	delete(allProviders, name)

	return existing
}

func AllProviders() []Provider {
	providerMutex.RLock()
	defer providerMutex.RUnlock()

	result := make([]Provider, len(allProviders))

	var i int
	for _, p := range allProviders {
		result[i] = p
		i++
	}

	return result
}
