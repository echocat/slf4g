//go:build go1.21

package sdk

import (
	sdk "log/slog"

	"github.com/echocat/slf4g/fields"
)

type attrs []sdk.Attr

const indexedAttrsThreshold = 64

func (instance attrs) ForEach(consumer func(key string, value any) error) error {
	if consumer == nil {
		return nil
	}
	for _, a := range instance {
		if err := consumer(a.Key, resolvedValueOf(a.Value).Any()); err != nil {
			return err
		}
	}
	return nil
}

func (instance attrs) Get(key string) (any, bool) {
	for _, a := range instance {
		if a.Key == key {
			return resolvedValueOf(a.Value).Any(), true
		}
	}
	return nil, false
}

func resolvedValueOf(value sdk.Value) sdk.Value {
	value = value.Resolve()
	if value.Kind() != sdk.KindGroup {
		return value
	}

	type frame struct {
		group       []sdk.Attr
		next        int
		parentIndex int
	}
	newFrame := func(value sdk.Value, parentIndex int) frame {
		originalGroup := value.Group()
		group := make([]sdk.Attr, len(originalGroup))
		copy(group, originalGroup)
		return frame{group: group, parentIndex: parentIndex}
	}

	stack := []frame{newFrame(value, 0)}
	for {
		current := &stack[len(stack)-1]
		if current.next < len(current.group) {
			index := current.next
			current.next++
			resolved := current.group[index].Value.Resolve()
			if resolved.Kind() == sdk.KindGroup {
				stack = append(stack, newFrame(resolved, index))
			} else {
				current.group[index].Value = resolved
			}
			continue
		}

		resolved := sdk.GroupValue(current.group...)
		parentIndex := current.parentIndex
		stack = stack[:len(stack)-1]
		if len(stack) == 0 {
			return resolved
		}
		stack[len(stack)-1].group[parentIndex].Value = resolved
	}
}

func (instance attrs) With(key string, value any) fields.Fields {
	return instance.asParentOf(fields.With(key, value))
}

func (instance attrs) Withf(key string, format string, args ...any) fields.Fields {
	return instance.asParentOf(fields.Withf(key, format, args...))
}

func (instance attrs) WithAll(of map[string]any) fields.Fields {
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

func (instance *attrs) add(keyPrefix string, vs ...sdk.Attr) {
	var indexes map[string]int
	if len(vs) >= indexedAttrsThreshold {
		indexes = make(map[string]int, len(*instance)+len(vs))
		for i, existing := range *instance {
			if _, exists := indexes[existing.Key]; !exists {
				indexes[existing.Key] = i
			}
		}
	}
	for _, v := range vs {
		nv := sdk.Attr{
			Key:   keyPrefix + v.Key,
			Value: v.Value,
		}
		if indexes != nil {
			if i, exists := indexes[nv.Key]; exists {
				(*instance)[i] = nv
			} else {
				indexes[nv.Key] = len(*instance)
				*instance = append(*instance, nv)
			}
			continue
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
