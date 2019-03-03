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
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
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

		entry := res.Entries[key]
		buf := bytes.NewBufferString("")

		io.WriteString(buf, fmt.Sprintf(`    <plurals name="%s">`+"\n", entry.name))

		for _, item := range entry.items {
			io.WriteString(buf, fmt.Sprintf(`        <item quantity="%s">%s</item>`+"\n", item[0], item[1]))
		}
		io.WriteString(buf, fmt.Sprintf(`    </plurals>`+"\n"))

		io.WriteString(f, buf.String())

	}
	io.WriteString(f, fmt.Sprintf("</resources>\n"))
}
