# List available recipes
default:
    @just --list

# Creates Mac OSX ARM binary
osx-build-arm:
    GOOS=darwin go build -o build/comics-downloader-osx-arm ./cmd/downloader

# Creates Mac OSX x86-64 binary
osx-build-x86-64:
    GOOS=darwin GOARCH=amd64 go build -o build/comics-downloader-osx-x86-64 ./cmd/downloader

# Creates Windows x86-64 binary
windows-x86-64-build:
    GOOS=windows GOARCH=amd64 go build -o build/comics-downloader-win-x86-64.exe ./cmd/downloader

# Creates Windows 386 binary
windows-386-build:
    GOOS=windows GOARCH=386 go build -o build/comics-downloader-win-386.exe ./cmd/downloader

# Creates Linux x86-64 binary
linux-x86-64-build:
    GOOS=linux GOARCH=amd64 go build -o build/comics-downloader-linux-x86-64 ./cmd/downloader

# Creates Linux 386 binary
linux-386-build:
    GOOS=linux GOARCH=386 go build -o build/comics-downloader-linux-386 ./cmd/downloader

# Creates Linux ARM binary
linux-arm-build:
    GOOS=linux GOARCH=arm go build -o build/comics-downloader-linux-arm ./cmd/downloader

# Creates Linux ARM64 binary
linux-arm64-build:
    GOOS=linux GOARCH=arm64 go build -o build/comics-downloader-linux-arm64 ./cmd/downloader

# Creates OSX GUI binary
osx-gui-build:
    GOOS=darwin go build -o build/comics-downloader-gui-osx ./cmd/gui

# Creates Windows GUI executable
windows-gui-build:
    fyne-cross windows -output comics-downloader-gui-windows.exe ./cmd/gui

# Creates Linux GUI executable
linux-gui-build:
    fyne-cross linux -output comics-downloader-gui ./cmd/gui

# Creates executables for OSX/Windows/Linux
builds: linux-386-build linux-arm-build linux-arm64-build linux-x86-64-build osx-build-arm osx-build-x86-64 windows-x86-64-build windows-386-build osx-gui-build linux-gui-build windows-gui-build

# Remove build artifacts
remove-builds:
    rm -rf build/

# Run static analysis
lint:
    golangci-lint run

# Run GUI unit tests (headless, no display required)
test-gui:
    go test -v ./cmd/gui/...

# Run the GUI
run-gui:
    go run ./cmd/gui
