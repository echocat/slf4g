package level

// NewProviderFacade creates a new facade of Provider with the given
// function that provides the actual Provider to use.
func NewProviderFacade(provider func() Provider) Provider {
	return providerFacade(provider)
}

type providerFacade func() Provider

func (instance providerFacade) GetName() string {
	return instance.Unwrap().GetName()
}

func (instance providerFacade) GetLevels() Levels {
	return instance.Unwrap().GetLevels()
}

func (instance providerFacade) Unwrap() Provider {
	return instance()
}
