package main

import (
	_ "embed"
	"fmt"
	"os"
)

// bundle css for styling files
//go:embed styles/gh.css
var GHStyle string

const version = "0.95"

func printErrExit(a ...interface{}) {
	fmt.Fprintf(os.Stderr, "mdconv: ")
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func main() {
	config := createConfiguration()
	config.parseMDAndBundleStyles()

	if config.pdf {
		config.createPDFFile()
	} else {
		config.createHTMLFile()
	}
}
