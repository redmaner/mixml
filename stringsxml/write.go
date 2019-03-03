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
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Write writes android strings.xml
func (res *Resources) Write() {

	if len(res.Keys) == 0 {
		return
	}

	if _, err := os.Stat(res.FilePath); err == nil {
		err := os.Remove(res.FilePath)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
	}

	f, err := os.Create(res.FilePath)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	defer f.Close()

	io.WriteString(f, fmt.Sprintf("<?xml version='1.0' encoding='UTF-8'?>\n"))

	if res.Comment != "" {
		io.WriteString(f, res.Comment)
	}

	io.WriteString(f, fmt.Sprintf("<resources>\n"))

	for _, key := range res.Keys {

		var formatString string
		var writeString string
		stringItem := res.Entries[key]

		if stringItem.formatted {
			formatString = ` formatted="false"`
		}

		value := stringItem.value

		// We fix apostrophe errors, by adding a \ in front of it
		// this slash is only added when it does not yet exist
		if stringItem.apostropheFix && value[0] != '"' {
			var newValue string

			strSlice := strings.Split(value, "'")
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
			value = newValue
		}

		writeString = fmt.Sprintf(`  <string name="%s"%s>%s</string>`+"\n", stringItem.name, formatString, value)

		f.WriteString(writeString)

	}
	io.WriteString(f, fmt.Sprintf("</resources>\n"))

}
