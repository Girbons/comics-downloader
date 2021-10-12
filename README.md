# Comics Downloader

[![Build Status](https://app.travis-ci.com/Girbons/comics-downloader.svg?branch=master)](https://app.travis-ci.com/Girbons/comics-downloader)
[![Coverage Status](https://img.shields.io/coveralls/github/Girbons/comics-downloader.svg?style=flat-square)](https://coveralls.io/github/Girbons/comics-downloader?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Girbons/comics-downloader)](https://goreportcard.com/report/github.com/Girbons/comics-downloader)
[![Github All Releases](https://img.shields.io/github/downloads/Girbons/comics-downloader/total.svg?style=flat-square)]()
[![Release](https://img.shields.io/github/release/Girbons/comics-downloader.svg?style=flat-square)](https://github.com/Girbons/comics-downloader/releases/latest)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)

## Supported Sites

- https://www.comicextra.com/
- https://readcomiconline.li/
- https://www.mangareader.tv/
- https://www.mangatown.com/
- https://mangadex.org/
- https://mangakakalot.com/
- https://manganato.com/

## Getting Started

### Installing

Download the latest release:

- [Linux](https://github.com/Girbons/comics-downloader/releases/download/v0.32.0/comics-downloader)
- [Mac OSX](https://github.com/Girbons/comics-downloader/releases/download/v0.32.0/comics-downloader-osx)
- [Windows](https://github.com/Girbons/comics-downloader/releases/download/v0.32.0/comics-downloader.exe)
- [Linux ARM](https://github.com/Girbons/comics-downloader/releases/download/v0.32.0/comics-downloader-linux-arm)
- [Linux ARM64](https://github.com/Girbons/comics-downloader/releases/download/v0.32.0/comics-downloader-linux-arm64)

Download the latest GUI release:

- [Linux](https://github.com/Girbons/comics-downloader/releases/download/v0.32.0/comics-downloader-gui)
- [Mac OSX](https://github.com/Girbons/comics-downloader/releases/download/v0.32.0/comics-downloader-gui-osx)
- [Windows](https://github.com/Girbons/comics-downloader/releases/download/v0.32.0/comics-downloader-gui-windows.exe)

## Usage

<img src="img/usage.gif?raw=true" />

You can invoke the `--help`:

```
Usage:
  -all
        Download all issues of the Comic or Comics
  -country string
        Set the country to retrieve a manga, Used by MangaDex which uses ISO 3166-1 codes
  -create-default-path comics/[source]/[name]/
        Using this flag your comics/issue will be downloaded without prepending the default folder structure, comics/[source]/[name]/ (default true)
  -custom-comic-name string
        Use a custom name for the comic output.
  -daemon
        Run the download as daemon
  -daemon-timeout int
        DaemonTimeout (seconds), specifies how often the downloader runs (default 600)
  -debug
    	Shows Debug log
  -format string
        Comic format output, supported formats are pdf,epub,cbr,cbz (default "pdf")
  -force-aspect
        Force images to A4 Portrait aspect ratio
  -images-format
        To use with images-only flag, choose the image format, available png,jpeg,img (default "jpg")
  -images-only
        Download comic/manga images
  -issue-number-only
        Force only saving with issue number instead of chapter name + issue number.
  -last
        Download the last Comic issue
  -output string
        Folder where the comics will be saved
  -range
        Range of issues to download, example 3-9
  -url string
        Comic URL or Comic URLS by separating each site with a comma without the use of spaces
  -version
        Display release version
```

## Options supported

| Source                      | all      | country  | last     |
| --------------------------- | -------- | -------- | -------- |
| http://readallcomics.com    | &#x2713; | &#x2717; | &#x2713; |
| http://www.comicextra.com/  | &#x2713; | &#x2717; | &#x2713; |
| http://www.mangatown.com/   | &#x2713; | &#x2717; | &#x2713; |
| https://mangadex.org/       | &#x2713; | &#x2713; | &#x2717; |
| https://readcomiconline.li/ | &#x2713; | &#x2717; | &#x2713; |
| https://www.mangareader.tv/ | &#x2713; | &#x2717; | &#x2713; |
| https://www.mangakalot.com/ | &#x2713; | &#x2717; | &#x2713; |
| https://www.manganato.com/  | &#x2713; | &#x2717; | &#x2713; |

### Checking for mangas using a Raspberry Pi

If you'd like to track your favourite mangas you can use this bash [script](https://gist.github.com/nestukh/5397b836c8e5f34f6feb4ec4efe6b86a).

### Multiple URLs

```bash
./comics-downloader -url=url,url2,url3
```

### Specify the format output

available formats:

- pdf
- epub
- cbr
- cbz

Default format is **pdf**.

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

### Download the range of issues

Provide the comic url and use the `-range` flag. The url provided can be any issue of the comic, or the main comic page url.

example:

```bash
./comics-downloader -url=[your url] -range=[start-end]
```

### Download the last comic issue

To download the last comic issue use `-last` flag.

```bash
./comics-downloader -url=[your url] -last
```

### Download to custom folder

To download to a custom folder use the `-output` flag.
The folder will be created if not already existing.

```bash
./comics-downloader -url=[your url] -output=[your path]
```

### Run as daemon

You can run the CLI downloader as daemon using `-daemon` flag.
works only if `-all` or `-last` flags are specified.

```bash
./comics-downloader -url=[your url] -daemon
```

You can customize the daemon timeout using the `-daemon-timeout` flag.

```bash
./comics-downloader -url=[your url] -daemon -daemon-timeout=300
```

### Download Only the Images

You can download only the images using `-images-only` flag.

```bash
./comics-downloader -url=[your url] -images-only
```

To choose the format use `-images-format` flag, the available formats are:

- img
- png
- jpg

Default is **jpg**.

```bash
./comics-downloader -url=[your url] -images-only -images-format=jpg
```

### Avoid Default Folder Structure

The default folder structure that will be created is: `/comics/[source]/[name]/`.
To avoid that use `-create-default-path` flag.

```bash
./comics-downloader -url=[your url] -create-default-path=false
```

## Built With

- [go](https://github.com/golang/go)
- [gofpdf](https://github.com/jung-kurt/gofpdf)
- [go-epub](http://github.com/bmaupin/go-epub)
- [soup](https://github.com/anaskhan96/soup)
- [progressbar](https://github.com/schollz/progressbar)
- [logrus](https://github.com/sirupsen/logrus)
- [mangadex](https://github.com/bake/mangadex)
- [archiver](https://github.com/mholt/archiver)
- [regexp2](https://github.com/dlclark/regexp2)

## Contributing

Feel free to submit a pull request, a guide to setup the development environment is available [here](docs/dev.md)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
