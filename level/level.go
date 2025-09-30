// Package level provides Level which identifies the severity of an event to
// be logged. As higher as more important is the event. Trace is the less
// severe and Fatal the most severe.
//
// # Customization
//
// Different implementations of Provider could introduce more Levels. All
// ordinals are unique over all instances of Level. Info = 3000 will be
// always Info = 3000. Another Level which uses the ordinal 3000 can be
// just assumed as an alias to Info.  Customization only means added new
// instances of Level. Standard levels always remains available.
package level

import "errors"

const (
	// Trace defines the lowest possible level. This is usually only used
	// in cases where really detailed logs are required. For example to document
	// the whole communication with a server, including each request with its
	// headers and so on.
	Trace Level = 1000

	// Debug is used to document information which is used to debug
	// problems. This information are in regular operation mode are not
	// required; but could help once enabled to track down common issues.
	Debug Level = 2000

	// Info is the regular level where everything is logged which is of
	// interest for everybody and should be always visible and imply regular
	// operations of the system. Usually this shows that one operation succeeded
	// successfully; like a user was created.
	Info Level = 3000

	// Warn implies that something happened which failed but could be
	// recovered gracefully. In best case the user does not notice anything.
	// But personal should investigate as soon as possible to prevent something
	// like this happening again.
	Warn Level = 4000

	// Error implies that something happened which failed and cannot be
	// recovered gracefully, but it only affects one or a small amount of users
	// and the rest of the system can continue to work. Personal should
	// investigate right now to prevent such stuff happening again and recover
	// broken users.
	Error Level = 5000

	// Fatal implies that something happened which failed and cannot be
	// recovered gracefully, which might affect every user of the system. This
	// implies that the whole system is no longer operable and should/will be
	// shutdown (if possible gracefully) right now. Personal is required to
	// investigate right now to prevent such stuff happening again, bring the
	// the system back to operations and recover broken users.
	Fatal Level = 6000
)

type Level uint16

// CompareTo helps to compare the severity of two levels. It returns how much
// severe the provided Level is compared to the actual one. Bigger means more
// severe.
func (instance Level) CompareTo(o Level) int {
	return int(instance) - int(o)
}

// Levels represents a slice of 0...n Level.
//
// Additionally, it implements the sort.Interface to enable to call sort.Sort
// to order the contents of this slice by its severity.
type Levels []Level

// Len implements sort.Interface
func (instance Levels) Len() int {
	return len(instance)
}

// Swap implements sort.Interface
func (instance Levels) Swap(i, j int) {
	instance[i], instance[j] = instance[j], instance[i]
}

// Less implements sort.Interface
func (instance Levels) Less(i, j int) bool {
	return instance[i].CompareTo(instance[j]) < 0
}

// ToProvider transforms the current Levels into a Provider instance with the
// given name.
func (instance Levels) ToProvider(name string) Provider {
	return levelsAsProvider{instance, name}
}

type levelsAsProvider struct {
	values Levels
	name   string
}

func (instance levelsAsProvider) GetName() string {
	return instance.name
}

func (instance levelsAsProvider) GetLevels() Levels {
	return instance.values
}

// ErrIllegalLevel represents that an illegal level.Level value/name was
// provided.
var ErrIllegalLevel = errors.New("illegal level")
