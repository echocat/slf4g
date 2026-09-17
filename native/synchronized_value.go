package native

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type synchronizedValue[T any] struct {
	mutexPointer unsafe.Pointer
}

func (instance *synchronizedValue[T]) load(configured *T) T {
	mutex := instance.mutex()
	mutex.RLock()
	defer mutex.RUnlock()
	return *configured
}

func (instance *synchronizedValue[T]) store(configured *T, value T) {
	mutex := instance.mutex()
	mutex.Lock()
	defer mutex.Unlock()
	*configured = value
}

func (instance *synchronizedValue[T]) mutex() *sync.RWMutex {
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
