package log

import (
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
)

type Provider interface {
	// GetRootLogger returns the root Logger. The Provider guarantees that there
	// is always the same Logger.
	GetRootLogger() Logger

	// GetLogger returns a Logger for the given name. The Provider guarantees
	// that there is always the same Logger for the same name returned.
	GetLogger(name string) Logger

	// GetName returns the name of this Provider instance.
	GetName() string

	// GetAllLevels returns all available level.Levels which are supported by
	// this Provider.
	GetAllLevels() level.Levels

	// GetFieldKeysSpec returns the fields.KeysSpec which describes the keys of
	// fields.Fields which are supported by this Provider.
	GetFieldKeysSpec() fields.KeysSpec
}

// NewProviderFacade creates a new facade of Provider with the given
// function that provides the actual Provider to use.
func NewProviderFacade(provider func() Provider) Provider {
	return providerFacade(provider)
}
