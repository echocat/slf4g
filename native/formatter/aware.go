package formatter

// Aware describes an object that is aware of a Formatter and exports its current
// state.
type Aware interface {
	// GetFormatter returns the current formatter.
	GetFormatter() Formatter
}

// MutableAware is similar to Aware but additionally is able to modify the
// Formatter by calling SetFormatter(Formatter).
type MutableAware interface {
	Aware

	// SetFormatter modifies the current formatter to the given one.
	SetFormatter(Formatter)
}
