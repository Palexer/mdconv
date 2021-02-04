package main

import (
	"bytes"
	"io"
	"log"
	"os"
)

// bundle the style.css
func main() {
	style, err := os.Open("style.css")
	if err != nil {
		log.Fatal("failed to open css file: ", err)
	}

	file, err := os.Create("bundled.go")
	_, err = file.Write([]byte("package main \n\nconst style = `"))
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(file, style)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write([]byte("`"))
	if err != nil {
		log.Fatal(err)
	}

}
