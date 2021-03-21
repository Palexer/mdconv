# mdconv

## About

MDConv is a markdown converter written in Go.
It is able to create PDF and HTML files from Markdown without using LaTeX. 
Instead MDConv uses the Blackfriday (v2) markdown processor and go-wkhtmltopdf to convert the HTML
to PDF.

## Installation

1. Install [wkhtmltopdf](https://wkhtmltopdf.org/downloads.html) for PDF conversions.
2. Download mdconv from the [releases section](https://github.com/Palexer/mdconv/releases).

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
|-o out.pdf|You can specify an output file with either the .html or .pdf extension. If -o is not provided it defaults to the markdown file name and the .html file extension|
|-c style.css|You can specify an additional stylesheet for your output file, which will be linked to in the HTML head.|
|-overwrite|If the -overwrite flag is parsed, the default stylesheet is not included in the output file.|
|-f sans|Specify a custom font family (sans, serif, monospace), default: sans|
|-help / -h|Get general usage inforamtion.|

_Note: The path to the input file must be provided _after_ the flags are specified._

## Features

- [x] HTML and PDF output 
- [x] Standard Markdown features like headings, images, lists, code blocks, embedded HTML, tables, etc.
- [x] Custom CSS stylesheet for output files
- [x] Github-like default stylesheet

## Contributing

Contributions of all kinds are very welcome. See the Github Issue Tracker for things you might want work on.

## ToDo

- add TOC support
- fix PDF background/border

## License

[MIT License](https://raw.githubusercontent.com/Palexer/mdconv/master/LICENSE)
