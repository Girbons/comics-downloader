package core

import "errors"

type ComicOutputFormat string

func (c ComicOutputFormat) String() string {
	return string(c)
}

// manga output format supported
const (
	CBR  ComicOutputFormat = "cbr"
	CBZ  ComicOutputFormat = "cbz"
	EPUB ComicOutputFormat = "epub"
	PDF  ComicOutputFormat = "pdf"
)

var InvalidOutputFormatError = errors.New("Invalid output format")

func ToComicOutputFormat(format string) (ComicOutputFormat, error) {
	switch format {
	case "cbr":
		return CBR, nil
	case "cbz":
		return CBZ, nil
	case "epub":
		return EPUB, nil
	case "pdf":
		return PDF, nil
	default:
		return "", InvalidOutputFormatError
	}
}
