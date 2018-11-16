package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/jung-kurt/gofpdf"
	"github.com/schollz/progressbar"
)

// Comic contains all the information about a comic.
type Comic struct {
	Name        string
	IssueNumber string
	ImageRegex  string
	Source      string
}

func isUrlValid(url string) bool {
	return !strings.Contains(url, ".gif") && !strings.Contains(url, "logo")
}

func imagesLink(res, regex string) []string {
	re := regexp.MustCompile(regex)
	match := re.FindAllStringSubmatch(res, -1)

	out := make([]string, len(match))

	for i := range out {
		url := match[i][1]
		if isUrlValid(url) {
			out[i] = url
		}
	}
	return out
}

// NewComic return a Comic filled with all the needed data.
func NewComic(name, issueNumber, imageRegex, source string) *Comic {
	return &Comic{
		Name:        name,
		IssueNumber: issueNumber,
		ImageRegex:  imageRegex,
		Source:      source,
	}
}

// MakeComic create the pdf file.
func (c *Comic) MakeComic(url, source string, splittedUrl []string) {
	res, _ := soup.Get(url)

	// retrieve links
	links := imagesLink(res, c.ImageRegex)
	pdf := gofpdf.New("P", "mm", "A4", "")

	bar := progressbar.New(len(links))
	bar.RenderBlank()

	for i, link := range links {
		bar.Add(i)
		if link != "" {
			rsp, err := http.Get(link)

			if err == nil {
				pdf.AddPage()
				tp := pdf.ImageTypeFromMime(rsp.Header["Content-Type"][0])
				pdf.RegisterImageReader(link, tp, rsp.Body)
				pdf.Image(link, 0, 0, 210, 0, false, "", 0, "")
			} else {
				pdf.SetError(err)
			}
		}
	}
	dir, _ := filepath.Abs(fmt.Sprintf("%s/%s/%s/%s/", filepath.Dir(os.Args[0]), "comics", c.Source, c.Name))
	os.MkdirAll(dir, os.ModePerm)

	pdfErr := pdf.OutputFileAndClose(fmt.Sprintf("%s/%s.pdf", dir, c.IssueNumber))
	log.Fatal(pdfErr)
}
