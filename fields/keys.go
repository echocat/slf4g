package fields

var DefaultKeys Keys = &defaultFieldKeys{}

type Keys interface {
	GetTimestamp() string
	GetMessage() string
	GetError() string
	GetLogger() string
}

type defaultFieldKeys struct{}

func (instance *defaultFieldKeys) GetTimestamp() string {
	return "timestamp"
}

func (instance *defaultFieldKeys) GetMessage() string {
	return "message"
}

func (instance *defaultFieldKeys) GetError() string {
	return "error"
}

func (instance *defaultFieldKeys) GetLogger() string {
	return "logger"
}
