package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/redmaner/mixml/stringsxml"
)

const version = "r1"

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
				if f.Name() == "strings.xml" {
					files = append(files, path)
				}
			}
			return nil
		})
	}

	for _, v := range files {
		res := stringsxml.NewResources(v, format, *parASCII)
		res.Load()
		res.Write()
		if !*parQuiet {
			fmt.Printf("Formatted %s\n", v)
		}
	}
}
