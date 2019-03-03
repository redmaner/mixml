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

package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/redmaner/mixml/arraysxml"
	"github.com/redmaner/mixml/pluralsxml"
	"github.com/redmaner/mixml/stringsxml"
	"github.com/redmaner/mixml/utils"
)

const version = "r5"

var (
	cmdFmt     = flag.Bool("fmt", false, "Format strings.xml in a directory")
	cmdSort    = flag.Bool("sort", false, "Sort strings.xml in a directory")
	cmdVersion = flag.Bool("version", false, "Show version of mixml")

	parDir   = flag.String("dir", "./", "Path to directory")
	parQuiet = flag.Bool("quiet", false, "Suppress log messages")
	parASCII = flag.Bool("ascii", false, "Allow only ascii characters")
)

func main() {

	flag.Parse()
	switch {
	case *cmdFmt:
		format(*parDir, true)
	case *cmdSort:
		format(*parDir, false)
	case *cmdVersion:
		fmt.Printf("mixml version %s\nDeveloped by Redmaner\n", version)
	default:
		flag.PrintDefaults()
	}
}

// Format function
func format(dir string, format bool) {

	var apks []string
	filepath.Walk(dir, func(path string, f os.FileInfo, _ error) error {
		if filepath.Ext(path) == ".apk" {
			apks = append(apks, path)
		}
		return nil
	})

	var files []string
	for _, v := range apks {
		filepath.Walk(v, func(path string, f os.FileInfo, _ error) error {
			if !f.IsDir() {
				if f.Name() == "strings.xml" || f.Name() == "arrays.xml" || f.Name() == "plurals.xml" {
					files = append(files, path)
				}
			}
			return nil
		})
	}

	for _, v := range files {

		// Do basic integrity check of XML
		err := utils.XMLIntegrity(v)
		if err != nil {
			if !*parQuiet {
				fmt.Printf("Skipped: %s\n", err)
			}
			continue
		}

		switch path.Base(v) {

		case "strings.xml":
			if !*parQuiet {
				fmt.Printf("Formatting %s\n", v)
			}
			res := stringsxml.NewResources(v, format, *parASCII)
			res.Load()
			res.Write()

		case "arrays.xml":
			if !*parQuiet {
				fmt.Printf("Formatting %s\n", v)
			}
			res := arraysxml.NewResources(v, format, *parASCII)
			res.Load()
			res.Write()

		case "plurals.xml":
			if !*parQuiet {
				fmt.Printf("Formatting %s\n", v)
			}
			res := pluralsxml.NewResources(v, format, *parASCII)
			res.Load()
			res.Write()
		}
	}
}
