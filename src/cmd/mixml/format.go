package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/redmaner/mixml/src/miuires"
)

// Format function
func format() {

	var apks []string
	filepath.Walk(argDir, func(path string, f os.FileInfo, _ error) error {
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

	// Apply filter if defined
	var fc *miuires.FilterConfig
	var filter bool
	if argFilter && argFilterConfig != "" {
		if f, err := os.Open(argFilterConfig); err == nil {
			defer f.Close()
			fc, err = miuires.GetFilterConfigFromFile(f)
			if err != nil {
				fmt.Printf("Yaml unmarshal error: %v\n", err)
			}
			filter = true
		} else {
			fmt.Printf("Couldn't open filter configuration: %v\n", err)
		}
	}

	for _, v := range files {
		res := miuires.NewResources(v)
		if err := res.Load(); err != nil {
			fmt.Printf("An error occurred when loading %s: %v\n", v, err)
			continue
		}

		if filter {
			res.Filter(fc)
		}

		res.Write()
		if argVerbose {
			fmt.Printf("Formatted %s\n", v)
		}
	}
}
