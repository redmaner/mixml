package miuires

import (
	"bytes"
	"fmt"
	"strings"
)

// Entry represents a string in Android strings.xml
// <string name="example" formatted="false">Hello %s, this is an example of %s</string>
type ElementStrings struct {
	name          string
	value         string
	formatted     bool
	apostropheFix bool
}

func (es *ElementStrings) GetName() (name string) {
	return es.name
}

func (es *ElementStrings) GetValue() (value string) {

	// We fix apostrophe errors, by adding a \ in front of it
	// this slash is only added when it does not yet exist
	if es.apostropheFix && es.value[0] != '"' {
		es.apostropheFix = false
		var newValue string

		strSlice := strings.Split(es.value, "'")
		splits := len(strSlice)
		for i, v := range strSlice {

			if v == "" {
				continue
			}

			if i == splits-1 {
				newValue = newValue + v
				break
			}

			lastChar := len(v) - 1
			if v[lastChar] == 92 {
				newValue = newValue + v + "'"
				continue
			}
			newValue = newValue + v + `\'`
		}
		es.value = newValue
	}

	return es.value
}

func (es *ElementStrings) Parse(base string) (ok bool) {
	// We do not parse empty strings
	if base == "" {
		return false
	}

	// We remove comments
	if strings.Contains(base, "<!--") {
		return false
	}

	// Trim spaces
	base = trimSpace(base)

	// Trim prefix
	base = strings.TrimPrefix(base, "<string ")

	// Trim suffix
	base = strings.TrimSuffix(base, "</string>")

	// Get the name and value
	var baseSlice []string
	switch {
	case strings.Contains(base, ` formatted="false"`):
		baseSlice = strings.Split(base, `" formatted="false">`)
		es.formatted = true
	default:
		baseSlice = strings.Split(base, `">`)
	}
	es.name = strings.TrimPrefix(baseSlice[0], `name="`)
	es.value = baseSlice[1]

	// Determine if apostrophe's need to be fixed
	es.apostropheFix = strings.IndexByte(es.value, 39) >= 0

	// Determine if string needs to be formatted
	if strings.Count(es.value, "%s") >= 2 {
		es.formatted = true
	}
	if strings.Count(es.value, "%d") >= 2 {
		es.formatted = true
	}
	return true
}

func (es *ElementStrings) Write() []byte {
	var formatString string
	var writeString string

	w := bytes.NewBuffer([]byte{})

	if es.formatted {
		formatString = ` formatted="false"`
	}

	writeString = fmt.Sprintf(`    <string name="%s"%s>%s</string>`+"\n", es.name, formatString, es.GetValue())

	w.WriteString(writeString)

	return w.Bytes()
}
