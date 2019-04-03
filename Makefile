help:                              # this command
	# [generating help from tasks header]
	@egrep '^[A-Za-z0-9_-]+:' Makefile

osx-build: # Create the executable for OSX
	@GOOS=darwin GOARCH=amd64 go build -o comics-downloader-osx ./cmd/comics-downloader/

windows-build: # Create the executable for WINDOWS
	@GOOS=windows GOARCH=amd64 go build -o comics-downloader.exe ./cmd/comics-downloader/

linux-build: # Create the executable for LINUX
	@GOOS=linux GOARCH=amd64 go build -o comics-downloader ./cmd/comics-downloader/

builds: # Create the executables for OSX/Windows/Linux
	@make osx-build
	@make windows-build
	@make linux-build

remove-builds: # Remove executables
	@rm comics-downloader comics-downloader.exe comics-downloader-osx
