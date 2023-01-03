package fields

// KeysSpec defines the keys for common usages inside a Fields instance.
type KeysSpec interface {
	// GetTimestamp returns the key where the timestamp is stored inside, if
	// available.
	GetTimestamp() string

	// GetMessage returns the key where the message is stored inside, if
	// available.
	GetMessage() string

	// GetError returns the key where an error is stored inside, if
	// available.
	GetError() string

	// GetLogger returns the key where the Logger is stored inside, which is
	// managed a Fields instance. (if available)
	GetLogger() string
}

// NewKeysSpecFacade creates a facade of KeysSpec using the given provider.
func NewKeysSpecFacade(provider func() KeysSpec) KeysSpec {
	return keysSpecFacade(provider)
}

type keysSpecFacade func() KeysSpec

func (instance keysSpecFacade) GetTimestamp() string {
	return instance.Unwrap().GetTimestamp()
}

func (instance keysSpecFacade) GetMessage() string {
	return instance.Unwrap().GetMessage()
}

func (instance keysSpecFacade) GetError() string {
	return instance.Unwrap().GetError()
}

func (instance keysSpecFacade) GetLogger() string {
	return instance.Unwrap().GetLogger()
}

func (instance keysSpecFacade) Unwrap() KeysSpec {
	return instance()
}
