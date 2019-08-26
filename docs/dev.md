# Environment Setup

## Go

You must have [go](https://golang.org/doc/install).

## Installing dependencies

Dependencies are handled with [modules](https://github.com/golang/go/wiki/Modules).

## Build

There's a [Makefile](https://github.com/Girbons/comics-downloader/blob/master/Makefile) which should make this step easier.

### CLI version:

```
go build -o comics-downloader ./cmd/downloader
```

### GUI version

[Prerequisites](https://fyne.io//develop/compiling.html) then you can run:

```
go build -o comics-downloader-gui ./cmd/gui
```

if you don't want to install extra dependencies to build the GUI version you could use [fyne-cross](https://github.com/lucor/fyne-cross)
which requires [Docker](https://www.docker.com/get-started).

## Run Tests

```
go test -v ./...
```
