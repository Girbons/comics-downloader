# Comics Downloader

[![CircleCI](https://circleci.com/gh/Girbons/comics-downloader/tree/master.svg?style=svg)](https://circleci.com/gh/Girbons/comics-downloader/tree/master)
[![Coverage Status](https://img.shields.io/coveralls/github/Girbons/comics-downloader.svg?style=flat-square)](https://coveralls.io/github/Girbons/comics-downloader?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Girbons/comics-downloader)](https://goreportcard.com/report/github.com/Girbons/comics-downloader)
[![Github All Releases](https://img.shields.io/github/downloads/Girbons/comics-downloader/total.svg?style=flat-square)]()
[![Release](https://img.shields.io/github/release/Girbons/comics-downloader.svg?style=flat-square)](https://github.com/Girbons/comics-downlowader/releases/latest)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)

## Supported Sites

- http://www.comicextra.com/
- https://mangarock.com/
- https://www.mangareader.net/
- http://www.mangatown.com/

## Getting Started

### Installing

Download the latest release:

- [Linux](https://github.com/Girbons/comics-downloader/releases/download/v0.11/comics-downloader)
- [OSX](https://github.com/Girbons/comics-downloader/releases/download/v0.11/comics-downloader-osx)
- [Windows](https://github.com/Girbons/comics-downloader/releases/download/v0.11/comics-downloader.exe)

Put the script under a folder.

## Usage

<img src="img/usage.gif?raw=true" />

You can invoke the `--help`:

```
  -all
        Download all issues of the Comic or Comics
  -country string
        Set the country to retrieve a manga, Used by MangaRock
  -format string
        Comic format output, supported formats are pdf,epub,cbr,cbz (default "pdf")
  -url string
        Comic URL or Comic URLS by separating each site with a comma without the use of spaces
  -version
        Display build date and release informations
```

### Multiple URLs

Without `url` parameter:

```bash
./comics-downloader url1 url2 url3
```

With `url` parameter:

```bash
./comics-downloader -url=url,url2,url3
```

### Specify the format output

available formats:

- pdf
- epub
- cbr
- cbz

Default format is __pdf__.

example:

```bash
./comics-downloader -url=[your url] -format=epub
```

### Download the whole comic

Provide the comic url and use the `-all` flag. The url provided can be any issue of the comic, or the main comic page url.

example:

```bash
./comics-downloader -url=[your url] -all
```

## Config file

To avoid to specify everytime the output format you can create a `config.yaml` file in the same path of the executable.

Add the string below and substitute `"cbr"` with your favourite format.

**NOTE**: if `--format` or `-format` is specified, the value in `config.yaml` will be ignored.

```
default_output_format: "cbr"
```

## Built With

- [go](https://github.com/golang/go)
- [gofpdf](https://github.com/jung-kurt/gofpdf)
- [go-epub](http://github.com/bmaupin/go-epub)
- [soup](https://github.com/anaskhan96/soup)
- [progressbar](https://github.com/schollz/progressbar)
- [logrus](https://github.com/sirupsen/logrus)
- [mri](https://github.com/BakeRolls/mri/blob/master/mri.go)
- [archiver](https://github.com/mholt/archiver)
- [viper](https://github.com/spf13/viper)

## Contribuiting

Feel free to submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
