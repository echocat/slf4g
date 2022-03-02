package fields

// Exclude marks values of fields which might be present in the list
// of fields but should not be printed somewhere.
//
// Cause: `nil` values are always printed. This value is an exception.
var Exclude = struct{}{}
