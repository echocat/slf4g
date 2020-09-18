package native

import (
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/native/location"
)

var DefaultFieldsKeySpec = NewFieldsKeySpec()

type FieldsKeysSpec interface {
	fields.KeysSpec

	GetLocation() string
}

func NewFieldsKeySpec() FieldsKeysSpec {
	def := fields.DefaultKeysSpec
	return &FieldsKeySpecImpl{
		KeysSpec: def,
		Location: location.Field,
		Logger:   def.GetLogger(),
		Error:    def.GetError(),
	}
}

type FieldsKeySpecImpl struct {
	fields.KeysSpec

	Location string
	Logger   string
	Error    string
}

func (instance *FieldsKeySpecImpl) GetLocation() string {
	return instance.Location
}

func (instance *FieldsKeySpecImpl) GetLogger() string {
	return instance.Location
}

func (instance *FieldsKeySpecImpl) GetError() string {
	return instance.Error
}
