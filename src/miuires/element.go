package miuires

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

// Elementer is an interface that holds common behavior for MIUI resources
type Elementer interface {
	GetName() (name string)
	GetValue() (value string)
	Parse(base string) (ok bool)
	Write() []byte
}

// ElementArrays implements the Elementer interface, and holds information and behavior
// to handle MIUI arrays.xml
type ElementArrays struct {
	name  string
	form  string
	items []string
}

// GetName returns the name (key) of the arrays element
func (ea *ElementArrays) GetName() (name string) {
	return ea.name
}

// GetValue returns the value (body) of the arrays element
func (ea *ElementArrays) GetValue() (value string) {
	return ""
}

// Parse parses a string and converts it into an arrays element if possible
// It returns ok if it was succesful
func (ea *ElementArrays) Parse(base string) (ok bool) {

	// We do not parse empty strings
	if base == "" {
		return false
	}

	// We remove comments
	if strings.Contains(base, "<!--") {
		return false
	}

	var parsedFirstLine bool
	var strBuffer string

	scanner := bufio.NewScanner(bytes.NewBufferString(base))

	for scanner.Scan() {

		str := scanner.Text()
		if !parsedFirstLine {
			parsedFirstLine = true

			switch {

			case strings.Contains(str, "<array"):
				ea.form = "array"

			case strings.Contains(str, "<string-array"):
				ea.form = "string-array"

			case strings.Contains(str, "<integer-array"):
				ea.form = "integer-array"
			}

			// Trim prefix and suffix
			str = trimSpace(str)
			str = strings.TrimPrefix(str, fmt.Sprintf(`<%s name="`, ea.form))
			str = strings.TrimSuffix(str, `">`)
			ea.name = str
			continue
		}

		// The final line of the string is the closure of the entry.
		// We ignore this
		if strings.Contains(str, ea.form) {
			continue
		}

		if strings.Contains(str, "<item></item>") {
			ea.items = append(ea.items, "")
			continue
		}

		if strings.Contains(str, "</item>") {
			if strBuffer != "" {
				str = strBuffer + "\n" + str
			}
			strBuffer = ""
			str = trimSpace(str)
			str = strings.TrimPrefix(str, "<item>")
			str = strings.TrimSuffix(str, "</item>")
			ea.items = append(ea.items, str)
			continue
		}
		if strBuffer != "" {
			strBuffer = strBuffer + "\n"
		}
		strBuffer = strBuffer + str
	}

	return true
}

// Write writes the contents of the arrays element to a slice of bytes
func (ea *ElementArrays) Write() []byte {
	w := bytes.NewBuffer([]byte{})
	buf := bytes.NewBufferString("")
	buf.WriteString(fmt.Sprintf(`    <%s name="%s">`+"\n", ea.form, ea.name))
	for _, item := range ea.items {
		buf.WriteString(fmt.Sprintf(`        <item>%s</item>`+"\n", item))
	}
	buf.WriteString(fmt.Sprintf(`    </%s>`+"\n", ea.form))
	w.WriteString(buf.String())
	return w.Bytes()
}

// ElementPlurals implements the Elementer interface, and holds information and behavior
// to handle MIUI plurals.xml
type ElementPlurals struct {
	name  string
	items [][]string
}

// GetName returns the name (key) of the plurals element
func (ep *ElementPlurals) GetName() (name string) {
	return ep.name
}

// GetValue returns the value (body) of the plurals element
func (ep *ElementPlurals) GetValue() (value string) {
	return ""
}

// Parse parses a string and converts it into an plurals element if possible
// It returns ok if it was succesful
func (ep *ElementPlurals) Parse(base string) (ok bool) {

	// We do not parse empty strings
	if base == "" {
		return false
	}

	// We remove comments
	if strings.Contains(base, "<!--") {
		return false
	}

	var parsedFirstLine bool
	var strBuffer string

	scanner := bufio.NewScanner(bytes.NewBufferString(base))

	for scanner.Scan() {

		str := scanner.Text()
		if !parsedFirstLine {
			parsedFirstLine = true

			// Trim prefix and suffix
			str = trimSpace(str)
			str = strings.TrimPrefix(str, `<plurals name="`)
			str = strings.TrimSuffix(str, `">`)
			ep.name = str
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
			str = trimSpace(str)
			str = strings.TrimPrefix(str, `<item quantity="`)
			str = strings.TrimSuffix(str, "</item>")

			strSlice := strings.Split(str, `">`)
			quantity := strSlice[0]
			value := strSlice[1]

			ep.items = append(ep.items, []string{quantity, value})
			continue
		}
		if strBuffer != "" {
			strBuffer = strBuffer + "\n"
		}
		strBuffer = strBuffer + str
	}

	return true
}

// Write writes the contents of the plurals element to a slice of bytes
func (ep *ElementPlurals) Write() []byte {
	w := bytes.NewBuffer([]byte{})
	buf := bytes.NewBufferString("")

	buf.WriteString(fmt.Sprintf(`    <plurals name="%s">`+"\n", ep.name))

	for _, item := range ep.items {
		buf.WriteString(fmt.Sprintf(`        <item quantity="%s">%s</item>`+"\n", item[0], item[1]))
	}
	buf.WriteString(fmt.Sprintf(`    </plurals>` + "\n"))

	w.WriteString(buf.String())
	return w.Bytes()
}

// ElementStrings implements the Elementer interface, and holds information and behavior
// to handle MIUI strings.xml
type ElementStrings struct {
	name          string
	value         string
	formatted     bool
	apostropheFix bool
}

// GetName returns the name (key) of the strings element
func (es *ElementStrings) GetName() (name string) {
	return es.name
}

// GetValue returns the value (body) of the strings element
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

// Parse parses a string and converts it into a strings element if possible
// It returns ok if it was succesful
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

// Write writes the contents of the element strings to a slice of bytes
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
