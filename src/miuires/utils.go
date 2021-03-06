package miuires

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// trimSpace removes all space prefixes and suffixes from a string
func trimSpace(base string) string {
	if base == "" {
		return base
	}

	// Trim tabs and spaces as prefixes
	for len(base) > 0 && base[0] == '\t' || len(base) > 0 && base[0] == ' ' || len(base) > 0 && base[0] == '\r' || len(base) > 0 && base[0] == '\n' || len(base) > 0 && base[0] == '+' {
		base = base[1:]
	}

	// Trim tabs and spaces as suffixes
	for len(base) > 0 && base[len(base)-1] == '\t' || len(base) > 0 && base[len(base)-1] == ' ' || len(base) > 0 && base[len(base)-1] == '\r' || len(base) > 0 && base[len(base)-1] == '\n' {
		base = base[:len(base)-1]
	}
	return base
}

func fixApostrophe(base string) (fixed string) {

	// If there are no apostrophes, return base
	apostropheIndex := strings.IndexByte(base, 39)
	if apostropheIndex < 0 {
		return base
	}

	// If strings are encapsulated with quotes, return base
	if base[0] == '"' {
		return base
	}

	// If apostrophe is escaped with backslash, return base
	if apostropheIndex > 0 && base[apostropheIndex-1] == 92 {
		return base
	}

	// We fix the apostrophe's by escaping it with a backslash
	splits := strings.Split(base, "'")
	firstParsed := false
	for _, split := range splits {
		if !firstParsed {
			firstParsed = true
			fixed = split
			continue
		}
		fixed = fixed + `\'` + split
	}

	return fixed
}

// getElementParameter extracts elementer parameters, like the name of a string or
// the quantity of a plural
func getElementParameter(base string, parameter string) (value string) {
	parameter = parameter + `="`
	start := strings.Index(base, parameter)
	if start < 0 {
		return ""
	}

	end := strings.IndexByte(base[start+len(parameter):], '"')
	end += start + len(parameter)
	value = strings.TrimPrefix(base[start:end], parameter)
	return
}

// getElementValue extract the elemement value
func getElementValue(base string, suffix string) (value string) {
	startValue := strings.IndexByte(base, '>')
	value = base[startValue+1:]
	value = strings.TrimSuffix(value, suffix)
	value = fixApostrophe(value)
	return
}

// getPluralsItem extracts the value and quantity of a plural item
func getPluralsItem(base string) (quantity string, value string) {
	quantity = getElementParameter(base, "quantity")
	value = getElementValue(base, "</item>")
	return
}

// getStringsNameValue extract the name and value of strings
func getStringsNameValue(base string) (name string, value string, formatted bool) {
	formatted = strings.Contains(base, `formatted="false"`)
	name = getElementParameter(base, "name")
	value = getElementValue(base, "</string>")
	return
}

// CheckIntegrity checks basic XML integrity of strings.xml and arrays.xml
func (res *Resources) CheckIntegrity() (err error) {

	// Open file, defer closure
	f, err := os.Open(res.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Read the file in a slice of bytes
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	// Count the occurrence of <resources></resources> pairs. These should be equal.
	resourceOpenCount := bytes.Count(data, []byte("<resources>"))
	resourceCloseCount := bytes.Count(data, []byte("</resources>"))

	if resourceOpenCount != resourceCloseCount {
		return fmt.Errorf("%s: basic XML integrity check failed (resources)", res.FilePath)
	}

	switch res.FileType {

	// strings.xml checks
	case FileTypeStrings:

		// Count the occurrence of <string></string> pairs. These should be equal.
		stringOpenCount := bytes.Count(data, []byte("<string name="))
		stringCloseCount := bytes.Count(data, []byte("</string>"))
		stringCloseCount = stringCloseCount + bytes.Count(data, []byte(`"/>`))
		stringCloseCount = stringCloseCount + bytes.Count(data, []byte(`" />`))

		if stringOpenCount != stringCloseCount {
			return fmt.Errorf("%s: basic XML integrity check failed (string mismatch)", res.FilePath)
		}

		// arrays.xml checks
	case FileTypeArrays:

		// Count the occurrence of :
		// * <array></array> pairs
		// * <string-array></string-array> pairs
		// * <integer-array></integer-array> pairs
		// These should all be equal
		arrayOpenCount := bytes.Count(data, []byte("<array name="))
		arrayOpenCount = arrayOpenCount + bytes.Count(data, []byte("<string-array name="))
		arrayOpenCount = arrayOpenCount + bytes.Count(data, []byte("<integer-array name="))
		arrayCloseCount := bytes.Count(data, []byte("</array>"))
		arrayCloseCount = arrayCloseCount + bytes.Count(data, []byte("</string-array>"))
		arrayCloseCount = arrayCloseCount + bytes.Count(data, []byte("</integer-array>"))
		arrayCloseCount = arrayCloseCount + bytes.Count(data, []byte(`"/>`))
		arrayCloseCount = arrayCloseCount + bytes.Count(data, []byte(`" />`))

		if arrayOpenCount != arrayCloseCount {
			return fmt.Errorf("%s: basic XML integrity check failed (array mismatch)", res.FilePath)
		}

		// Count the occurrence of <item></item> pairs. These should be equal.
		itemOpenCount := bytes.Count(data, []byte("<item>"))
		itemCloseCount := bytes.Count(data, []byte("</item>"))

		if itemOpenCount != itemCloseCount {
			return fmt.Errorf("%s: basic XML integrity check failed (item mismatch)", res.FilePath)
		}

		// Checks for plurals.xml
	case FileTypePlurals:

		// Count the occurrence of <plurals></plurals> pairs. These should be equal.
		pluralsOpenCount := bytes.Count(data, []byte("<plurals name="))
		pluralsCloseCount := bytes.Count(data, []byte("</plurals>"))

		if pluralsOpenCount != pluralsCloseCount {
			return fmt.Errorf("%s: basic XML integrity check failed (plurals mismatch)", res.FilePath)
		}

		// Count the occurrence of <item></item> pairs. These should be equal.
		itemOpenCount := bytes.Count(data, []byte("<item"))
		itemCloseCount := bytes.Count(data, []byte("</item>"))

		if itemOpenCount != itemCloseCount {
			return fmt.Errorf("%s: basic XML integrity check failed (item mismatch)", res.FilePath)
		}
	}

	return nil
}
