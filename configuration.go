package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	pdf "github.com/adrg/go-wkhtmltopdf"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	_ "github.com/yuin/goldmark-emoji/definition"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

// page size options, height and width in mm
type pageSize struct {
	height int
	width  int
}

// switch width and height to get landscape options in HTML document
var (
	HTMLSizeLetterPortrait = pageSize{height: 279, width: 216}
	HTMLSizeLegalPortrait  = pageSize{height: 356, width: 216}
	HTMLSizeA1Portrait     = pageSize{height: 841, width: 594}
	HTMLSizeA2Portrait     = pageSize{height: 594, width: 420}
	HTMLSizeA3Portrait     = pageSize{height: 420, width: 297}
	HTMLSizeA4Portrait     = pageSize{height: 297, width: 210}
	HTMLSizeA5Portrait     = pageSize{height: 210, width: 148}
	HTMLSizeA6Portrait     = pageSize{height: 148, width: 105}
)

type configuration struct {
	inputfile           string
	outputfile          string
	customCSSPath       string
	overwriteDefaultCSS bool

	// output to PDF
	pdf bool

	// actual content
	HTMLContent []byte

	customCSS []byte

	// margins
	marginTop    uint
	marginRight  uint
	marginBottom uint
	marginLeft   uint

	dpi       uint64
	title     string
	grayscale bool

	// CSS strings
	margins     string
	orientation string
	size        string
}

func (c *configuration) getCustomCSS() {
	if c.customCSSPath == "" {
		return
	}

	var err error
	c.customCSS, err = os.ReadFile(c.customCSSPath)
	if err != nil {
		printErrExit("failed to open custom CSS file: ", err.Error())
	}
}

func (c *configuration) parseMDAndBundleStyles() {
	// read input file
	file, err := os.ReadFile(c.inputfile)
	if err != nil {
		printErrExit("error: failed to open file: ")
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, emoji.Emoji),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	var content bytes.Buffer
	md.Convert(file, &content)

	// bundle custom CSS provided by the user
	c.getCustomCSS()

	// embed CSS in file
	if c.overwriteDefaultCSS {
		// overwrite: only include custom CSS
		css := []byte(fmt.Sprintf("<style>\n%s\n</style>", string(c.customCSS)))
		c.HTMLContent = append(css, content.Bytes()...)
	} else {
		// no overwrite: include default and custom CSS, don't include margins if it outputs to PDF
		var css []byte
		if c.pdf {
			css = []byte(fmt.Sprintf("<style>\n%s\n</style>\n\n<style>%s</style>\n", GHStyle, c.customCSS))
		} else {
			css = []byte(fmt.Sprintf("<style>\n%s\n%s\n</style>\n\n<style>%s</style>\n", c.margins, GHStyle, c.customCSS))
		}
		c.HTMLContent = append(css, content.Bytes()...)
	}

}

func (c *configuration) createHTMLFile() {
	// create html output file

	// page size
	var width, height int

	switch strings.ToLower(c.size) {
	case "letter":
		width = HTMLSizeLetterPortrait.width
		height = HTMLSizeLetterPortrait.height
	case "legal":
		width = HTMLSizeLegalPortrait.width
		height = HTMLSizeLegalPortrait.height
	case "a1":
		width = HTMLSizeA1Portrait.width
		height = HTMLSizeA1Portrait.height
	case "a2":
		width = HTMLSizeA2Portrait.width
		height = HTMLSizeA2Portrait.height
	case "a3":
		width = HTMLSizeA3Portrait.width
		height = HTMLSizeA3Portrait.height
	case "a4":
		width = HTMLSizeA4Portrait.width
		height = HTMLSizeA4Portrait.height
	case "a5":
		width = HTMLSizeA5Portrait.width
		height = HTMLSizeA5Portrait.height
	case "a6":
		width = HTMLSizeA6Portrait.width
		height = HTMLSizeA6Portrait.height
	}

	// switch width and height if the page is in landscape mode
	if c.orientation == "landscape" {
		tmp := width
		width = height
		height = tmp
	}

	sizeStr := fmt.Sprintf("<style>\n\thtml { width: %dmm; height: %dmm; }\n</style>", width, height)
	c.HTMLContent = append([]byte(sizeStr), c.HTMLContent...)

	file, err := os.Create(c.outputfile)
	if err != nil {
		printErrExit("error: failed to create HTML file: ", err)
	}

	if _, err := file.Write(c.HTMLContent); err != nil {
		printErrExit("error: failed to write to HTML file: ", err)
	}

}

