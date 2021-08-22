# Mdconv

## About

Mdconv is a Markdown converter written in Go.
The goal of this project is to create a Markdown to PDF converter, that doesn't need big dependencies like
headless Chromium or LaTeX and is written in a compiled language, so that NodeJS is not necessary.
It is able to create PDF and HTML files from Markdown without using LaTeX or any other big dependency. 
Instead, mdconv uses the goldmark Markdown processor and relies on wkhtmltopdf (go-wkhtmltopdf) to convert the HTML
to PDF.

## Installation

### Download

1. Install [wkhtmltopdf](https://wkhtmltopdf.org/downloads.html) for PDF conversions.
2. Download mdconv from the [releases section](https://github.com/Palexer/mdconv/releases).

### Compile from source

1. Clone this repository to your local machine
2. Install the dependencies: go, wkhtmltopdf
3. Run ```sudo make install```

_Note: Run can also ```sudo make uninstall``` to remove the program_

## Usage

### General Usage

**Convert a Markdown document to HTML:**


```mdconv path/to/markdowndocument.md```


**Convert a Markdown document to PDF:**


```mdconv -o output.pdf path/to/markdowndocument.md```

_Note: The output file type is defined by the file extension of the output file
specified with ```-o```. Consequently you'll always have to use the ```-o``` flag to output to PDF_

For all available options see ```mdconv -h```

## Features

- [x] HTML and PDF output 
- [x] standard Markdown features like headings, images, lists, code blocks, embedded HTML, tables, etc.
- [x] custom CSS stylesheet for output files
- [x] configurable via command line flags

## ToDo

- write more tests and test pagesize and orientation
- PDF emoji support
- more options as flags with new converter lib
- correct display of checkboxes

## Contributing

Contributions of all kinds are very welcome.

## License

[MIT License](https://raw.githubusercontent.com/Palexer/mdconv/master/LICENSE)
