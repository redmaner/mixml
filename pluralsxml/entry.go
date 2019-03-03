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

package pluralsxml

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/redmaner/mixml/utils"
)

// Entry represents a plural in Android plurals.xml
type Entry struct {
	name  string
	items [][]string
}

// ParseEntry parses an Entry from a string. It returns true and the Entry if it was able
// to parse an Entry from the string. Otherwise it returns false and an empty item.
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

	var parsedFirstLine bool
	var strBuffer string
	var arrayName string
	var arrayItems [][]string

	sc := bufio.NewScanner(bytes.NewBufferString(base))

	for sc.Scan() {

		str := sc.Text()
		if !parsedFirstLine {
			parsedFirstLine = true

			// Trim prefix and suffix
			str = utils.TrimSpace(str)
			str = strings.TrimPrefix(str, `<plurals name="`)
			str = strings.TrimSuffix(str, `">`)
			arrayName = str
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
			str = utils.TrimSpace(str)
			str = strings.TrimPrefix(str, `<item quantity="`)
			str = strings.TrimSuffix(str, "</item>")

			strSlice := strings.Split(str, `">`)
			quantity := strSlice[0]
			value := strSlice[1]

			arrayItems = append(arrayItems, []string{quantity, value})
			continue
		}
		if strBuffer != "" {
			strBuffer = strBuffer + "\n"
		}
		strBuffer = strBuffer + str
	}

	return true, Entry{
		name:  arrayName,
		items: arrayItems,
	}
}
