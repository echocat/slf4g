package consumer

import (
	"fmt"
	"sync"

	log "github.com/echocat/slf4g"
)

// Recorder is an implementation of a Consumer which only records all logged
// events and makes it able to Get() them afterwards from this Recorder.
type Recorder struct {
	// Synchronized defines if this instance can be used in concurrent
	// environments; which is meaningful in the most context. It might have
	// additional performance costs.
	Synchronized bool

	recorded []log.Event

	mutex sync.RWMutex
}

// NewRecorder creates a new instance of Recorder which can be customized using
// customizer and is ready to use. The created instance is synchronized by
// default (See Recorder.Synchronized).
func NewRecorder(customizer ...func(*Recorder)) *Recorder {
	result := &Recorder{
		Synchronized: true,
	}
	for _, c := range customizer {
		c(result)
	}
	return result
}

// Consume implements Consumer.Consume()
func (instance *Recorder) Consume(event log.Event, _ log.CoreLogger) {
	if instance.Synchronized {
		instance.mutex.Lock()
		defer instance.mutex.Unlock()
	}

	instance.recorded = append(instance.recorded, event)
}

// Len returns the amount of recorded events so far.
func (instance *Recorder) Len() int {
	instance.mutex.RLock()
	defer instance.mutex.RUnlock()

	return len(instance.recorded)
}

// Get return an event at the given index. If this index does not exists this
// method will panic.
func (instance *Recorder) Get(index uint) log.Event {
	instance.mutex.RLock()
	defer instance.mutex.RUnlock()

	if uint(len(instance.recorded)) <= index {
		panic(fmt.Sprintf("Index %d requested but the amount of recorded events is only %d", index, len(instance.recorded)))
	}

	return instance.recorded[index]
}

// GetAll returns all recorded events.
func (instance *Recorder) GetAll() []log.Event {
	instance.mutex.RLock()
	defer instance.mutex.RUnlock()

	result := make([]log.Event, len(instance.recorded))
	copy(result, instance.recorded)

	return result
}

// Reset will remove all recorded events of this Consumer.
func (instance *Recorder) Reset() {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()

	instance.recorded = []log.Event{}
}
