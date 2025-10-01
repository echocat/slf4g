package eventlog

import (
	"github.com/echocat/slf4g/fields"
)

// DefaultFieldKeysSpec is the default instance of FieldKeysSpec which should
// cover the majority of cases.
var DefaultFieldKeysSpec = &FieldKeysSpecImpl{}

// FieldKeysSpec defines the field keys supported this implementation of slf4g.
//
// It is an extension of the default fields.KeysSpec.
type FieldKeysSpec interface {
	fields.KeysSpec

	// GetLocation defines the key location information of logged event are
	// stored inside. Such as the calling method, ...
	GetLocation() string
}

// FieldKeysSpecImpl is a default implementation of FieldKeysSpec.
type FieldKeysSpecImpl struct {
	fields.KeysSpecImpl

	// Location defines the used key of an location.
	// If empty "location" will be used instead.
	Location string
}

// GetLocation implements FieldKeysSpec#GetLocation()
func (instance *FieldKeysSpecImpl) GetLocation() string {
	if v := instance.Location; v != "" {
		return v
	}
	return "location"
}

// NewFieldKeysSpecFacade creates a facade of FieldKeysSpec using the given
// provider.
func NewFieldKeysSpecFacade(provider func() FieldKeysSpec) FieldKeysSpec {
	return &fieldKeysSpecFacade{
		KeysSpec: fields.NewKeysSpecFacade(func() fields.KeysSpec {
			return provider()
		}),
		provider: provider,
	}
}

type fieldKeysSpecFacade struct {
	fields.KeysSpec
	provider func() FieldKeysSpec
}

func (instance *fieldKeysSpecFacade) GetLocation() string {
	return instance.Unwrap().GetLocation()
}

func (instance *fieldKeysSpecFacade) Unwrap() FieldKeysSpec {
	return instance.provider()
}
