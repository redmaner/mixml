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
	"bufio"
	"log"
	"os"
	"sort"
	"strings"
)

// Load loads android strings.xml
func (res *Resources) Load() {

	// load file
	f, err := os.Open(res.FilePath)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	defer f.Close()

	// We scan the file with a bufio.Scanner. Each string item is stored in a slice.
	sc := bufio.NewScanner(f)
	var stringEntries []string
	var strPlaceholder string

	for sc.Scan() {

		// We scan each line of the file.
		str := sc.Text()

		// If a string item has trash in it, we skip the line
		if strings.Contains(str, "resources>") || strings.Contains(str, "<?xml") || strings.Contains(str, "/>") {
			continue
		}

		// We want to join strings together <string> </string>
		if strings.Contains(str, "<string") || strings.Contains(str, "<!--") {
			stringEntries = append(stringEntries, strPlaceholder)
			strPlaceholder = ""
			strPlaceholder = strPlaceholder + str
			continue
		}
		strPlaceholder = strPlaceholder + "\n" + str
	}
	stringEntries = append(stringEntries, strPlaceholder)

	// We put every string in a map. This makes sure we have unique keys.
	// This way we remove double string items.
	for _, v := range stringEntries {
		if ok, val := res.ParseEntry(v); ok {
			res.Entries[val.name] = val
		}
	}

	// We store xmlKeys in a separte slice and sort it, this way we can rebuild
	// the file in a ordered way.
	for k := range res.Entries {
		res.Keys = append(res.Keys, k)
	}
	sort.Strings(res.Keys)
}
