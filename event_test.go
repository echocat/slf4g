package log

type entries []entry

func (instance *entries) add(key string, value interface{}) {
	*instance = append(*instance, entry{key, value})
}

func (instance *entries) consumer() func(key string, value interface{}) error {
	return func(key string, value interface{}) error {
		instance.add(key, value)
		return nil
	}
}

type entry struct {
	key   string
	value interface{}
}
