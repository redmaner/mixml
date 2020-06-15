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

	// We do not parse empty strings
	if base == "" {
		return false, nil
	}

	// We remove comments
	if strings.Contains(base, "<!--") {
		return false, nil
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

			// Handle empty array
			if strings.Contains(str, "/>") {
				strSlice := strings.Split(str, `"`)
				ea.name = strSlice[0]
				return true, &ea
			}

			// Handle normal array
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
			ea.items = append(ea.items, fixApostrophe(str))
			continue
		}
		if strBuffer != "" {
			strBuffer = strBuffer + "\n"
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
		return []byte(fmt.Sprintf(`  <%s name="%s"/>`+"\n", ea.form, ea.name))
	}

	// Handle normal arrays
	w := bytes.NewBuffer([]byte{})
	buf := bytes.NewBufferString("")
	buf.WriteString(fmt.Sprintf(`  <%s name="%s">`+"\n", ea.form, ea.name))
	for _, item := range ea.items {
		buf.WriteString(fmt.Sprintf(`    <item>%s</item>`+"\n", item))
	}
	buf.WriteString(fmt.Sprintf(`  </%s>`+"\n", ea.form))
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

	// We do not parse empty strings
	if base == "" {
		return false, nil
	}

	// We remove comments
	if strings.Contains(base, "<!--") {
		return false, nil
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
			value := fixApostrophe(strSlice[1])

			ep.items = append(ep.items, value)
			ep.quantities = append(ep.quantities, quantity)
			continue
		}
		if strBuffer != "" {
			strBuffer = strBuffer + "\n"
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
	buf.WriteString(fmt.Sprintf(`  <plurals name="%s">`+"\n", ep.name))
	for index, item := range ep.items {
		buf.WriteString(fmt.Sprintf(`    <item quantity="%s">%s</item>`+"\n", item, ep.quantities[index]))
	}
	buf.WriteString(fmt.Sprintf(`  </plurals>` + "\n"))
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

	// We do not parse empty strings
	if base == "" {
		return false, nil
	}

	// We remove comments
	if strings.Contains(base, "<!--") {
		return false, nil
	}

	// Trim spaces
	base = trimSpace(base)

	// Trim prefix
	base = strings.TrimPrefix(base, "<string ")

	// Handle empty strings
	if strings.Contains(base, `"/>`) || strings.Contains(base, `" />`) {
		baseSlice := strings.Split(base, `name="`)
		baseSlice = strings.Split(baseSlice[1], `"`)
		es.name = baseSlice[0]
		return true, &es
	}

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
	es.value = fixApostrophe(baseSlice[1])

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
		return []byte(fmt.Sprintf(`  <string name="%s"/>`, es.name) + "\n")
	}

	// Handle normal strings
	var formatString string
	if es.formatted {
		formatString = ` formatted="false"`
	}
	return []byte(fmt.Sprintf(`  <string name="%s"%s>%s</string>`+"\n", es.name, formatString, es.GetValue()))
}
