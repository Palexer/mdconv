.DEFAULT_GOAL := build
build:
	go build -o mdconv 

install:
	sudo mkdir -p /usr/local/bin
	sudo mkdir -p /usr/local/share/man/man1/
	sudo cp doc/mdconv.1 /usr/local/share/man/man1/mdconv.1
	sudo mv mdconv /usr/local/bin

clean:
	if [ -f "mdconv" ]; then rm mdconv; fi
	if [ -d "dist" ]; then rm -r dist/; fi
	if [ -d "testoutput" ]; then rm -r testoutput/; fi

test: build
	if [ ! -d "testoutput" ]; then mkdir testoutput; fi

	# ususal pdf/html convertions
	./mdconv -o testoutput/main_test.html testdata/main_test.md
	./mdconv -o testoutput/main_test.pdf testdata/main_test.md

	if [ -f "mdconv" ]; then rm mdconv; fi

testall: build
	# create folder for output files
	if [ ! -d "testoutput" ]; then mkdir testoutput; fi

	# ususal pdf/html convertions
	./mdconv -o testoutput/main_test.html testdata/main_test.md
	./mdconv -o testoutput/main_test.pdf testdata/main_test.md

	# custom and default CSS
	./mdconv -o testoutput/custom_test.html -c testdata/custom.css testdata/main_test.md
	./mdconv -o testoutput/custom_test.pdf -c testdata/custom.css testdata/main_test.md

	# only custom CSS
	./mdconv -o testoutput/overwrite_test.html -c testdata/custom.css -overwrite testdata/main_test.md
	./mdconv -o testoutput/overwrite_test.pdf -c testdata/custom.css -overwrite testdata/main_test.md

	# no style
	./mdconv -o testoutput/nostyle_test.html -overwrite testdata/main_test.md
	./mdconv -o testoutput/nostyle_test.pdf -overwrite testdata/main_test.md

	# custom fonts (only HTML tests)
	./mdconv -f sans -o testoutput/font_sans.html testdata/main_test.md
	./mdconv -f serif -o testoutput/font_serif.html testdata/main_test.md
	./mdconv -f monospace -o testoutput/font_monospace.html testdata/main_test.md

	# remove binary
	rm mdconv

buildall:
	# create directories
	mkdir -p dist/win/amd64
	mkdir -p dist/darwin/amd64
	mkdir -p dist/darwin/arm64
	mkdir -p dist/linux/amd64
	mkdir -p dist/linux/arm64

	# compile for windows
	GOOS=windows GOARCH=amd64 go build -o dist/win/amd64/mdconv.exe -ldflags "-s -w"

	# compile for mac/darwin
	GOOS=darwin GOARCH=amd64 go build -o dist/darwin/amd64/mdconv -ldflags "-s -w"
	GOOS=darwin GOARCH=arm64 go build -o dist/darwin/arm64/mdconv -ldflags "-s -w"

	# compile for linux
	GOOS=linux GOARCH=amd64 go build -o dist/linux/amd64/mdconv -ldflags "-s -w"
	GOOS=linux GOARCH=arm64 go build -o dist/linux/arm64/mdconv -ldflags "-s -w"

	# zip binaries
	# windows
	zip -r dist/win/amd64/mdconv-win-amd64.zip dist/win/amd65/mdconv/exe doc/mdconv.1

	# mac/darwin
	zip -r dist/darwin/amd64/mdconv-darwin-amd64.zip dist/darwin/amd64/mdconv doc/mdconv.1
	zip -r dist/darwin/arm64/mdconv-darwin-arm64.zip dist/darwin/arm64/mdconv doc/mdconv.1

	# linux
	zip -r dist/linux/amd64/mdconv-linux-amd64.zip dist/linux/amd64/mdconv doc/mdconv.1
	zip -r dist/linux/arm64/mdconv-linux-arm64.zip dist/linux/arm64/mdconv doc/mdconv.1
