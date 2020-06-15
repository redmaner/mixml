package main

import (
	"fmt"
	"os"
)

const helpMessage = `
mixml version: %s (by redmaner)

Usage:
    mixml <command> <options>

Commands:
    format             Format MIUI resources
    help               Show this help

`

const helpMessageFormat = `
mixml version: %s (by redmaner)

Usage:
    mixml format <options>

Options:
    --dir     | -d      Path of directory to format
    --filter  | -f      Enable filter when formatting
    --config  | -c      Path to the filter configuration YAML file
    --verbose | -v      Show verbose logging
    --help    | -h      Show this help

`

func showHelp() {
	fmt.Printf(helpMessage, version)
	os.Exit(10)
}

func showHelpFormat() {
	fmt.Printf(helpMessageFormat, version)
	os.Exit(10)
}
