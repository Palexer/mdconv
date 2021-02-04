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

func main() {
	// get output file
	outFileName := flag.String("o", "", "output file")
	flag.Parse()

	// get output file type
	var pdf bool
	if filepath.Ext(*outFileName) == ".pdf" {
		pdf = true
	} else if filepath.Ext(*outFileName) == ".html" || *outFileName == "" {
		pdf = false
	} else {
		log.Fatal("output file format not supported\nsupported formats: .pdf .html")
	}

	if *outFileName == "" {
		*outFileName = strings.TrimSuffix(filepath.Base(flag.Arg(0)), filepath.Ext(flag.Arg(0)))
		*outFileName = *outFileName + ".html"
	}

	// get input file
	file, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatal("failed to open file: ", err)
	}

	// parse markdown
	content := blackfriday.Run(file, blackfriday.WithRenderer(blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: blackfriday.CompletePage,
		CSS:   "style.css",
	})))

	// embedd CSS in html file
	output := append([]byte("<style>\n"), []byte(style)...)
	output = append(output, []byte("</style>\n")...)
	output = append(output, content...)

	if pdf {
		// create pdf output file
		pdfg, err := wkhtmltopdf.NewPDFGenerator()
		if err != nil {
			log.Fatal("failed to convert html to pdf: ", err)
		}

		pdfg.Dpi.Set(300)
		pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)

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
