package fields

type KeysSpec interface {
	GetTimestamp() string
	GetMessage() string
	GetError() string
	GetLogger() string
}

type KeysSpecProvider func() KeysSpec
