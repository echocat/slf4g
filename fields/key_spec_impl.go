package fields

// KeysSpecImpl is a default implementation of KeysSpec.
type KeysSpecImpl struct {
	// Timestamp defines the used key of a timestamp.
	// If empty "timestamp" will be used instead.
	Timestamp string

	// Message defines the used key of a message.
	// If empty "message" will be used instead.
	Message string

	// Logger defines the used key of a logger.
	// If empty "logger" will be used instead.
	Logger string

	// Error defines the used key of an error.
	// If empty "error" will be used instead.
	Error string
}

// GetTimestamp implements KeysSpec#GetTimestamp()
func (instance KeysSpecImpl) GetTimestamp() string {
	if v := instance.Timestamp; v != "" {
		return v
	}
	return "timestamp"
}

// GetMessage implements KeysSpec#GetMessage()
func (instance KeysSpecImpl) GetMessage() string {
	if v := instance.Message; v != "" {
		return v
	}
	return "message"
}

// GetError implements KeysSpec#GetError()
func (instance KeysSpecImpl) GetError() string {
	if v := instance.Error; v != "" {
		return v
	}
	return "error"
}

// GetLogger implements KeysSpec#GetLogger()
func (instance KeysSpecImpl) GetLogger() string {
	if v := instance.Logger; v != "" {
		return v
	}
	return "logger"
}
