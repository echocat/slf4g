// Fields represents a collection of key value pairs.
//
// Key and Values
//
// The keys are always of type string and should be only printable characters
// which can be printed in any context. Recommended are everything that matches:
//   ^[a-zA-Z0-9._-]+$
//
// The values could be everything including nils.
//
// Immutability
//
// Fields are defined as immutable. Calling the methods With, Withf, Without and
// WithAll always results in a new instance of Fields that could be either
// brand new, a copy of the source one or do inherit some stuff of the original
// called one; but it never modifies the called instance.
package fields

type Fields interface {
	// ForEach will call the provided Consumer for each field which is provided
	// by this Fields instance.
	ForEach(consumer Consumer) error

	// Get will return for the given key the corresponding value if exists.
	// Otherwise it will return nil.
	Get(key string) interface{}

	// With returns an variant of this Fields with the given key
	// value pair contained inside. If the given key already exists in the
	// current instance this means it will be overwritten.
	With(key string, value interface{}) Fields

	// Withf is similar to With but it adds classic fmt.Printf functions to it.
	// It is defined that the format itself will not be executed before the
	// consumption of the value. (See ForEach() and Get())
	Withf(key string, format string, args ...interface{}) Fields

	// WithAll is similar to With but it can consume more than one field at
	// once. Be aware: There is neither a guarantee that this instance will be
	// copied or not.
	WithAll(map[string]interface{}) Fields

	// Without returns a variant of this Fields without the given
	// key contained inside. In other words: If someone afterwards tries to
	// call either ForEach() or Get() nothing with this key(s) will be returned.
	Without(keys ...string) Fields
}

// Consumer will be called on each field that needs to be consumed.
//
// See Fields.ForEach() for more details.
type Consumer func(key string, value interface{}) error
