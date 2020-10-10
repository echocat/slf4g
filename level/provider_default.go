package level

var defaultProviderV = &defaultProvider{"default"}

type defaultProvider struct {
	name string
}

func (instance *defaultProvider) GetName() string {
	return instance.name
}

func (instance *defaultProvider) GetLevels() Levels {
	return Levels{Trace, Debug, Info, Warn, Error, Fatal}
}