func (c *configuration) createPDFFile() {
	if err := pdf.Init(); err != nil {
		printErrExit(err)
	}
	defer pdf.Destroy()

	bytes.NewReader(c.HTMLContent)

	obj, err := pdf.NewObjectFromReader(bytes.NewReader(c.HTMLContent))
	if err != nil {
		printErrExit(err)
	}

	converter, err := pdf.NewConverter()
	if err != nil {
		printErrExit(err)
	}
	defer converter.Destroy()

	converter.Add(obj)

	converter.MarginTop = fmt.Sprintf("%dmm", c.marginLeft)
	converter.MarginRight = fmt.Sprintf("%dmm", c.marginRight)
	converter.MarginBottom = fmt.Sprintf("%dmm", c.marginBottom)
	converter.MarginLeft = fmt.Sprintf("%dmm", c.marginLeft)

	converter.DPI = c.dpi
	converter.Title = c.title

	if c.grayscale {
		converter.Colorspace = pdf.Grayscale
	}

	// page size
	switch strings.ToLower(c.size) {
	case "letter":
		converter.PaperSize = pdf.Letter
	case "legal":
		converter.PaperSize = pdf.Legal
	case "a1":
		converter.PaperSize = pdf.A1
	case "a2":
		converter.PaperSize = pdf.A2
	case "a3":
		converter.PaperSize = pdf.A3
	case "a4":
		converter.PaperSize = pdf.A4
	case "a5":
		converter.PaperSize = pdf.A5
	case "a6":
		converter.PaperSize = pdf.A6
	}

	switch strings.ToLower(c.orientation) {
	case "portrait":
		converter.Orientation = pdf.Portrait
	case "landscape":
		converter.Orientation = pdf.Landscape
	}

	file, err := os.Create(c.outputfile)
	if err != nil {
		printErrExit("failed to create output file: ", err)
	}
	defer file.Close()

	if err := converter.Run(file); err != nil {
		printErrExit("failed to write output file: ", err)
	}
}

func createConfiguration() *configuration {
	// define flags
	outFileName := flag.String("o", "", "Output file, file extension is used to determine the output file type (default HTML)")
	cssPath := flag.String("css", "", "Specify path to custom CSS file")
	overwrite := flag.Bool("overwrite", false, "Overwrite default CSS")

	orientation := flag.String("orientation", "portrait", "PDF orientation (portrait (default) / landscape)")
	size := flag.String("size", "A4", "The size of a PDF page (A4 (default), A5)")
	dpi := flag.Uint64("dpi", 300, "Specify the DPI of the PDF file (e.g. 96; default: 300)")
	grayscale := flag.Bool("grayscale", false, "Choose whether the PDF file should be grayscale only (default: false)")

	marginLeft := flag.Int("margin-left", 20, "Specify a left margin in mm")
	marginRight := flag.Int("margin-right", 20, "Specify a right margin in mm")
	marginTop := flag.Int("margin-top", 20, "Specify a top margin in mm")
	marginBottom := flag.Int("margin-bottom", 20, "Specify a bottom margin in mm")

	versionFlag := flag.Bool("v", false, "Show currently used mdconv version")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("mdconv version %s\n", version)
		os.Exit(0)
	}

	// get the input file
	input := flag.Arg(0)
	if filepath.Ext(input) != ".md" {
		printErrExit("error: input file type not supported or file not found (please use a *.md file)", "\nSee 'mdconv -h' or 'man mdconv' for more information")
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

	return &configuration{
		inputfile:           input,
		outputfile:          *outFileName,
		customCSSPath:       *cssPath,
		overwriteDefaultCSS: *overwrite,

		pdf: pdf,

		orientation: *orientation,
		size:        *size,
		dpi:         *dpi,
		grayscale:   *grayscale,

		marginTop:    uint(*marginTop),
		marginRight:  uint(*marginRight),
		marginBottom: uint(*marginBottom),
		marginLeft:   uint(*marginLeft),
		margins: fmt.Sprintf("@page { margin-top: %dmm; margin-right: %dmm; margin-bottom: %dmm; margin-left: %dmm; }\n @media screen { html {  margin-top: %dmm; margin-right: %dmm; margin-bottom: %dmm; margin-left: %dmm; } }",
			*marginTop, *marginRight, *marginBottom, *marginLeft,
			*marginTop, *marginRight, *marginBottom, *marginLeft),
	}
}
