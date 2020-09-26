package fields

// KeysSpec defines the keys for common usages inside of a Fields instance.
type KeysSpec interface {
	GetTimestamp() string
	GetMessage() string
	GetError() string
	GetLogger() string
}

// KeysSpecProvider provides a populated instance of KeysSpec.
type KeysSpecProvider func() KeysSpec
