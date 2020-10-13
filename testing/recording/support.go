package recording

// BeTrue is a utility function with returns a pointer to a true bool to
// easier use with Provider.SetIfAbsent and/or CoreLogger.SetIfAbsent.
func BeTrue() *bool {
	v := true
	return &v
}

// BeFalse is a utility function with returns a pointer to a false bool to
// easier use with Provider.SetIfAbsent and/or CoreLogger.SetIfAbsent.
func BeFalse() *bool {
	v := false
	return &v
}
