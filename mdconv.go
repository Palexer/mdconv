package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/russross/blackfriday/v2"
)

// bundle css for styling files
//go:generate go run script/bundlecss.go

func getCustomCSS(path string) []byte {
	if path == "" {
		return []byte{}
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("failed to open custom css file: ", err.Error())
	}
	return file
}

func main() {
	// get output file
	outFileName := flag.String("o", "", "output file (optional, default: HTML)")
	cssPath := flag.String("style", "", "pathToCSSFile (optional)")
	overwrite := flag.Bool("overwrite", false, "(optional, overwrites default CSS stylesheet)")
	flag.Parse()

	if filepath.Ext(flag.Arg(0)) != ".md" {
		log.Fatal("file type not supported: please use a .md input file")
	}

	// get output file type
	var pdf bool
	if filepath.Ext(*outFileName) == ".pdf" {
		pdf = true
	} else if filepath.Ext(*outFileName) == ".html" || *outFileName == "" {
		pdf = false
	} else {
		log.Fatal("output file format not supported\nsupported formats: .pdf .html")
	}

	// default: HTML output
	if *outFileName == "" {
		*outFileName = strings.TrimSuffix(filepath.Base(flag.Arg(0)), filepath.Ext(flag.Arg(0)))
		*outFileName = *outFileName + ".html"
	}

	// read input file
	file, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatal("failed to open file: ", err)
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
		output = append([]byte("<style>\n"), customCSS...)
		output = append(output, []byte("</style>\n")...)
		output = append(output, content...)
	} else {
		// no overwrite: include default and custom css
		output = append([]byte("<style>\n"), []byte(style)...)
		output = append(output, []byte("</style>\n")...)
		output = append(output, []byte("<style>\n")...)
		output = append(output, customCSS...)
		output = append(output, []byte("</style>\n")...)
		output = append(output, content...)
	}

	if pdf {
		// create pdf output file
		pdfg, err := wkhtmltopdf.NewPDFGenerator()
		if err != nil {
			log.Fatal("failed to convert html to pdf: ", err)
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
			log.Fatal("failed to create pdf in internal buffer: ", err)
		}

		// write pdf file to output path
		if err := pdfg.WriteFile(*outFileName); err != nil {
			log.Fatal("failed to write pdf file: ", err)
		}

	} else {
		// create html output file
		file, err := os.Create(*outFileName)
		if err != nil {
			log.Fatal("failed to create file: ", err)
		}

		if _, err := file.Write(output); err != nil {
			log.Fatal("failed to copy content from tmpfile: ", err)
		}
	}
}
