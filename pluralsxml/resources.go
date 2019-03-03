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

// Resources represents a set of strings in an android strings.xml
type Resources struct {
	FilePath  string
	Keys      []string
	Entries   map[string]Entry
	Comment   string
	Format    bool
	ASCIIOnly bool
}

// NewResources returns an empty resources struct
func NewResources(filePath string, format bool, asciiOnly bool) *Resources {
	return &Resources{
		FilePath:  filePath,
		Keys:      []string{},
		Entries:   make(map[string]Entry),
		Format:    format,
		ASCIIOnly: asciiOnly,
	}
}
