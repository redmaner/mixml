package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "r6"

// Commands
var cmdFormat = flag.NewFlagSet("format", flag.ExitOnError)
var cmdCheck = flag.NewFlagSet("check", flag.ExitOnError)

// Arguments
var argDir string
var argFilter bool
var argFilterConfig string
var argVerbose bool
var argHelp bool

func init() {

	// Arguments for format
	cmdFormat.StringVar(&argDir, "dir", "./", "Directory of MIUI resources")
	cmdFormat.StringVar(&argDir, "d", "./", "Directory of MIUI resources")
	cmdFormat.BoolVar(&argFilter, "filter", false, "Filter MIUI resources")
	cmdFormat.BoolVar(&argFilter, "f", false, "Filter MIUI resources")
	cmdFormat.StringVar(&argFilterConfig, "config", "", "Path to filter configuration")
	cmdFormat.StringVar(&argFilterConfig, "c", "", "Path to filter configuration")
	cmdFormat.BoolVar(&argHelp, "help", false, "Show help")
	cmdFormat.BoolVar(&argHelp, "h", false, "Show help")
	cmdFormat.BoolVar(&argVerbose, "verbose", false, "Print verbose logging")
	cmdFormat.BoolVar(&argVerbose, "v", false, "Print verbose logging")

	// Arguments for check
	cmdCheck.StringVar(&argDir, "dir", "./", "Directory of MIUI resources")
	cmdCheck.StringVar(&argDir, "d", "./", "Directory of MIUI resources")
}

func main() {

	args := os.Args
	if len(args) < 2 {
		showHelp()
	}

	switch args[1] {
	case "format":
		if err := cmdFormat.Parse(args[2:]); err != nil {
			fmt.Println(err)
			showHelp()
		}
		format()
	case "check":
	default:
		showHelp()
	}
}
