package stringsxml

// Resources represents a set of strings in an android strings.xml
type Resources struct {
	FilePath  string
	Keys      []string
	Items     map[string]Item
	Format    bool
	ASCIIOnly bool
}

// NewResources returns an empty resources struct
func NewResources(filePath string, format bool, asciiOnly bool) *Resources {
	return &Resources{
		FilePath:  filePath,
		Keys:      []string{},
		Items:     make(map[string]Item),
		Format:    format,
		ASCIIOnly: asciiOnly,
	}
}
