package fields

var DefaultKeysSpec KeysSpec = &defaultFieldKeysSpec{}

type KeysSpec interface {
	GetTimestamp() string
	GetMessage() string
	GetError() string
	GetLogger() string
}

type defaultFieldKeysSpec struct{}

func (instance *defaultFieldKeysSpec) GetTimestamp() string {
	return "timestamp"
}

func (instance *defaultFieldKeysSpec) GetMessage() string {
	return "message"
}

func (instance *defaultFieldKeysSpec) GetError() string {
	return "error"
}

func (instance *defaultFieldKeysSpec) GetLogger() string {
	return "logger"
}
