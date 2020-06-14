package main

import (
	"flag"
	"os"

	"github.com/redmaner/mixml/src/miuires"
	"gopkg.in/yaml.v2"
)

const version = "r6"

// Commands
var cmdFormat bool
var cmdIntegrity bool

// Arguments
var argDir string
var argFilter bool
var argFilterConfig string
var argVerbose bool

func init() {
	flag.BoolVar(&cmdFormat, "format", false, "Format MIUI resources")
	flag.BoolVar(&cmdIntegrity, "check", false, "Check basic integrity of MIUI resources")
	flag.StringVar(&argDir, "dir", "./", "Directory of MIUI resources")
	flag.StringVar(&argDir, "d", "./", "Directory of MIUI resources")
	flag.BoolVar(&argFilter, "filter", false, "Filter MIUI resources")
	flag.BoolVar(&argFilter, "f", false, "Filter MIUI resources")
	flag.StringVar(&argFilterConfig, "config", "", "Path to filter configuration")
	flag.StringVar(&argFilterConfig, "c", "", "Path to filter configuration")
	flag.BoolVar(&argVerbose, "verbose", false, "Print verbose logging")
	flag.BoolVar(&argVerbose, "v", false, "Print verbose logging")
}

func main() {

	// Parse flags
	flag.Parse()

	switch {
	case cmdFormat:
		format()
	case cmdIntegrity:
		fc := miuires.FilterConfig{
			StringsKeyRules:   make(map[string][]miuires.FilterRules),
			StringsValueRules: make(map[string][]miuires.FilterRules),
		}

		fc.StringsKeyRules["all"] = append(fc.StringsKeyRules["all"], miuires.FilterRules{
			Match: "com.",
			Mode:  "prefix",
		})

		fc.StringsValueRules["all"] = append(fc.StringsValueRules["all"], miuires.FilterRules{
			Match: "@drawable",
			Mode:  "prefix",
		})

		fc.StringsValueRules["all"] = append(fc.StringsValueRules["all"], miuires.FilterRules{
			Match: "@string",
			Mode:  "prefix",
		})

		fc.StringsValueRules["all"] = append(fc.StringsValueRules["all"], miuires.FilterRules{
			Match: "@color",
			Mode:  "prefix",
		})

		fc.StringsValueRules["Settings.apk"] = append(fc.StringsValueRules["Settings.apk"], miuires.FilterRules{
			Match: "@drawable",
			Mode:  "prefix",
		})

		fc.StringsValueRules["Settings.apk"] = append(fc.StringsValueRules["Settings.apk"], miuires.FilterRules{
			Match: "@string",
			Mode:  "prefix",
		})

		fc.StringsValueRules["Settings.apk"] = append(fc.StringsValueRules["Settings.apk"], miuires.FilterRules{
			Match: "@color",
			Mode:  "prefix",
		})

		f, err := os.Create("./example-filter.yaml")
		if err != nil {
			panic(err)
		}

		out, _ := yaml.Marshal(&fc)
		f.Write(out)
		f.Close()

	default:
		flag.PrintDefaults()
		os.Exit(0)
	}
}
