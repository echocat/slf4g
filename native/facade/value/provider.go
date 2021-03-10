package value

// ProviderTarget defines an object that receives the Level and Formatter managed
// by the Provider value facade.
type ProviderTarget interface {
	LevelTarget
	ConsumerTarget
}

// Provider is a value facade for transparent setting of native.Provider
// for the slf4g/native implementation. This is quite handy for usage
// with flags package of the SDK or similar flag libraries. This might
// be usable, too in contexts where serialization might be required.
type Provider struct {
	// Level is the corresponding level.Level facade.
	Level Level

	// Consumer is the corresponding consumer.Consumer facade.
	Consumer Consumer
}

// NewProvider create a new instance of Provider with the given target ProviderTarget instance.
func NewProvider(target ProviderTarget, customizer ...func(*Provider)) Provider {
	result := Provider{
		Level:    NewLevel(target),
		Consumer: NewConsumer(target),
	}

	for _, c := range customizer {
		c(&result)
	}

	return result
}
