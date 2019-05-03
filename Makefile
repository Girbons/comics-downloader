#get latest git tag
TAG = $$(git describe --tags)
#get current date and time
NOW = $$(date +'%Y-%m-%d_%T')
#use the linker to inject the informations
LDFLAGS = -ldflags="-X main.release=${TAG} -X main.buildTime=${NOW}"

help:                              # this command
	# [generating help from tasks header]
	@egrep '^[A-Za-z0-9_-]+:' Makefile

osx-build: # Create the executable for OSX
	@GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o comics-downloader-osx ./cmd/

windows-build: # Create the executable for WINDOWS
	@GOOS=windows GOARCH=amd64 go build ${LDFLAGS}  -o comics-downloader.exe ./cmd/

linux-build: # Create the executable for LINUX
	@GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o comics-downloader ./cmd/

builds: # Create the executables for OSX/Windows/Linux
	@make osx-build
	@make windows-build
	@make linux-build

remove-builds: # Remove executables
	@rm comics-downloader comics-downloader.exe comics-downloader-osx
