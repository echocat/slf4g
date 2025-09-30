package sdk

import (
	sdk "log/slog"

	"github.com/echocat/slf4g/fields"
)

type attrs []sdk.Attr

func (instance attrs) ForEach(consumer func(key string, value interface{}) error) error {
	if consumer == nil {
		return nil
	}
	for _, a := range instance {
		if err := consumer(a.Key, a.Value.Any()); err != nil {
			return err
		}
	}
	return nil
}

func (instance attrs) Get(key string) (interface{}, bool) {
	for _, a := range instance {
		if a.Key == key {
			return a.Value.Any(), true
		}
	}
	return nil, false
}

func (instance attrs) With(key string, value interface{}) fields.Fields {
	return instance.asParentOf(fields.With(key, value))
}

func (instance attrs) Withf(key string, format string, args ...interface{}) fields.Fields {
	return instance.asParentOf(fields.Withf(key, format, args...))
}

func (instance attrs) WithAll(of map[string]interface{}) fields.Fields {
	return instance.asParentOf(fields.WithAll(of))
}

func (instance attrs) Without(keys ...string) fields.Fields {
	return fields.NewWithout(instance, keys...)
}

func (instance attrs) asParentOf(fds fields.Fields) fields.Fields {
	return fields.NewLineage(fds, instance)
}

func (instance attrs) Len() (result int) {
	return len(instance)
}

func (instance attrs) clone() attrs {
	result := make(attrs, len(instance))
	copy(result, instance)
	return result
}

func (instance *attrs) add(keyPrefix string, vs ...sdk.Attr) {
	for _, v := range vs {
		nv := sdk.Attr{
			Key:   keyPrefix + v.Key,
			Value: v.Value,
		}

		replaced := false
		for i, existing := range *instance {
			if existing.Key == nv.Key {
				(*instance)[i] = nv
				replaced = true
				break
			}
		}
		if !replaced {
			*instance = append(*instance, nv)
		}
	}
}
