package fields

import "github.com/echocat/slf4g/level"

// Filtered is a value which will be executed on usage to retrieve the actual
// value or exclude it.
//
// This is useful in context where fields should be only respected based on a
// specific log level, another field has a specific value, ...
type Filtered interface {
	// Filter is the method which will be called at the moment where the value
	// should be consumed.
	//
	// Only if shouldBeRespected is true it will be respected by the consumers.
	Filter(FilterContext) (value interface{}, shouldBeRespected bool)

	// Get will return the original value (unfiltered).
	Get() interface{}
}

// FilterContext provides information about the context where a field exists
// within.
type FilterContext interface {
	// GetLevel provides the current level.Level of this context.
	GetLevel() level.Level

	// Get provides access to other fields within this context.
	Get(key string) (value interface{}, exists bool)
}

// RequireMaximalLevel represents a filtered value which will only be consumed if the
// level.Level of the current context (for example logging events) is not bigger than
// the requested maximalLevel.
func RequireMaximalLevel(maximalLevel level.Level, value interface{}) Filtered {
	return RequireMaximalLevelLazy(maximalLevel, LazyFunc(func() interface{} {
		return value
	}))
}

// RequireMaximalLevelLazy represents a filtered Lazy value which will only be consumed
// if the level.Level of the current context (for example logging events) is not bigger
// than requested maximalLevel.
func RequireMaximalLevelLazy(minimalLevel level.Level, value Lazy) Filtered {
	return requireMaximalLevel{value, minimalLevel}
}

type requireMaximalLevel struct {
	Lazy
	level level.Level
}

func (instance requireMaximalLevel) Filter(ctx FilterContext) (value interface{}, shouldBeRespected bool) {
	if ctx.GetLevel() > instance.level {
		return nil, false
	}

	return instance.Get(), true
}

// IgnoreLevels represents a filtered value which will only be consumed if the
// level.Level of the current context (for example logging events) is smaller than
// fromLevel or equal/bigger than toLevel (fromLevel:inclusive, toLevel:exclusive).
func IgnoreLevels(fromLevel, toLevel level.Level, value interface{}) Filtered {
	return IgnoreLevelsLazy(fromLevel, toLevel, LazyFunc(func() interface{} {
		return value
	}))
}

// IgnoreLevelsLazy represents a filtered Lazy value which will only be consumed
// if the level.Level of the current context (for example logging events) is smaller
// than fromLevel or equal/bigger than toLevel (fromLevel:inclusive, toLevel:exclusive).
func IgnoreLevelsLazy(fromLevel, toLevel level.Level, value Lazy) Filtered {
	return ignoreLevels{value, fromLevel, toLevel}
}

type ignoreLevels struct {
	Lazy
	fromLevel level.Level
	toLevel   level.Level
}

func (instance ignoreLevels) Filter(ctx FilterContext) (value interface{}, shouldBeRespected bool) {
	lvl := ctx.GetLevel()
	if lvl >= instance.fromLevel && lvl < instance.toLevel {
		return nil, false
	}

	return instance.Get(), true
}
