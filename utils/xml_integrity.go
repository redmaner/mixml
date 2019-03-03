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

package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// XMLIntegrity checks basic XML integrity of strings.xml and arrays.xml
func XMLIntegrity(filePath string) error {

	// Open file, defer closure
	f, err := os.Open(filePath)
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
		return fmt.Errorf("%s: basic XML integrity check failed (resources)", filePath)
	}

	switch path.Base(filePath) {

	// strings.xml checks
	case "strings.xml":

		// Count the occurrence of <string></string> pairs. These should be equal.
		stringOpenCount := bytes.Count(data, []byte("<string name="))
		stringCloseCount := bytes.Count(data, []byte("</string>"))
		stringCloseCount = stringCloseCount + bytes.Count(data, []byte(`"/>`))
		stringCloseCount = stringCloseCount + bytes.Count(data, []byte(`" />`))

		if stringOpenCount != stringCloseCount {
			return fmt.Errorf("%s: basic XML integrity check failed (string mismatch)", filePath)
		}

		// arrays.xml checks
	case "arrays.xml":

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
			return fmt.Errorf("%s: basic XML integrity check failed (array mismatch)", filePath)
		}

		// Count the occurrence of <item></item> pairs. These should be equal.
		itemOpenCount := bytes.Count(data, []byte("<item>"))
		itemCloseCount := bytes.Count(data, []byte("</item>"))

		if itemOpenCount != itemCloseCount {
			return fmt.Errorf("%s: basic XML integrity check failed (item mismatch)", filePath)
		}

	// Checks for plurals.xml
	case "plurals.xml":

		// Count the occurrence of <plurals></plurals> pairs. These should be equal.
		pluralsOpenCount := bytes.Count(data, []byte("<plurals name="))
		pluralsCloseCount := bytes.Count(data, []byte("</plurals>"))

		if pluralsOpenCount != pluralsCloseCount {
			return fmt.Errorf("%s: basic XML integrity check failed (plurals mismatch)", filePath)
		}

		// Count the occurrence of <item></item> pairs. These should be equal.
		itemOpenCount := bytes.Count(data, []byte("<item"))
		itemCloseCount := bytes.Count(data, []byte("</item>"))

		if itemOpenCount != itemCloseCount {
			return fmt.Errorf("%s: basic XML integrity check failed (item mismatch)", filePath)
		}
	}

	return nil

}
