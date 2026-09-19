package fields

const minimumCompactionCost = 64

type compactedFields struct {
	values    map[string]any
	keys      []string
	hidden    keySet
	base      Fields
	compactAt int
}

func newDerivedLineage(target Fields, parent Fields) Fields {
	if isEmpty(parent) {
		return target
	}
	if isEmpty(target) {
		return parent
	}
	if !isCompactableLeaf(target) {
		return &lineage{target: target, parent: parent}
	}
	pendingCost, compactAt := compactionProgressOf(parent)
	result := &lineage{
		target:                target,
		parent:                parent,
		pendingCompactionCost: pendingCost + fieldsCost(target),
		compactAt:             compactAt,
	}
	if result.pendingCompactionCost >= result.compactAt {
		return compactFields(result)
	}
	return result
}

func newDerivedWithout(target Fields, keys ...string) Fields {
	if isEmpty(target) {
		return Empty()
	}
	if len(keys) == 0 {
		return target
	}
	excludedKeys := make(keySet, len(keys))
	for _, key := range keys {
		excludedKeys[key] = keyPresent
	}
	pendingCost, compactAt := compactionProgressOf(target)
	result := &without{
		fields:                target,
		excludedKeys:          excludedKeys,
		pendingCompactionCost: pendingCost + maxInt(1, len(excludedKeys)),
		compactAt:             compactAt,
	}
	if result.pendingCompactionCost >= result.compactAt {
		return compactFields(result)
	}
	return result
}

func compactionProgressOf(parent Fields) (int, int) {
	switch current := parent.(type) {
	case *lineage:
		if current != nil && current.pendingCompactionCost > 0 {
			return current.pendingCompactionCost, current.compactAt
		}
	case *without:
		if current != nil && current.pendingCompactionCost > 0 {
			return current.pendingCompactionCost, current.compactAt
		}
	case *compactedFields:
		if current != nil {
			return 0, current.compactAt
		}
	}
	return 0, minimumCompactionCost
}

func fieldsCost(value Fields) int {
	switch current := value.(type) {
	case mapped:
		return maxInt(1, len(current))
	case *mapped:
		if current != nil {
			return maxInt(1, len(*current))
		}
	}
	return 1
}

func isCompactableLeaf(value Fields) bool {
	switch value.(type) {
	case *single, mapped, *mapped:
		return true
	default:
		return false
	}
}

func compactFields(root Fields) Fields {
	values := map[string]any{}
	var keys []string
	seen := keySet{}
	current := root

	add := func(key string, value any) {
		if _, exists := seen[key]; exists {
			return
		}
		seen[key] = keyPresent
		values[key] = value
		keys = append(keys, key)
	}
	exclude := func(key string) {
		if _, exists := seen[key]; !exists {
			seen[key] = keyPresent
		}
	}
	addFields := func(value Fields) bool {
		switch source := value.(type) {
		case *single:
			if source != nil {
				add(source.key, source.value)
			}
			return true
		case mapped:
			for key, value := range source {
				add(key, value)
			}
			return true
		case *mapped:
			if source != nil {
				for key, value := range *source {
					add(key, value)
				}
			}
			return true
		}
		return false
	}

	for {
		switch source := current.(type) {
		case *lineage:
			if source == nil || !addFields(source.target) {
				return finishCompaction(values, keys, seen, current)
			}
			current = source.parent
		case *without:
			if source == nil {
				return finishCompaction(values, keys, seen, current)
			}
			for key := range source.excludedKeys {
				exclude(key)
			}
			current = source.fields
		case *compactedFields:
			if source == nil {
				return finishCompaction(values, keys, seen, nil)
			}
			for _, key := range source.keys {
				add(key, source.values[key])
			}
			for key := range source.hidden {
				exclude(key)
			}
			return finishCompaction(values, keys, seen, source.base)
		case *single, mapped, *mapped:
			_ = addFields(current)
			return finishCompaction(values, keys, seen, nil)
		case nil, *empty:
			return finishCompaction(values, keys, seen, nil)
		default:
			return finishCompaction(values, keys, seen, current)
		}
	}
}

func finishCompaction(values map[string]any, keys []string, seen keySet, base Fields) Fields {
	if len(values) == 0 && base == nil {
		return Empty()
	}
	var hidden keySet
	if base != nil && len(seen) > 0 {
		hidden = seen
	}
	stateSize := len(values)
	if base != nil {
		stateSize = len(hidden)
	}
	return &compactedFields{
		values:    values,
		keys:      keys,
		hidden:    hidden,
		base:      base,
		compactAt: maxInt(minimumCompactionCost, stateSize),
	}
}

func (instance *compactedFields) ForEach(consumer func(key string, value any) error) error {
	if instance == nil || consumer == nil {
		return nil
	}
	for _, key := range instance.keys {
		if err := consumer(key, instance.values[key]); err != nil {
			return err
		}
	}
	if instance.base == nil {
		return nil
	}
	seen := make(keySet, len(instance.hidden))
	for key := range instance.hidden {
		seen[key] = keyPresent
	}
	return instance.base.ForEach(func(key string, value any) error {
		if _, handled := seen[key]; handled {
			return nil
		}
		seen[key] = keyPresent
		return consumer(key, value)
	})
}

func (instance *compactedFields) Get(key string) (any, bool) {
	if instance == nil {
		return nil, false
	}
	if value, exists := instance.values[key]; exists {
		return value, true
	}
	if _, hidden := instance.hidden[key]; hidden || instance.base == nil {
		return nil, false
	}
	return instance.base.Get(key)
}

func (instance *compactedFields) Len() int {
	if instance == nil {
		return 0
	}
	result := len(instance.keys)
	if instance.base == nil {
		return result
	}
	seen := make(keySet, len(instance.hidden))
	for key := range instance.hidden {
		seen[key] = keyPresent
	}
	_ = forEachField(instance.base, func(key string, _ any) error {
		if _, handled := seen[key]; handled {
			return nil
		}
		seen[key] = keyPresent
		result++
		return nil
	}, true)
	return result
}

func (instance *compactedFields) With(key string, value any) Fields {
	return newDerivedLineage(With(key, value), instance)
}

func (instance *compactedFields) Withf(key string, format string, args ...any) Fields {
	return newDerivedLineage(Withf(key, format, args...), instance)
}

func (instance *compactedFields) WithAll(values map[string]any) Fields {
	return newDerivedLineage(WithAll(values), instance)
}

func (instance *compactedFields) Without(keys ...string) Fields {
	return newDerivedWithout(instance, keys...)
}

func maxInt(left, right int) int {
	if left > right {
		return left
	}
	return right
}
