package fields

type traversalFrame struct {
	fields       Fields
	excludedKeys keySet
}

func forEachField(root Fields, consumer func(key string, value any) error, ignoreSourceErrors bool) error {
	if root == nil || consumer == nil {
		return nil
	}

	var handledKeys keySet
	var excludedKeys map[string]int
	stack := []traversalFrame{{fields: root}}
	for len(stack) > 0 {
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]

		if current.excludedKeys != nil {
			for key := range current.excludedKeys {
				excludedKeys[key]--
				if excludedKeys[key] == 0 {
					delete(excludedKeys, key)
				}
			}
			continue
		}

		switch currentFields := current.fields.(type) {
		case *lineage:
			if currentFields == nil {
				continue
			}
			if handledKeys == nil {
				handledKeys = keySet{}
			}
			stack = append(stack,
				traversalFrame{fields: currentFields.parent},
				traversalFrame{fields: currentFields.target},
			)
		case *without:
			if currentFields == nil {
				continue
			}
			if len(currentFields.excludedKeys) > 0 {
				if excludedKeys == nil {
					excludedKeys = make(map[string]int, len(currentFields.excludedKeys))
				}
				for key := range currentFields.excludedKeys {
					excludedKeys[key]++
				}
				stack = append(stack, traversalFrame{excludedKeys: currentFields.excludedKeys})
			}
			stack = append(stack, traversalFrame{fields: currentFields.fields})
		default:
			if current.fields == nil {
				continue
			}
			err := current.fields.ForEach(func(key string, value any) error {
				if excludedKeys[key] > 0 {
					return nil
				}
				if handledKeys != nil {
					if _, handled := handledKeys[key]; handled {
						return nil
					}
					handledKeys[key] = keyPresent
				}
				return consumer(key, value)
			})
			if err != nil && !ignoreSourceErrors {
				return err
			}
		}
	}
	return nil
}

func getField(root Fields, key string) (any, bool) {
	stack := []Fields{root}
	for len(stack) > 0 {
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]

		switch currentFields := current.(type) {
		case *lineage:
			if currentFields == nil {
				continue
			}
			stack = append(stack, currentFields.parent, currentFields.target)
		case *without:
			if currentFields == nil {
				continue
			}
			if _, excluded := currentFields.excludedKeys[key]; !excluded {
				stack = append(stack, currentFields.fields)
			}
		default:
			if current == nil {
				continue
			}
			if value, exists := current.Get(key); exists {
				return value, true
			}
		}
	}
	return nil, false
}

func fieldCount(root Fields) (result int) {
	_ = forEachField(root, func(string, any) error {
		result++
		return nil
	}, true)
	return
}
