package native

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/echocat/slf4g/level"
)

type levelState struct {
	mutexPointer unsafe.Pointer
}

func (instance *levelState) load(configured *level.Level) level.Level {
	mutex := instance.mutex()
	mutex.RLock()
	defer mutex.RUnlock()
	return *configured
}

func (instance *levelState) store(configured *level.Level, value level.Level) {
	mutex := instance.mutex()
	mutex.Lock()
	defer mutex.Unlock()
	*configured = value
}

func (instance *levelState) mutex() *sync.RWMutex {
	for {
		if current := (*sync.RWMutex)(atomic.LoadPointer(&instance.mutexPointer)); current != nil {
			return current
		}
		created := &sync.RWMutex{}
		if atomic.CompareAndSwapPointer(&instance.mutexPointer, nil, unsafe.Pointer(created)) {
			return created
		}
	}
}
