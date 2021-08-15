package main

import (
	_ "embed"
	"fmt"
	"os"
)

// bundle css for styling files
//go:embed styles/gh.css
var style string

const version = "0.95"

func printErrExit(a ...interface{}) {
	fmt.Fprintf(os.Stderr, "mdconv: ")
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func getCustomCSS(path string) []byte {
	if path == "" {
		return []byte{}
	}
	file, err := os.ReadFile(path)
	if err != nil {
		printErrExit("failed to open custom css file: ", err.Error())
	}
	return file
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
