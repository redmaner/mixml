package miuires

import (
	"bytes"
	"fmt"
	"strings"
)

// Elementer is an interface that holds common behavior for MIUI resources
type Elementer interface {
	GetName() (name string)
	GetItems() (items []string)
	GetValue() (value string)
	Write() []byte
}

// ElementArrays implements the Elementer interface, and holds information and behavior
// to handle MIUI arrays.xml
type ElementArrays struct {
	name  string
	form  string
	items []string
}

// NewArrays parses a string and converts it into an arrays element if possible
// It returns ok if it was succesful, and a pointer to the new ElementArrays
func NewArrays(base string) (bool, *ElementArrays) {

	var ea ElementArrays

	// We remove comments
	if strings.Contains(base, "<!--") {
		return false, nil
	}

	var parsedFirstLine bool
	var strBuffer string

	// We do not parse empty arrays
	base = trimSpace(base)
	if base == "" {
		return false, nil
	}

	scanner := NewElementScanner(bytes.NewBufferString(base))

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

			// extract name
			ea.name = getElementParameter(str, "name")

			// If array is empty we break out, otherwise we continue
			if strings.Contains(str, "/>") {
				return true, &ea
			}
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
				str = strBuffer + str
			}
			strBuffer = ""
			str = trimSpace(str)
			ea.items = append(ea.items, getElementValue(str, "</item>"))
			continue
		}
		strBuffer = strBuffer + str
	}

	return true, &ea
}

// GetName returns the name (key) of the arrays element
func (ea *ElementArrays) GetName() (name string) {
	return ea.name
}

// GetItems returns the items of the arrays element
func (ea *ElementArrays) GetItems() (items []string) {
	return ea.items
}

// GetValue returns the value (body) of the arrays element
func (ea *ElementArrays) GetValue() (value string) {
	return ""
}

// Write writes the contents of the arrays element to a slice of bytes
func (ea *ElementArrays) Write() []byte {

	// Handle empty array
	if len(ea.items) == 0 {
		return []byte(fmt.Sprintf(`    <%s name="%s"/>`+"\n", ea.form, ea.name))
	}

	// Handle normal arrays
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
	name       string
	items      []string
	quantities []string
}

// NewPlurals parses a string and converts it into an plurals element if possible
// It returns ok if it was succesful, and a pointer to the new ElementPlurals
func NewPlurals(base string) (bool, *ElementPlurals) {

	var ep ElementPlurals

	// We remove comments
	if strings.Contains(base, "<!--") {
		return false, nil
	}

	var parsedFirstLine bool
	var strBuffer string

	// We do not parse empty plurals
	base = trimSpace(base)
	if base == "" {
		return false, nil
	}

	scanner := NewElementScanner(bytes.NewBufferString(base))

	for scanner.Scan() {

		str := scanner.Text()
		if !parsedFirstLine {
			parsedFirstLine = true

			// Trim prefix and suffix
			str = trimSpace(str)
			ep.name = getElementParameter(str, "name")
			continue
		}

		// The final line of the string is the closure of the entry.
		// We ignore this
		if strings.Contains(str, "</plurals>") {
			continue
		}

		if strings.Contains(str, "</item>") {
			if strBuffer != "" {
				str = strBuffer + str
			}
			strBuffer = ""
			str = trimSpace(str)
			quantity, value := getPluralsItem(str)
			ep.items = append(ep.items, value)
			ep.quantities = append(ep.quantities, quantity)
			continue
		}
		strBuffer = strBuffer + str
	}

	return true, &ep
}

// GetName returns the name (key) of the plurals element
func (ep *ElementPlurals) GetName() (name string) {
	return ep.name
}

// GetItems returns the items of the plurals element
func (ep *ElementPlurals) GetItems() (items []string) {
	return ep.items
}

// GetValue returns the value (body) of the plurals element
func (ep *ElementPlurals) GetValue() (value string) {
	return ""
}

// Write writes the contents of the plurals element to a slice of bytes
func (ep *ElementPlurals) Write() []byte {
	w := bytes.NewBuffer([]byte{})
	buf := bytes.NewBufferString("")
	buf.WriteString(fmt.Sprintf(`    <plurals name="%s">`+"\n", ep.name))
	for index, item := range ep.items {
		buf.WriteString(fmt.Sprintf(`        <item quantity="%s">%s</item>`+"\n", ep.quantities[index], item))
	}
	buf.WriteString(fmt.Sprintf(`    </plurals>` + "\n"))
	w.WriteString(buf.String())
	return w.Bytes()
}

// ElementStrings implements the Elementer interface, and holds information and behavior
// to handle MIUI strings.xml
type ElementStrings struct {
	name      string
	value     string
	formatted bool
}

// NewStrings parses a string and converts it into a strings element if possible
// It returns ok if it was succesful, and a pointer to the new ElementStrings
func NewStrings(base string) (bool, *ElementStrings) {

	var es ElementStrings

	// We remove comments
	if strings.Contains(base, "<!--") {
		return false, nil
	}

	// Trim spaces
	base = trimSpace(base)

	// We do not parse empty strings
	if base == "" {
		return false, nil
	}

	// Trim prefix
	base = strings.TrimPrefix(base, "<string ")

	// Handle empty strings
	if strings.Contains(base, `"/>`) || strings.Contains(base, `" />`) {
		baseSlice := strings.Split(base, `name="`)
		baseSlice = strings.Split(baseSlice[1], `"`)
		es.name = baseSlice[0]
		return true, &es
	}

	// Get the name and value
	es.name, es.value, es.formatted = getStringsNameValue(base)

	// Determine if string needs to be formatted
	if strings.Count(es.value, "%s") >= 2 {
		es.formatted = true
	}
	if strings.Count(es.value, "%d") >= 2 {
		es.formatted = true
	}
	return true, &es
}

// GetName returns the name (key) of the strings element
func (es *ElementStrings) GetName() (name string) {
	return es.name
}

// GetItems returns the items of the strings element
func (es *ElementStrings) GetItems() (items []string) {
	return []string{}
}

// GetValue returns the value (body) of the strings element
func (es *ElementStrings) GetValue() (value string) {
	return es.value
}

// Write writes the contents of the element strings to a slice of bytes
func (es *ElementStrings) Write() []byte {

	// Handle empty strings
	if es.value == "" {
		return []byte(fmt.Sprintf(`    <string name="%s"/>`, es.name) + "\n")
	}

	// Handle normal strings
	var formatString string
	if es.formatted {
		formatString = ` formatted="false"`
	}
	return []byte(fmt.Sprintf(`    <string name="%s"%s>%s</string>`+"\n", es.name, formatString, es.GetValue()))
}
