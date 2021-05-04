package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/russross/blackfriday/v2"
)

// bundle css for styling files
//go:embed styles/gh.css
var style string

//go:embed styles/fonts/monospace.css
var fontMonospace string

//go:embed styles/fonts/serif.css
var fontSerif string

//go:embed styles/fonts/sansserif.css
var fontSansSerif string

const version = "0.9"

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
	// define flags
	outFileName := flag.String("o", "", "output file, file extension is used to determine the output file type (default HTML)")
	cssPath := flag.String("c", "", "path to custom CSS file")
	font := flag.String("f", "sans", "specify the font for the output document (sans, serif, monospace)")
	overwrite := flag.Bool("overwrite", false, "overwrites default CSS stylesheet")
	versionShort := flag.Bool("V", false, "show currently used mdconv version")
	versionLong := flag.Bool("version", false, "show currently used mdconv version")
	flag.Parse()

	if *versionShort || *versionLong {
		fmt.Printf("mdconv version %s\n", version)
		os.Exit(0)
	}

	// get the input file
	input := flag.Arg(0)
	if filepath.Ext(input) != ".md" {
		printErrExit("error (wrong input file): file type not supported (please use a .md input file) or file not found", "\nSee mdconv -h or man mdconv for more information")
	}

	// get output file type
	var pdf bool
	if filepath.Ext(*outFileName) == ".pdf" {
		pdf = true
	} else if filepath.Ext(*outFileName) == ".html" || *outFileName == "" {
		pdf = false
	} else {
		printErrExit("error (wrong output file format): format not supported\nsupported formats: .pdf .html")
	}

	// default: HTML output
	if *outFileName == "" {
		*outFileName = strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))
		*outFileName = *outFileName + ".html"
	}

	// include correct font
	var fontStyle string
	switch strings.ToLower(*font) {
	case "sans":
		fontStyle = fontSansSerif
	case "serif":
		fontStyle = fontSerif
	case "monospace":
		fontStyle = fontMonospace
	default:
		printErrExit("font family not supported, see mdconv --help for more information")
	}

	// read input file
	file, err := os.ReadFile(input)
	if err != nil {
		printErrExit("error: failed to open file: ")
	}

	// parse markdown
	content := blackfriday.Run(file, blackfriday.WithRenderer(blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: blackfriday.CompletePage,
	})))

	// bundle custom CSS provided by the user
	customCSS := getCustomCSS(*cssPath)

	// embedd CSS in html file
	var output []byte
	if *overwrite {
		// overwrite: only include custom css
		css := []byte(fmt.Sprintf("<style>\n%s\n</style>", string(customCSS)))
		output = append(css, content...)
	} else {
		// no overwrite: include default and custom css
		css := []byte(fmt.Sprintf("<style>\n%s\n%s\n</style>\n\n<style>%s</style>\n", fontStyle, style, customCSS))
		output = append(css, content...)
	}

	if pdf {
		// create pdf output file
		pdfg, err := wkhtmltopdf.NewPDFGenerator()
		if err != nil {
			printErrExit("error: failed to convert html to pdf: ", err)
		}

		// define page options
		pdfg.Dpi.Set(300)
		pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
		pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
		pdfg.MarginTop.Set(20)
		pdfg.MarginBottom.Set(20)
		pdfg.MarginLeft.Set(20)
		pdfg.MarginRight.Set(20)

		pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(output)))

		if err := pdfg.Create(); err != nil {
			printErrExit("error: failed to create pdf in internal buffer: ", err)
		}

		// write pdf file to output path
		if err := pdfg.WriteFile(*outFileName); err != nil {
			printErrExit("error: failed to write pdf file: ", err)
		}

	} else {
		// create html output file
		file, err := os.Create(*outFileName)
		if err != nil {
			printErrExit("error: failed to create HTML file: ", err)
		}

		if _, err := file.Write(output); err != nil {
			printErrExit("error: failed to write to HTML file: ", err)
		}
	}
}
