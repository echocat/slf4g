package consumer

// Aware describes an object that is aware of a Consumer and exports its current
// state.
type Aware interface {
	// GetConsumer returns the current consumer.
	GetConsumer() Consumer
}

// MutableAware is similar to Aware but additionally is able to modify the
// Consumer by calling SetConsumer(Consumer).
type MutableAware interface {
	Aware

	// SetConsumer modifies the current consumer to the given one.
	SetConsumer(Consumer)
}
