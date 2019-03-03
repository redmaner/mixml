// Copyright 2019 Jake van der Putten.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stringsxml

import (
	"strings"
	"unicode/utf8"

	"github.com/redmaner/mixml/utils"
)

// Entry represents a string in Android strings.xml
// <string name="example" formatted="false">Hello %s, this is an example of %s</string>
type Entry struct {
	name          string
	value         string
	formatted     bool
	apostropheFix bool
}

// ParseEntry parses an Entry from a string. It returns true and the Entry if it was able
// to parse an Entry from the string. Otherwise it returns false and an empty Entry.
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

	// Trim spaces
	base = utils.TrimSpace(base)

	// Trim prefix
	base = strings.TrimPrefix(base, "<string ")

	// Trim suffix
	base = strings.TrimSuffix(base, "</string>")

	// Get the name and value
	var baseSlice []string
	var formatted bool
	switch {
	case strings.Contains(base, ` formatted="false"`):
		baseSlice = strings.Split(base, `" formatted="false">`)
		formatted = true
	default:
		baseSlice = strings.Split(base, `">`)
	}
	name := strings.TrimPrefix(baseSlice[0], `name="`)
	value := baseSlice[1]

	if res.Format {

		// If value contains multiple _ and doesn't contain spaces we skip it
		if strings.Count(value, "_") >= 2 && strings.Count(value, " ") == 0 {
			return false, Entry{}
		}

		// If value contains multiple . and doesn't contain spaces we skip it
		if strings.Count(value, ".") > 2 && strings.Count(value, " ") == 0 {
			return false, Entry{}
		}
	}

	if res.ASCIIOnly {
		lenValue := len(value) - 1
		switch {
		case lenValue == -1, lenValue == 0:
			// We do nothing
		default:
			testOne, _ := utf8.DecodeRune([]byte{value[0]})
			testTwo, _ := utf8.DecodeRune([]byte{value[lenValue]})
			if testOne > 591 && testTwo > 591 {
				return false, Entry{}
			}
		}
	}

	// Determine if apostrophe's need to be fixed
	apostropheFix := strings.IndexByte(value, 39) >= 0

	// Determine if string needs to be formatted
	if strings.Count(value, "%s") >= 2 {
		formatted = true
	}
	if strings.Count(value, "%d") >= 2 {
		formatted = true
	}

	return true, Entry{
		name:          name,
		value:         value,
		apostropheFix: apostropheFix,
		formatted:     formatted,
	}
}
