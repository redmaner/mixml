package miuires

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

// Entry represents an array in Android arrays.xml
type ElementArrays struct {
	name  string
	form  string
	items []string
}

func (ea *ElementArrays) GetName() (name string) {
	return ea.name
}

func (ea *ElementArrays) GetValue() (value string) {
	return ""
}

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
