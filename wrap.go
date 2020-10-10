package log

// UnwrapLogger unwraps a wrapped Logger inside of another Logger. For example
// by NewLoggerFacade.
func UnwrapLogger(in Logger) Logger {
	type unwrapper interface {
		Unwrap() Logger
	}
	type coreUnwrapper interface {
		Unwrap() CoreLogger
	}
	if u, ok := in.(unwrapper); ok {
		return u.Unwrap()
	} else if u, ok := in.(coreUnwrapper); ok {
		if c := u.Unwrap(); c != nil {
			return NewLogger(c)
		}
	}
	return nil
}

// UnwrapCoreLogger unwraps a wrapped CoreLogger inside of another CoreLogger.
// For example by NewLoggerFacade.
func UnwrapCoreLogger(in CoreLogger) CoreLogger {
	type unwrapper interface {
		Unwrap() Logger
	}
	type coreUnwrapper interface {
		Unwrap() CoreLogger
	}
	if u, ok := in.(unwrapper); ok {
		return u.Unwrap()
	}
	if u, ok := in.(coreUnwrapper); ok {
		return u.Unwrap()
	}
	return nil
}

// UnwrapProvider unwraps a wrapped Provider inside of another Provider. For
// example by NewProviderFacade.
func UnwrapProvider(in Provider) Provider {
	type unwrapper interface {
		Unwrap() Provider
	}
	if u, ok := in.(unwrapper); ok {
		return u.Unwrap()
	}
	return nil
}
