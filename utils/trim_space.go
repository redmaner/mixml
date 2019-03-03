package utils

// TrimSpace removes all space prefixes and suffixes from a string
func TrimSpace(base string) string {
	if base == "" {
		return base
	}

	// Trim tabs and spaces as prefixes
	for base[0] == '\t' || base[0] == ' ' {
		base = base[1:]
	}

	// Trim tabs and spaces as suffixes
	for base[len(base)-1] == '\t' || base[len(base)-1] == ' ' {
		base = base[:len(base)-1]
	}
	return base
}
