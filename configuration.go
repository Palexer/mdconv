package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/russross/blackfriday/v2"
)

const (
	HTMLOrientationPortrait  = ""
	HTMLOrientationLandscape = ""
	HTMLSizeLetter           = ""
	HTMLSizeLegal            = ""
	HTMLSizeA1               = ""
	HTMLSizeA2               = ""
	HTMLSizeA3               = ""
	HTMLSizeA4               = ""
	HTMLSizeA5               = ""
	HTMLSizeA6               = ""

	FontSerif     = "html { font-family: Times New Roman, Times, serif, Georgia, Gramond; }"
	FontSansSerif = "html { font-family: helvetica, arial, freesans, clean, sans-serif, Liberation Sans, Calibri; }"
	FontMonospace = "html {	font-family: monospace, Courier New, Lucida Console, Monaco; } "
)

type configuration struct {
	inputfile           string
	outputfile          string
	csspath             string
	overwriteDefaultCSS bool

	pdf bool

	HTMLContent []byte

	fonttype    string
	orientation string
	size        string

	marginTop    int
	marginRight  int
	marginBottom int
	marginLeft   int
}

func (c *configuration) parseMDAndBundleStyles() {
	// read input file
	file, err := os.ReadFile(c.inputfile)
	if err != nil {
		printErrExit("error: failed to open file: ")
	}

	// parse markdown
	content := blackfriday.Run(file, blackfriday.WithRenderer(blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: blackfriday.CompletePage,
	})))

	// bundle custom CSS provided by the user
	var customCSS []byte
	if c.csspath != "" {
		customCSS = getCustomCSS(c.csspath)
	}

	// embed CSS in file
	if c.overwriteDefaultCSS {
		// overwrite: only include custom CSS
		css := []byte(fmt.Sprintf("<style>\n%s\n</style>", string(customCSS)))
		c.HTMLContent = append(css, content...)
	} else {
		// no overwrite: include default and custom CSS
		css := []byte(fmt.Sprintf("<style>\n%s\n%s\n</style>\n\n<style>%s</style>\n", c.fonttype, style, customCSS))
		c.HTMLContent = append(css, content...)
	}

}

func (c *configuration) createHTMLFile() {
	// create html output file
	file, err := os.Create(c.outputfile)
	if err != nil {
		printErrExit("error: failed to create HTML file: ", err)
	}

	if _, err := file.Write(c.HTMLContent); err != nil {
		printErrExit("error: failed to write to HTML file: ", err)
	}

}

func (c *configuration) createPDFFile() {
	// create pdf output file
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		printErrExit("error: failed to convert html to pdf: ", err)
	}

	// define page options
	// DPI
	pdfg.Dpi.Set(300)

	// page orientation
	switch strings.ToLower(c.orientation) {
	case "portrait":
		pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	case "landscape":
		pdfg.Orientation.Set(wkhtmltopdf.OrientationLandscape)

	}

	// page size
	switch strings.ToLower(c.size) {
	case "letter":
		pdfg.PageSize.Set(wkhtmltopdf.PageSizeLetter)
	case "legal":
		pdfg.PageSize.Set(wkhtmltopdf.PageSizeLegal)
	case "a1":
		pdfg.PageSize.Set(wkhtmltopdf.PageSizeA1)
	case "a2":
		pdfg.PageSize.Set(wkhtmltopdf.PageSizeA2)
	case "a3":
		pdfg.PageSize.Set(wkhtmltopdf.PageSizeA3)
	case "a4":
		pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	case "a5":
		pdfg.PageSize.Set(wkhtmltopdf.PageSizeA5)
	case "a6":
		pdfg.PageSize.Set(wkhtmltopdf.PageSizeA6)
	}

	// set default page margins
	pdfg.MarginTop.Set(20)
	pdfg.MarginBottom.Set(20)
	pdfg.MarginLeft.Set(5)
	pdfg.MarginRight.Set(20)

	// set parsed page margins
	if c.marginLeft > -1 {
		pdfg.MarginLeft.Set(uint(c.marginLeft))
	}

	if c.marginRight > -1 {
		pdfg.MarginRight.Set(uint(c.marginRight))
	}

	if c.marginTop > -1 {
		pdfg.MarginTop.Set(uint(c.marginTop))
	}

	if c.marginBottom > -1 {
		pdfg.MarginBottom.Set(uint(c.marginBottom))
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(c.HTMLContent)))

	if err := pdfg.Create(); err != nil {
		printErrExit("error: failed to create pdf in internal buffer: ", err)
	}

	// write pdf file to output path
	if err := pdfg.WriteFile(c.outputfile); err != nil {
		printErrExit("error: failed to write pdf file: ", err)
	}

}

func createConfiguration() *configuration {
	// define flags
	outFileName := flag.String("o", "", "Output file, file extension is used to determine the output file type (default HTML)")
	cssPath := flag.String("c", "", "Specify path to custom CSS file")
	overwrite := flag.Bool("overwrite", false, "Overwrite default CSS")
	font := flag.String("f", "sans", "Specify the font for the output document (sans, serif, monospace)")

	orientation := flag.String("orientation", "portrait", "PDF orientation (portrait (default) / landscape)")
	pagesize := flag.String("pagesize", "A4", "The size of a PDF page (A4 (default), A5)")

	marginLeft := flag.Int("margin-left", -1, "Specify a left margin in mm")
	marginRight := flag.Int("margin-right", -1, "Specify a right margin in mm")
	marginTop := flag.Int("margin-top", -1, "Specify a top margin in mm")
	marginBottom := flag.Int("margin-bottom", -1, "Specify a bottom margin in mm")

	versionShort := flag.Bool("V", false, "Show currently used mdconv version")
	versionLong := flag.Bool("version", false, "Show currently used mdconv version")

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

	// HTML file with input name if no output file was specified
	if *outFileName == "" {
		*outFileName = strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))
		*outFileName = *outFileName + ".html"
	}

	// include correct font
	var fontStyle string
	switch strings.ToLower(*font) {
	case "sans":
		fontStyle = FontSansSerif
	case "serif":
		fontStyle = FontSerif
	case "monospace":
		fontStyle = FontMonospace
	default:
		printErrExit("font family not supported, see mdconv --help for more information")
	}

	return &configuration{
		inputfile:           input,
		outputfile:          *outFileName,
		csspath:             *cssPath,
		overwriteDefaultCSS: *overwrite,

		pdf: pdf,

		fonttype:    fontStyle,
		orientation: *orientation,
		size:        *pagesize,

		marginTop:    *marginTop,
		marginRight:  *marginRight,
		marginBottom: *marginBottom,
		marginLeft:   *marginLeft,
	}
}
