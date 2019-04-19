#get latest git tag
TAG = $$(git describe --abbrev=0 --tags)
#get current date and time
NOW = $$(date +'%Y-%m-%d_%T')
#use the linker to inject the informations
LDFLAGS = -ldflags="-X main.release=${TAG} -X main.buildTime=${NOW}"

help:                              # this command
	# [generating help from tasks header]
	@egrep '^[A-Za-z0-9_-]+:' Makefile

osx-build: # Create the executable for OSX
	@GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o build/comics-downloader-osx ./cmd/downloader

windows-build: # Create the executable for WINDOWS
	@GOOS=windows GOARCH=amd64 go build ${LDFLAGS}  -o build/comics-downloader.exe ./cmd/downloader

linux-build: # Create the executable for LINUX
	@GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o build/comics-downloader ./cmd/downloader

gui-builds: # Creates GUI executables for LINUX OSX WINDOWS
	fyne-cross --output comics-downloader-gui --targets=linux/amd64,windows/amd64,darwin/amd64 ./cmd/gui

builds: # Create the executables for OSX/Windows/Linux
	@make osx-build
	@make windows-build
	@make linux-build
	@make gui-builds

remove-builds: # Remove executables
	@rm -rf build/
