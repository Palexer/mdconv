package main 

const style = `/* general */

/* Note: Since webkit doesn't support @page, this will be ignored 
   when converting to pdf using wkhtmltopdf. The page margins for pdfs are defined in mdconv.go 
   @page defines the margins when printing through a web browser */
@page {
    margin-top: 2cm;
    margin-bottom: 2cm;
    margin-left: 2cm;
    margin-right: 2cm;
}

/* page margins when viewing on a screen */
@media screen {
	html {
	    margin-top: 2cm;
        margin-bottom: 2cm;
        margin-left: 2cm;
        margin-right: 2cm;
	}
}

html {
    font-size: 12pt;
    line-height: 1.5;
    font-family: helvetica, arial, freesans, clean, sans-serif;
    color: black;
    word-wrap: break-word;
    text-align: justify;
	height: 297mm;
    width: 210mm;
}

body {
    margin: auto;
}

/* headings */
h1 {
    margin: 15px 0;
    padding-bottom: 2px;
    font-size: 24pt;
    border-bottom: 1px solid #EEE;
}

h2 {
    margin: 20px 0 10px 0;
    font-size: 22pt;
}

h3 {
    margin: 20px 0 10px 0;
    padding-bottom: 2px;
    font-size: 18pt;
}

h4 {
    font-size: 16pt;
    line-height: 26px;
    padding: 14px 0 4px;
    font-weight: bold;
}

h5 {
    font-size: 14pt;
    line-height: 26px;
    padding: 14px 0 0;
    font-weight: bold;
}


h6 {
    font-size: 13pt;
    line-height: 26px;
    padding: 18px 0 0;
    font-weight: normal;
    font-style: italic;
}

/* horizontal line */

hr {
	border: 0;
	height: 0;
	border-top: 1px solid rgba(0, 0, 0, 0.1);
	border-bottom: 1px solid rgba(255, 255, 255, 0.3);
}

br+br {
    line-height: 0;
    height: 0;
    display: none;
}

p {
    margin: 1em 0;
}

blockquote {
    margin: 14px 0;
    border-left: 4px solid #DDD;
    padding-left: 11px;
    color: #555;
}

pre,
code {
    font-family: 'Bitstream Vera Sans Mono', 'Courier', monospace;
}

pre {
    background-color: #F8F8F8;
    border: 1px solid #CCC;
    font-size: 13px;
    line-height: 19px;
    overflow: auto;
    padding: 6px 10px;
    border-radius: 3px;
    color: black;
}

code {
    margin: 0 2px;
    padding: 2px 5px;
    white-space: nowrap;
    border: 1px solid #CCC;
    background-color: #F8F8F8;
    border-radius: 3px;
    font-size: 12px !important;
}

pre>code {
    margin: 0px;
    padding: 0px;
    white-space: pre;
    border: none;
    background-color: transparent;
    border-radius: 0;
}

a,
a code {
    color: #4183C4;
    text-decoration: none;
}

a:hover,
a code:hover {
    text-decoration: underline;
}

table {
    border-collapse: collapse;
    margin: 20px 0 0;
    padding: 0;
}

table tr th,
table tr td {
    border: 1px solid #CCC;
    text-align: left;
    margin: 0;
    padding: 6px 13px;
}

table tbody tr:nth-child(2n-1) {
    background-color: #F8F8F8;
}
`