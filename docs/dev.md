# Environment Setup

## Go

You must have:

- [go](https://golang.org/doc/install)
- GCC or [MinGW](http://tdm-gcc.tdragon.net/download)

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

### Optional: libjpeg backend for JPEG encode/decode

Comics Downloader supports an optional JPEG backend using [go-libjpeg](https://github.com/pixiv/go-libjpeg).
The default behavior is to use the Go standard library unless you pass the `libjpeg` build tag.

Additional requirements for `libjpeg` builds:

- libjpeg development headers and libraries (for example `libjpeg-dev` or `libjpeg-turbo-devel`)

Examples:

```bash
# CLI
go build -tags libjpeg -o comics-downloader ./cmd/downloader

# GUI
go build -tags libjpeg -o comics-downloader-gui ./cmd/gui
```

Using Makefile targets with build tags:

```bash
# Build all release artifacts using libjpeg backend
make builds BUILD_TAGS=libjpeg

# Build a specific target with libjpeg backend
make linux-x86-64-build BUILD_TAGS=libjpeg
```

## Run Tests

```
go test -v ./...
```

## Lint

Install golangci-lint once with:

```
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
```

Then run the full static analysis suite with:

```
golangci-lint run
```

### Cloudflare / Request Tweaks

Some sources require browser-like fingerprints. The CLI accepts:

```
--user-agents "UA1,UA2"          # rotate these agents per request
--session-cookie "cf_clearance=...; other=value"
```

Capture values from a working browser session when needed.
