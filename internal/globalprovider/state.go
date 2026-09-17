// Package globalprovider coordinates global provider installations and scoped
// overrides without exposing additional public API.
package globalprovider

import (
	"sync"
	"sync/atomic"
)

type entry struct {
	value    interface{}
	previous *entry
	scoped   bool
	active   bool
}

var (
	mutex   sync.Mutex
	current atomic.Value
)

func init() {
	current.Store(&entry{})
}

func Get() interface{} {
	return load().value
}

func Set(value interface{}) interface{} {
	mutex.Lock()
	defer mutex.Unlock()

	previous := load().value
	store(&entry{value: value})
	return previous
}

func Resolve(resolver func() interface{}) interface{} {
	if value := Get(); value != nil {
		return value
	}

	mutex.Lock()
	defer mutex.Unlock()

	if value := load().value; value != nil {
		return value
	}

	value := resolver()
	store(&entry{value: value})
	return value
}

func Invalidate(change func()) {
	mutex.Lock()
	defer mutex.Unlock()
	defer store(&entry{})

	change()
}

func Read(reader func()) {
	mutex.Lock()
	defer mutex.Unlock()

	reader()
}

func Push(value interface{}) func() {
	mutex.Lock()
	installed := &entry{
		value:    value,
		previous: load(),
		scoped:   true,
		active:   true,
	}
	store(installed)
	mutex.Unlock()

	return func() {
		mutex.Lock()
		defer mutex.Unlock()

		if installed == nil {
			return
		}
		target := installed
		installed = nil
		currentEntry := load()
		target.active = false
		if currentEntry != target {
			if !contains(currentEntry, target) {
				detach(target)
			}
			return
		}

		previous := target.previous
		for previous != nil && previous.scoped && !previous.active {
			skipped := previous
			previous = previous.previous
			skipped.previous = nil
		}
		target.previous = nil
		if previous == nil {
			previous = &entry{}
		}
		store(previous)
	}
}

func load() *entry {
	return current.Load().(*entry)
}

func store(value *entry) {
	current.Store(value)
}

func contains(current, expected *entry) bool {
	for current != nil {
		if current == expected {
			return true
		}
		current = current.previous
	}
	return false
}

func detach(current *entry) {
	for current != nil && current.scoped {
		previous := current.previous
		current.previous = nil
		current = previous
	}
}
