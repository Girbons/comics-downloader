help: # this command
	# [generating help from tasks header]
	@egrep '^[A-Za-z0-9_-]+:' Makefile

osx-build: # Creates Mac OSX
	@GOOS=darwin go build -o build/comics-downloader-osx ./cmd/downloader

windows-build: # Creates Windows
	@GOOS=windows go build -o build/comics-downloader.exe ./cmd/downloader

linux-build: # Creates Linux
	@GOOS=linux go build -o build/comics-downloader ./cmd/downloader

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
	@make osx-build
	@make windows-build
	@make linux-build
	@make linux-arm-build
	@make linux-arm64-build
	@make osx-gui-build
	@make windows-gui-build
	@make linux-gui-build

remove-builds: # Remove executables
	@rm -rf build/
