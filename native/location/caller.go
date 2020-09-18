package location

type CallerAware interface {
	Location

	GetFunction() string
	GetFile() string
	GetLine() int
}
