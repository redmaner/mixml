package miuires

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

// Entry represents a plural in Android plurals.xml
type ElementPlurals struct {
	name  string
	items [][]string
}

func (ep *ElementPlurals) GetName() (name string) {
	return ep.name
}

func (ep *ElementPlurals) GetValue() (value string) {
	return ""
}

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
