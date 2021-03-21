# mdconv

## About

mdconv is a markdown converter written in Go.
It is able to create PDF and HTML files from Markdown without using LaTeX. 
Instead, mdconv uses the Blackfriday (v2) Markdown processor and go-wkhtmltopdf to convert the HTML
to PDF.

## Installation

### Download

1. Install [wkhtmltopdf](https://wkhtmltopdf.org/downloads.html) for PDF conversions.
2. Download mdconv from the [releases section](https://github.com/Palexer/mdconv/releases).

### Compile from source

1. Clone this repository to your local machine
2. Install the dependencies: go, wkhtmltopdf
3. Run ```make``` followed by ```sudo make install```
_Note: Run can also ```sudo make uninstall``` to remove the program_

## Usage

### General Usage

**Convert a Markdown document to HTML:**


```mdconv path/to/markdowndocument.md```


**Convert a Markdown document to PDF:**


```mdconv -o output.pdf path/to/markdowndocument.md```

_Note: The output file type is defined by the file extension of the output file
specified with ```-o```._

### Flags

|Flag|Description|
|----|------|
|-o out.ext|Specify the output file name and file type|
|-c [FILE]|Specify the path to a custom CSS style sheet|
|-overwrite|Don't include the default style sheet|
|-f [FONT]|Specify a font family (sans, serif, monospace)|
|-h, -help|Get general usage information|
|-v, -V|print currently used mdconv version|

_Note: The path to the input file must be provided _after_ the flags are specified._

## Features

- [x] HTML and PDF output 
- [x] Standard Markdown features like headings, images, lists, code blocks, embedded HTML, tables, etc.
- [x] Custom CSS stylesheet for output files
- [x] Github-like default stylesheet

## Contributing

Contributions of all kinds are very welcome. See the GitHub Issue Tracker for things you might want work on.

## License

[MIT License](https://raw.githubusercontent.com/Palexer/mdconv/master/LICENSE)
