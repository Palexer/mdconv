.DEFAULT_GOAL := build
build:
	go generate
	go build -o mdconv

buildall:
	# create directories
	if [ ! -d "dist" ]; then mkdir dist; fi
	if [ ! -d "dist/win" ]; then mkdir dist/win; fi
	if [ ! -d "dist/darwin" ]; then mkdir dist/darwin; fi
	if [ ! -d "dist/linux" ]; then mkdir dist/linux; fi

	# compile for windows
	go generate
	GOOS=windows GOARCH=amd64 go build -o dist/win/mdconv-win-amd64.exe

	# compile for mac/darwin
	go generate
	GOOS=darwin GOARCH=amd64 go build -o dist/darwin/mdconv-darwin-amd64	
#	go generate
#	GOOS=darwin GOARCH=arm64 go build -o dist/darwin/mdconv-darwin-arm64

	# compile for linux
	go generate
	GOOS=linux GOARCH=amd64 go build -o dist/linux/mdconv-linux-amd64
	go generate
	GOOS=linux GOARCH=arm64 go build -o dist/linux/mdconv-linux-arm64
