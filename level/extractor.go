package level

// LineExtractor checks a given line for the contained Level and will return it.
// The implementation can decide to also return another level.
type LineExtractor interface {
	ExtractLevelFromLine([]byte) (Level, error)
}

// LineExtractorFunc like LineExtractor but as a func.
type LineExtractorFunc func([]byte) (Level, error)

func (instance LineExtractorFunc) ExtractLevelFromLine(in []byte) (Level, error) {
	return instance(in)
}

// FixedLevelExtractor is an implementation of LineExtractor which always
// returns the same Level value regardless what was in the line.
func FixedLevelExtractor(in Level) LineExtractor {
	return fixedLevelExtractor(in)
}

type fixedLevelExtractor Level

func (instance fixedLevelExtractor) ExtractLevelFromLine([]byte) (Level, error) {
	return Level(instance), nil
}
