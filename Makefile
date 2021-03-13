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
	if [ -d "test_output" ]; then rm -r test_output/; fi

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
