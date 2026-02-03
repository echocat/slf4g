package formatter

// QuoteType defines how values should be quoted.
type QuoteType uint8

const (
	// QuoteTypeNormal defines that values are quoted how it is expected due to
	// their types; string="<string>", int=<int>, ...
	QuoteTypeNormal QuoteType = 0

	// QuoteTypeMinimal defines that values are only quoted if required. When it
	// is possible quoting will be prevented; foo="hello\"" bar=world
	QuoteTypeMinimal QuoteType = 1

	// QuoteTypeEverything forces to everything if required or not.
	QuoteTypeEverything QuoteType = 2
)

func stringNeedsQuoting(text string) bool {
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == ',' || ch == '_' || ch == ':' || ch == ';' || ch == '!' || ch == '?' ||
			ch == '/' || ch == '\\' ||
			ch == '@' || ch == '^' || ch == '+' || ch == '#' || ch == '~' ||
			ch == '(' || ch == ')' ||
			ch == '[' || ch == ']' ||
			ch == '{' || ch == '}') {
			return true
		}
	}
	return false
}
