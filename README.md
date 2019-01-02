# Comics Downloader (Go version)

[![Build Status](https://travis-ci.org/Girbons/comics-downloader.svg?branch=master)](https://travis-ci.org/Girbons/comics-downloader)
[![Coverage Status](https://coveralls.io/repos/github/Girbons/comics-downloader/badge.svg?branch=master)](https://coveralls.io/github/Girbons/comics-downloader?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Girbons/go-comics-downloader)](https://goreportcard.com/report/github.com/Girbons/go-comics-downloader)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Supported Sites

- http://www.comicextra.com/
- http://www.mangahere.cc/
- https://mangarock.com/
- https://www.mangareader.net/

## Getting Started

### Installing

Download the latest release:

- [Linux](https://github.com/Girbons/comics-downloader/releases/download/v0.5.0/comics-downloader)
- [OSX](https://github.com/Girbons/comics-downloader/releases/download/v0.5.0/comics-downloader-osx)
- [Windows](https://github.com/Girbons/comics-downloader/releases/download/v0.5.0/comics-downloader.exe)

Put the script under a folder.

## Usage

<img src="img/usage.gif?raw=true" />

### Download as EPUB

```bash
./comics-downloader -url=[your url] -format=epub
```

## Built With

- [go](https://github.com/golang/go)
- [gofpdf](https://github.com/jung-kurt/gofpdf)
- [go-epub](http://github.com/bmaupin/go-epub)
- [soup](https://github.com/anaskhan96/soup)
- [progressbar](https://github.com/schollz/progressbar)
- [logrus](https://github.com/sirupsen/logrus)
- [mri](https://github.com/BakeRolls/mri/blob/master/mri.go)

## Contribuiting

Feel free to submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
