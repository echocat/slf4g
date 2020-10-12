package level

// Aware describes an object that is aware of a Level and exports its current
// state.
type Aware interface {
	// GetLevel returns the current level.
	GetLevel() Level
}

// MutableAware is similar to Aware but additionally is able to modify the
// Level by calling SetLevel(Level).
type MutableAware interface {
	Aware

	// SetLevel modifies the current level to the given one.
	SetLevel(Level)
}
