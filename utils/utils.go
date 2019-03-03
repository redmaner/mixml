package utils

// TrimSpacePrefix removes all space prefixes and suffixes from a string
func TrimSpace(base string) string {
	if base == "" {
		return base
	}
	for base[0] == '\t' {
		base = base[1:]
	}
	for base[0] == ' ' {
		base = base[1:]
	}
	for base[len(base)-1] == ' ' {
		base = base[:len(base)-1]
	}
	return base
}
