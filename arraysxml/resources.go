package arraysxml

// Resources represents a set of strings in an android strings.xml
type Resources struct {
	FilePath  string
	Keys      []string
	Entries   map[string]Entry
	Comment   string
	Format    bool
	ASCIIOnly bool
}

// NewResources returns an empty resources struct
func NewResources(filePath string, format bool, asciiOnly bool) *Resources {
	return &Resources{
		FilePath:  filePath,
		Keys:      []string{},
		Entries:   make(map[string]Entry),
		Format:    format,
		ASCIIOnly: asciiOnly,
	}
}
