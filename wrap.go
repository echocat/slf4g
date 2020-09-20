package log

func Unwrap(in Logger) Logger {
	type unwrapper interface {
		Unwrap() Logger
	}
	type coreUnwrapper interface {
		UnwrapCore() CoreLogger
	}
	if u, ok := in.(unwrapper); ok {
		return u.Unwrap()
	} else if u, ok := in.(coreUnwrapper); ok {
		if c := u.UnwrapCore(); c != nil {
			return NewLogger(c)
		}
	}
	return nil
}

func UnwrapCore(in CoreLogger) CoreLogger {
	type unwrapper interface {
		UnwrapCore() CoreLogger
	}
	if u, ok := in.(unwrapper); ok {
		return u.UnwrapCore()
	}
	return nil
}

func UnwrapProvider(in Provider) Provider {
	type unwrapper interface {
		UnwrapProvider() Provider
	}
	if u, ok := in.(unwrapper); ok {
		return u.UnwrapProvider()
	}
	return nil
}
