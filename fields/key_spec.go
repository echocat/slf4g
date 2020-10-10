package fields

// KeysSpec defines the keys for common usages inside of a Fields instance.
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

// KeysSpecProvider provides a populated instance of KeysSpec.
type KeysSpecProvider func() KeysSpec
