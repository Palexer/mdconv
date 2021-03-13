.DEFAULT_GOAL := build
build:
	go generate
	go build -o mdconv

install:
	sudo mkdir -p /usr/local/bin
	sudo mkdir -p /usr/local/share/man/man1/
	sudo cp doc/mdconv.man /usr/local/share/man/man1/mdconv.1
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
	go generate
	GOOS=windows GOARCH=amd64 go build -o dist/win/amd64/mdconv.exe

	# compile for mac/darwin
	go generate
	GOOS=darwin GOARCH=amd64 go build -o dist/darwin/amd64/mdconv
	go generate
	GOOS=darwin GOARCH=arm64 go build -o dist/darwin/arm64/mdconv

	# compile for linux
	go generate
	GOOS=linux GOARCH=amd64 go build -o dist/linux/amd64/mdconv
	go generate
	GOOS=linux GOARCH=arm64 go build -o dist/linux/arm64/mdconv
