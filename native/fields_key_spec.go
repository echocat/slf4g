package native

import (
	"github.com/echocat/slf4g/fields"
)

var (
	DefaultFieldsKeySpec       = NewFieldsKeySpec()
	DefaultFieldsKeySpecFacade = NewFieldsKeySpecFacade(func() FieldsKeysSpec {
		return DefaultFieldsKeySpec
	})
)

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

func NewFieldsKeySpecFacade(provider func() FieldsKeysSpec) FieldsKeysSpec {
	return fieldsKeySpecFacade(provider)
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

type fieldsKeySpecFacade func() FieldsKeysSpec

func (instance fieldsKeySpecFacade) GetTimestamp() string {
	return instance().GetTimestamp()
}

func (instance fieldsKeySpecFacade) GetMessage() string {
	return instance().GetMessage()
}

func (instance fieldsKeySpecFacade) GetError() string {
	return instance().GetError()
}

func (instance fieldsKeySpecFacade) GetLogger() string {
	return instance().GetLogger()
}

func (instance fieldsKeySpecFacade) GetLocation() string {
	return instance().GetLocation()
}
