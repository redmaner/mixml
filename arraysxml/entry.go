package arraysxml

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/redmaner/mixml/utils"
)

// Entry represents an array in Android arrays.xml
type Entry struct {
	name  string
	form  string
	items []string
}

// ParseEntry parses an Entry from a string. It returns true and the Entry if it was able
// to parse an Entry from the string. Otherwise it returns false and an empty item.
func ParseEntry(base string, format bool, asciiOnly bool) (bool, Entry) {

	if base == "" {
		return false, Entry{}
	}

	var parsedFirstLine bool
	var strBuffer string
	var arrayForm string
	var arrayName string
	var arrayItems []string

	sc := bufio.NewScanner(bytes.NewBufferString(base))

	for sc.Scan() {

		str := sc.Text()
		if !parsedFirstLine {
			parsedFirstLine = true

			switch {

			case strings.Contains(str, "<array"):
				arrayForm = "array"

			case strings.Contains(str, "<string-array"):
				arrayForm = "string-array"

			case strings.Contains(str, "<integer-array"):
				arrayForm = "integer-array"
			}

			// Trim prefix and suffix
			str = utils.TrimSpace(str)
			str = strings.TrimPrefix(str, fmt.Sprintf(`<%s name="`, arrayForm))
			str = strings.TrimSuffix(str, `">`)
			arrayName = str
		}

		// The final line of the string is the closure of the entry.
		// We ignore this
		if strings.Contains(str, arrayForm) {
			continue
		}

		if strings.Contains(str, "<item></item>") {
			arrayItems = append(arrayItems, "")
			continue
		}

		if strings.Contains(str, "</item>") {
			if strBuffer != "" {
				str = strBuffer + "\n" + str
			}
			strBuffer = ""
			str = utils.TrimSpace(str)
			str = strings.TrimPrefix(str, "<item>")
			str = strings.TrimSuffix(str, "</item>")
			arrayItems = append(arrayItems, str)
			continue
		}
		if strBuffer != "" {
			strBuffer = strBuffer + "\n"
		}
		strBuffer = strBuffer + str
	}

	return true, Entry{
		name:  arrayName,
		form:  arrayForm,
		items: arrayItems,
	}
}
