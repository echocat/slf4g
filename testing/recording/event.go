package recording

import (
	"fmt"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
)

type event struct {
	provider log.Provider
	fields   fields.Fields
	level    level.Level
}

func (instance *event) ForEach(consumer func(key string, value interface{}) error) error {
	return instance.fields.ForEach(consumer)
}

func (instance *event) Get(key string) (interface{}, bool) {
	return instance.fields.Get(key)
}

func (instance *event) Len() int {
	return instance.fields.Len()
}

func (instance *event) GetLevel() level.Level {
	return instance.level
}

func (instance *event) With(key string, value interface{}) log.Event {
	return instance.with(func(s fields.Fields) fields.Fields {
		return s.With(key, value)
	})
}

func (instance *event) Withf(key string, format string, args ...interface{}) log.Event {
	return instance.with(func(s fields.Fields) fields.Fields {
		return s.Withf(key, format, args...)
	})
}

func (instance *event) WithError(err error) log.Event {
	return instance.with(func(s fields.Fields) fields.Fields {
		return s.With(instance.provider.GetFieldKeysSpec().GetError(), err)
	})
}

func (instance *event) WithAll(of map[string]interface{}) log.Event {
	return instance.with(func(s fields.Fields) fields.Fields {
		return s.WithAll(of)
	})
}

func (instance *event) Without(keys ...string) log.Event {
	return instance.with(func(s fields.Fields) fields.Fields {
		return s.Without(keys...)
	})
}

func (instance *event) with(mod func(fields.Fields) fields.Fields) log.Event {
	return &event{
		provider: instance.provider,
		fields:   mod(instance.fields),
		level:    instance.level,
	}
}

func (instance *event) Format(f fmt.State, verb rune) {
	printf := func(format string, args ...interface{}) {
		_, _ = fmt.Fprintf(f, format, args...)
	}

	switch verb {
	case 'v':
		fds := instance.fields
		if f.Flag('+') && fds != nil && fds.Len() > 0 {
			printf("[%d] {", instance.level)
			_ = fds.ForEach(func(key string, value interface{}) error {
				if key == "logger" || key == "timestamp" {
					return nil
				}
				printf("\n\t%s=%+v", key, value)
				return nil
			})
			printf("\n}")
		} else {
			printf("[%d] {", instance.level)
			first := true
			_ = fds.ForEach(func(key string, value interface{}) error {
				if key == "logger" || key == "timestamp" {
					return nil
				}
				if first {
					first = false
				} else {
					printf(", ")
				}
				printf("%s=%v", key, value)
				return nil
			})
			printf("}")
		}
	default:
		printf("%%!%c(*types.Sym=%p)", verb, instance)
	}
}

func (instance *event) String() string {
	return fmt.Sprintf("%v", instance)
}
