package fields

import "fmt"

// Field represents a single field which are usually contained in Fields.
type Field interface {
	// Key returns the key of the field.
	Key() string
	// Value returns the actual value of the field.
	Value() interface{}

	// String returns a simple representation of this field in format <key>=<value>
	String() string
}

// NewField creates a new field from the given key and value.
func NewField(key string, value interface{}) Field {
	return field{key: key, value: value}
}

type field struct {
	key   string
	value interface{}
}

func (instance field) Key() string {
	return instance.key
}

func (instance field) Value() interface{} {
	return instance.value
}

func (instance field) String() string {
	return fmt.Sprintf("%s=%v", instance.key, instance.value)
}
