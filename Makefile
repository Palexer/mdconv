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
	GOOS=windows GOARCH=amd64 go build -o dist/win/win-amd64/mdconv.exe

	# compile for mac/darwin
	go generate
	GOOS=darwin GOARCH=amd64 go build -o dist/darwin/darwin-amd64/mdconv
#	go generate
#	GOOS=darwin GOARCH=arm64 go build -o dist/darwin/darwin-arm64/mdconv

	# compile for linux
	go generate
	GOOS=linux GOARCH=amd64 go build -o dist/linux/linux-amd64/mdconv
	go generate
	GOOS=linux GOARCH=arm64 go build -o dist/linux/linux-arm64/mdconv
