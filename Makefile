.DEFAULT_GOAL := build
build:
	@go build -ldflags "-s -w" -o mdconv 

install: build
	@sudo mkdir -p /usr/local/bin
	@sudo mkdir -p /usr/local/share/man/man1/
	@sudo cp doc/mdconv.1 /usr/local/share/man/man1/mdconv.1
	@sudo mv mdconv /usr/local/bin

clean:
	@echo "cleaning"
	@if [ -f "mdconv" ]; then rm mdconv; fi
	@if [ -f "out.html" ]; then rm out.html; fi
	@if [ -f "out.pdf" ]; then rm out.pdf; fi
	@if [ -d "dist" ]; then rm -r dist/; fi
	@if [ -d "testoutput" ]; then rm -r testoutput/; fi

test: build
	@mkdir -p testoutput

	@./mdconv -o testoutput/main_test.html testdata/main_test.md
	@./mdconv -o testoutput/main_test.pdf testdata/main_test.md

	@if [ -f "mdconv" ]; then rm mdconv; fi

testall: build
	# create folders
	@echo create folder for output files
	@mkdir -p testoutput/html
	@mkdir -p testoutput/pdf

	# HTML
	@echo HTML tests

	# main test
	@./mdconv -o testoutput/html/main_test.html testdata/main_test.md
	# custom and default CSS
	@./mdconv -o testoutput/html/custom_test.html -css testdata/custom.css testdata/main_test.md
	# only custom CSS
	@./mdconv -o testoutput/html/overwrite_test.html -css testdata/custom.css -overwrite testdata/main_test.md
	# no style
	@./mdconv -o testoutput/html/nostyle_test.html -overwrite testdata/main_test.md

	# PDF
	@echo PDF tests

	# main test
	@./mdconv -o testoutput/pdf/main_test.pdf testdata/main_test.md
	# custom and default CSS
	@./mdconv -o testoutput/pdf/custom_test.pdf -css testdata/custom.css testdata/main_test.md
	# only custom CSS
	@./mdconv -o testoutput/pdf/overwrite_test.pdf -css testdata/custom.css -overwrite testdata/main_test.md
	# no style
	@./mdconv -o testoutput/pdf/nostyle_test.pdf -overwrite testdata/main_test.md

	# remove binary
	@echo removing binary
	@rm mdconv

buildall:
	@echo creating directories
	@mkdir -p dist/win/amd64
	@mkdir -p dist/darwin/amd64
	@mkdir -p dist/darwin/arm64
	@mkdir -p dist/linux/amd64
	@mkdir -p dist/linux/arm64

	@echo start building

	@GOOS=windows GOARCH=amd64 go build -o dist/win/amd64/mdconv.exe -ldflags "-s -w"

	@GOOS=darwin GOARCH=amd64 go build -o dist/darwin/amd64/mdconv -ldflags "-s -w"
	@GOOS=darwin GOARCH=arm64 go build -o dist/darwin/arm64/mdconv -ldflags "-s -w"

	@GOOS=linux GOARCH=amd64 go build -o dist/linux/amd64/mdconv -ldflags "-s -w"
	@GOOS=linux GOARCH=arm64 go build -o dist/linux/arm64/mdconv -ldflags "-s -w"

	@echo zipping binaries
	@zip -r dist/win/amd64/mdconv-win-amd64.zip dist/win/amd64/mdconv.exe doc/mdconv.1

	@zip -r dist/darwin/amd64/mdconv-darwin-amd64.zip dist/darwin/amd64/mdconv doc/mdconv.1
	@zip -r dist/darwin/arm64/mdconv-darwin-arm64.zip dist/darwin/arm64/mdconv doc/mdconv.1

	@zip -r dist/linux/amd64/mdconv-linux-amd64.zip dist/linux/amd64/mdconv doc/mdconv.1
	@zip -r dist/linux/arm64/mdconv-linux-arm64.zip dist/linux/arm64/mdconv doc/mdconv.1
