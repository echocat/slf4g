package log

// UnwrapLogger unwraps a wrapped Logger inside of another Logger. For example
// by NewLoggerFacade.
func UnwrapLogger(in CoreLogger) Logger {
	if u, ok := in.(interface {
		Unwrap() Logger
	}); ok {
		return u.Unwrap()
	} else if u, ok := in.(interface {
		Unwrap() CoreLogger
	}); ok {
		if c := u.Unwrap(); c != nil {
			return NewLogger(c)
		}
	}
	return nil
}

// UnwrapCoreLogger unwraps a wrapped CoreLogger inside another CoreLogger.
// For example, by NewLoggerFacade.
func UnwrapCoreLogger(in CoreLogger) CoreLogger {
	if u, ok := in.(interface {
		Unwrap() Logger
	}); ok {
		return u.Unwrap()
	}
	if u, ok := in.(interface {
		Unwrap() CoreLogger
	}); ok {
		return u.Unwrap()
	}
	return nil
}

// UnwrapProvider unwraps a wrapped Provider inside another Provider. For
// example, by NewProviderFacade.
func UnwrapProvider(in Provider) Provider {
	if u, ok := in.(interface {
		Unwrap() Provider
	}); ok {
		return u.Unwrap()
	}
	return nil
}
