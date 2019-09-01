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

- [Linux](https://github.com/Girbons/comics-downloader/releases/download/v0.15/comics-downloader)
- [Mac OSX](https://github.com/Girbons/comics-downloader/releases/download/v0.15/comics-downloader-osx)
- [Windows](https://github.com/Girbons/comics-downloader/releases/download/v0.15/comics-downloader.exe)
- [Linux ARM](https://github.com/Girbons/comics-downloader/releases/download/v0.15/comics-downloader-linux-arm)
- [Linux ARM64](https://github.com/Girbons/comics-downloader/releases/download/v0.15/comics-downloader-linux-arm64)

Download the latest GUI release:

- [Linux](https://github.com/Girbons/comics-downloader/releases/download/v0.15/comics-downloader-gui-linux-amd64)
- [Mac OSX](https://github.com/Girbons/comics-downloader/releases/download/v0.15/comics-downloader-gui-osx)
- [Windows](https://github.com/Girbons/comics-downloader/releases/download/v0.15/comics-downloader-gui-windows.exe)

Put the script under a folder.

## Usage

<img src="img/usage.gif?raw=true" />

You can invoke the `--help`:

```
  -all
        Download all issues of the Comic or Comics
  -country string
        Set the country to retrieve a manga, Used by MangaRock
  -deamon
        Run the download as deamon
  -format string
        Comic format output, supported formats are pdf,epub,cbr,cbz (default "pdf")
  -images-format
        To use with images-only flag, choose the image format, available png,jpeg,img (default "jpg")
  -images-only
        Download comic/manga images
  -last
        Download the last Comic issue
  -timeout int
        Timeout (seconds), specifies how often the downloader runs (default 600)
  -url string
        Comic URL or Comic URLS by separating each site with a comma without the use of spaces
  -version
        Display build date and release informations
```

### Checking for mangas using a Raspberry Pi

If you'd like to track your favourite mangas you can use this bash [script](https://gist.github.com/nestukh/5397b836c8e5f34f6feb4ec4efe6b86a).

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

### Download the last comic issue

To download the last comic issue use `-last` flag.

```bash
./comics-downloader -url=[your url] -last
```

### Run as Deamon

You can run the CLI downloader as deamon using `-deamon` flag.
works only if `-all` or `-last` flags are specified.

```bash
./comics-downloader -url=[your url] -deamon
```

You can customize the deamon timeout using the `-timeout` flag.

```bash
./comics-downloader -url=[your url] -deamon -timeout=300
```

### Download Only the Images

You can download only the images using `-images-only` flag.

```bash
./comics-downloader -url=[your url] -images-only
```

To choose the format use `-images-format` flag, the available formats are:

* img
* png
* jpg

Default is __jpg__.

```bash
./comics-downloader -url=[your url] -images-only -images-format=jpg
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

## Contribuiting

Feel free to submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
