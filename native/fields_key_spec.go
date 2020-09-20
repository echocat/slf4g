package native

import (
	"github.com/echocat/slf4g/fields"
)

var DefaultFieldsKeySpec = NewFieldsKeySpec()

type FieldsKeysSpec interface {
	fields.KeysSpec
	GetLocation() string
}

func NewFieldsKeySpec() FieldsKeysSpec {
	return &FieldsKeySpecImpl{
		Timestamp: "timestamp",
		Message:   "message",
		Logger:    "logger",
		Error:     "error",
		Location:  "location",
	}
}

type FieldsKeySpecImpl struct {
	Timestamp string
	Message   string
	Logger    string
	Error     string
	Location  string
}

func (instance *FieldsKeySpecImpl) GetTimestamp() string {
	return instance.Timestamp
}

func (instance *FieldsKeySpecImpl) GetMessage() string {
	return instance.Message
}

func (instance *FieldsKeySpecImpl) GetLogger() string {
	return instance.Location
}

func (instance *FieldsKeySpecImpl) GetError() string {
	return instance.Error
}

func (instance *FieldsKeySpecImpl) GetLocation() string {
	return instance.Location
}
