package log

func Unwrap(in Logger) Logger {
	u, ok := in.(interface {
		Unwrap() Logger
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}

func UnwrapProvider(in Provider) Provider {
	u, ok := in.(interface {
		Unwrap() Provider
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}
