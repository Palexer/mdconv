#!/bin/bash

# generate new binary
make

# create folder for output files
if [ ! -d "test_output" ]; then mkdir test_output; fi

# ususal pdf/html convertions
./mdconv -o test_output/main_test.html testdata/main_test.md
./mdconv -o test_output/main_test.pdf testdata/main_test.md

# custom and default CSS
./mdconv -o test_output/custom_test.html -style testdata/custom.css testdata/main_test.md
./mdconv -o test_output/custom_test.pdf -style testdata/custom.css testdata/main_test.md

# only custom CSS
./mdconv -o test_output/overwrite_test.html -style testdata/custom.css -overwrite testdata/main_test.md
./mdconv -o test_output/overwrite_test.pdf -style testdata/custom.css -overwrite testdata/main_test.md

# remove binary
rm mdconv
