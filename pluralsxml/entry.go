package pluralsxml

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/redmaner/mixml/utils"
)

// Entry represents a plural in Android plurals.xml
type Entry struct {
	name  string
	items [][]string
}

// ParseEntry parses an Entry from a string. It returns true and the Entry if it was able
// to parse an Entry from the string. Otherwise it returns false and an empty item.
func (res *Resources) ParseEntry(base string) (bool, Entry) {

	// We do not parse empty strings
	if base == "" {
		return false, Entry{}
	}

	// We remove comments
	if strings.Contains(base, "<!--") {
		res.Comment = base + "\n"
		return false, Entry{}
	}

	var parsedFirstLine bool
	var strBuffer string
	var arrayName string
	var arrayItems [][]string

	sc := bufio.NewScanner(bytes.NewBufferString(base))

	for sc.Scan() {

		str := sc.Text()
		if !parsedFirstLine {
			parsedFirstLine = true

			// Trim prefix and suffix
			str = utils.TrimSpace(str)
			str = strings.TrimPrefix(str, `<plurals name="`)
			str = strings.TrimSuffix(str, `">`)
			arrayName = str
			continue
		}

		// The final line of the string is the closure of the entry.
		// We ignore this
		if strings.Contains(str, "</plurals>") {
			continue
		}

		if strings.Contains(str, "</item>") {
			if strBuffer != "" {
				str = strBuffer + "\n" + str
			}
			strBuffer = ""
			str = utils.TrimSpace(str)
			str = strings.TrimPrefix(str, `<item quantity="`)
			str = strings.TrimSuffix(str, "</item>")

			strSlice := strings.Split(str, `">`)
			quantity := strSlice[0]
			value := strSlice[1]

			arrayItems = append(arrayItems, []string{quantity, value})
			continue
		}
		if strBuffer != "" {
			strBuffer = strBuffer + "\n"
		}
		strBuffer = strBuffer + str
	}

	return true, Entry{
		name:  arrayName,
		items: arrayItems,
	}
}
