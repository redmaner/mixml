package miuires

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Resources contains MIUI resources. Resources are divided per file and file type
type Resources struct {
	FilePath string
	FileType string
	AppName  string
	Keys     []string
	Elements map[string]Elementer
	Comment  string
}

// NewResources returns new unloaded resources
func NewResources(filePath string) (res *Resources) {

	// Get app name
	var appName string
	slice := strings.Split(filePath, "/")
	for _, p := range slice {
		if strings.HasSuffix(p, ".apk") {
			appName = p
			break
		}
	}

	return &Resources{
		FilePath: filePath,
		FileType: filepath.Base(filePath),
		AppName:  appName,
		Keys:     []string{},
		Elements: make(map[string]Elementer),
	}
}

// Load loads the resources from res.FilePath
func (res *Resources) Load() (err error) {

	// Do XML integrity check first
	err = res.CheckIntegrity()
	if err != nil {
		return err
	}

	// Load the file
	f, err := os.Open(res.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// We scan the file with a bufio.Scanner. Each element item is stored in a slice.
	scanner := bufio.NewScanner(f)
	var elements []string
	var elementPlaceholder string

	for scanner.Scan() {

		// We scan each line of the file.
		str := scanner.Text()

		// If a string item has trash in it, we skip the line
		if strings.Contains(str, "resources>") || strings.Contains(str, "<?xml") {
			continue
		}

		// We want to join elements togehter
		if strings.Contains(str, "<string") || strings.Contains(str, "<array") || strings.Contains(str, "<string-array") || strings.Contains(str, "<integer-array") || strings.Contains(str, "<plurals") || strings.Contains(str, "<!--") {
			elements = append(elements, elementPlaceholder)
			elementPlaceholder = str
			continue
		}
		elementPlaceholder = elementPlaceholder + "\n" + str
	}
	elements = append(elements, elementPlaceholder)

	// We put every string in a map. This makes sure we have unique keys.
	// This way we remove double string items.
	for _, v := range elements {

		// Handle comment
		if strings.Contains(v, "<!--") && res.Comment == "" {
			res.Comment = v + "\n"
			continue
		}

		switch res.FileType {
		case FileTypeArrays:
			if ok, element := NewArrays(v); ok {
				res.Elements[element.GetName()] = element
			}
		case FileTypePlurals:
			if ok, element := NewPlurals(v); ok {
				res.Elements[element.GetName()] = element
			}
		case FileTypeStrings:
			if ok, element := NewStrings(v); ok {
				res.Elements[element.GetName()] = element
			}
		}
	}

	// We store xmlKeys in a separte slice and sort it, this way we can rebuild
	// the file in a ordered way.
	for k := range res.Elements {
		res.Keys = append(res.Keys, k)
	}
	sort.Strings(res.Keys)

	return nil
}

// Filter filters the resources using FilterConfig
func (res *Resources) Filter(fc *FilterConfig) error {

	for elementKey, element := range res.Elements {

		switch res.FileType {
		case FileTypeStrings:
			// Filter general key rules
			if rules, ok := fc.StringsKeyRules["all"]; ok {
				res.filterKey(rules, elementKey)
			}

			// Filter application key rules
			if rules, ok := fc.StringsKeyRules[res.AppName]; ok {
				res.filterKey(rules, elementKey)
			}

			// Filter general value rules
			if rules, ok := fc.StringsValueRules["all"]; ok {
				res.filterValue(rules, element.GetValue())
			}

			// Filter application value rules
			if rules, ok := fc.StringsValueRules[res.AppName]; ok {
				res.filterValue(rules, element.GetValue())
			}

		case FileTypeArrays:
			// Filter general key rules
			if rules, ok := fc.ArraysKeyRules["all"]; ok {
				res.filterKey(rules, elementKey)
			}

			// Filter application key rules
			if rules, ok := fc.ArraysKeyRules[res.AppName]; ok {
				res.filterKey(rules, elementKey)
			}

		case FileTypePlurals:
			// Filter general key rules
			if rules, ok := fc.PluralsKeyRules["all"]; ok {
				res.filterKey(rules, elementKey)
			}

			// Filter application key rules
			if rules, ok := fc.PluralsKeyRules[res.AppName]; ok {
				res.filterKey(rules, elementKey)
			}
		}
	}

	return nil
}

func (res *Resources) filterKey(rules []FilterRules, elementKey string) {
	for _, rule := range rules {
		switch rule.Mode {
		case FilterModeSuffix:
			if strings.HasSuffix(elementKey, rule.Match) {
				delete(res.Elements, elementKey)
			}
		case FilterModePrefix:
			if strings.HasPrefix(elementKey, rule.Match) {
				delete(res.Elements, elementKey)
			}
		case FilterModeContains:
			if strings.Contains(elementKey, rule.Match) {
				delete(res.Elements, elementKey)
			}
		}
	}
}

func (res *Resources) filterValue(rules []FilterRules, elementValue string) {
	for _, rule := range rules {
		switch rule.Mode {
		case FilterModeSuffix:
			if strings.HasSuffix(elementValue, rule.Match) {
				delete(res.Elements, elementValue)
			}
		case FilterModePrefix:
			if strings.HasSuffix(elementValue, rule.Match) {
				delete(res.Elements, elementValue)
			}
		case FilterModeContains:
			if strings.Contains(elementValue, rule.Match) {
				delete(res.Elements, elementValue)
			}
		default:
		}
	}
}

// Write writes resources to res.FilePath
func (res *Resources) Write() error {

	if len(res.Keys) == 0 {
		return errors.New("No resources to write")
	}

	f, err := os.Create(res.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	io.WriteString(f, fmt.Sprintf("<?xml version='1.0' encoding='UTF-8'?>\n"))

	if res.Comment != "" {
		io.WriteString(f, res.Comment)
	}

	io.WriteString(f, fmt.Sprintf("<resources>\n"))

	for _, key := range res.Keys {

		if val, ok := res.Elements[key]; ok {
			f.Write(val.Write())
		}
	}

	io.WriteString(f, fmt.Sprintf("</resources>\n"))
	return nil
}
