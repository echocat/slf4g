package level

// Names is used to make readable names out of level.Level or the other way
// around.
type Names interface {
	// ToName converts a given level.Level to a human-readable name. If this
	// level is unknown by this instance an error is returned. Most likely
	// ErrIllegalLevel.
	ToName(Level) (string, error)

	// ToLevel converts a given human-readable name to a level.Level. If this
	// name is unknown by this instance an error is returned. Most likely
	// ErrIllegalLevel.
	ToLevel(string) (Level, error)
}

// NamesAware represents an object that is aware of Names.
type NamesAware interface {
	// GetLevelNames returns an instance of level.Names that support by
	// formatting levels in a human-readable format.
	GetLevelNames() Names
}

// NewNamesFacade creates a facade of Names using the given provider.
func NewNamesFacade(provider func() Names) Names {
	return namesFacade(provider)
}

type namesFacade func() Names

func (instance namesFacade) ToName(lvl Level) (string, error) {
	return instance.Unwrap().ToName(lvl)
}

func (instance namesFacade) ToLevel(name string) (Level, error) {
	return instance.Unwrap().ToLevel(name)
}

func (instance namesFacade) Unwrap() Names {
	return instance()
}
