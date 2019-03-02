package stringsxml

import (
	"strings"
	"unicode/utf8"

	"github.com/redmaner/mixml/utils"
)

// Item represents a string in Android strings.xml
// <string name="example" formatted="false">Hello %s, this is an example of %s</string>
type Item struct {
	name          string
	value         string
	formatted     bool
	apostropheFix bool
}

// ParseItem parses a Item from a string. It returns true and the Item if it was able
// to parse an Item from the string. Otherwise it returns false and an empty item.
func ParseItem(base string, format bool, asciiOnly bool) (bool, Item) {

	if base == "" {
		return false, Item{}
	}

	// Trim spaces
	base = utils.TrimSpace(base)

	// Trim prefix
	base = strings.TrimPrefix(base, "<string ")

	// Trim suffix
	base = strings.TrimSuffix(base, "</string>")

	// Get the name and value
	var baseSlice []string
	switch {
	case strings.Contains(base, ` formatted="false"`):
		baseSlice = strings.Split(base, `" formatted="false">`)
	default:
		baseSlice = strings.Split(base, `">`)
	}
	name := strings.TrimPrefix(baseSlice[0], `name="`)
	value := baseSlice[1]

	if format {

		// If value contains multiple _ and doesn't contain spaces we skip it
		if strings.Count(value, "_") >= 2 && strings.Count(value, " ") == 0 {
			return false, Item{}
		}

		// If value contains multiple . and doesn't contain spaces we skip it
		if strings.Count(value, ".") > 2 && strings.Count(value, " ") == 0 {
			return false, Item{}
		}
	}

	if asciiOnly {
		lenValue := len(value) - 1
		switch {
		case lenValue == -1, lenValue == 0:
			// We do nothing
		default:
			testOne, _ := utf8.DecodeRune([]byte{value[0]})
			testTwo, _ := utf8.DecodeRune([]byte{value[lenValue]})
			if testOne > 591 && testTwo > 591 {
				return false, Item{}
			}
		}
	}

	// Determine if apostrophe's need to be fixed
	apostropheFix := strings.IndexByte(value, 39) >= 0

	// Determine if string needs to be formatted
	var formatted bool
	if strings.Count(value, "%s") >= 2 {
		formatted = true
	}
	if strings.Count(value, "%d") >= 2 {
		formatted = true
	}

	return true, Item{
		name:          name,
		value:         value,
		apostropheFix: apostropheFix,
		formatted:     formatted,
	}
}
