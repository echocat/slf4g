package level

// Provider provides all available Levels.
type Provider interface {
	GetName() string
	GetLevels() Levels
}
