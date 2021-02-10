# mdconv

## About

MDConv is a markdown converter written in Go.
It is able to create PDF and HTML files from Markdown without using LaTeX. 
Instead MDConv uses the Blackfriday (v2) markdown processor and go-wkhtmltopdf to convert the HTML
to PDF.

## Installation

1. Install [wkhtmltopdf](https://wkhtmltopdf.org/downloads.html)
2. Download mdconv from the [releases section](https://github.com/Palexer/mdconv/releases)

## Usage

### General usage

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
|-style style.css|You can specify an additional stylesheet for your output file, which will be linked to in the HTML head.|
|-overwrite|If the -overwrite flag is parsed, the default stylesheet is not included in the output file.|
|-help / -h|Get general usage inforamtion.|

_Note: The path to the input file must be provided _after_ the flags are specified._

## Contributing

Contributions of all kinds are very welcome. See the Github Issue Tracker or the ToDo below
for things you might want to do.

## ToDo

- add more testing
- fix PDF output, if -overwrite flag is invoked (problem: overwrite still keeps the default stylesheet and adds the custom one
to it (like -style without -overwrite))

## License

MIT License
