help: # this command
	# [generating help from tasks header]
	@egrep '^[A-Za-z0-9_-]+:' Makefile

osx-build-arm: # Creates Mac OSX
	@GOOS=darwin go build -o build/comics-downloader-osx-arm ./cmd/downloader

osx-build-x86-64: # Creates Mac OSX
	@GOOS=darwin GOARCH=amd64 go build -o build/comics-downloader-osx-x86-64 ./cmd/downloader

windows-x86-64-build: # Creates Windows
	@GOOS=windows GOARCH=amd64 go build -o build/comics-downloader-win-x86-64.exe ./cmd/downloader

windows-386-build: # Creates Windows
	@GOOS=windows GOARCH=386 go build -o build/comics-downloader-win-386.exe ./cmd/downloader

linux-x86-64-build: # Creates Linux
	@GOOS=linux GOARCH=amd64 go build -o build/comics-downloader-linux-x86-64 ./cmd/downloader

linux-386-build:
	@GOOS=linux GOARCH=386 go build -o build/comics-downloader-linux-386 ./cmd/downloader

linux-arm-build: # Creates Linux ARM
	@GOOS=linux GOARCH=arm go build -o build/comics-downloader-linux-arm ./cmd/downloader

linux-arm64-build: # Creates Linux ARM64
	@GOOS=linux GOARCH=arm64 go build -o build/comics-downloader-linux-arm64 ./cmd/downloader

osx-gui-build: # Creates osx GUI
	@GOOS=darwin go build -o build/comics-downloader-gui-osx ./cmd/gui

windows-gui-build: # Creates Window GUI executable
	@fyne-cross windows -output comics-downloader-gui-windows.exe ./cmd/gui

linux-gui-build: # Creates Linux Gui executable
	@fyne-cross linux -output comics-downloader-gui ./cmd/gui

builds: # Creates executables for OSX/Windows/Linux
	@make linux-386-build
	@make linux-arm-build
	@make linux-arm64-build
	@make linux-x86-64-build
	@make osx-build-arm
	@make osx-build-x86-64
	@make windows-x86-64-build
	@make windows-386-build
	@make osx-gui-build
	@make linux-gui-build
	@make windows-gui-build

remove-builds: # Remove executables
	@rm -rf build/
