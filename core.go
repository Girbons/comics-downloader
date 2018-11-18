package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/anaskhan96/soup"
	"github.com/jung-kurt/gofpdf"
	"github.com/schollz/progressbar"
	log "github.com/sirupsen/logrus"
)

// Comic struct all the information ab
type Comic struct {
	Name        string
	IssueNumber string
	ImageRegex  string
	Source      string
}

// imagesLink returns an array containing the image link
func imagesLink(source, response, regex string) []string {
	log.Debug("Retrieving the images link for ", source)

	re := regexp.MustCompile(regex)

	switch source {
	case "www.comicextra.com":
		match := re.FindAllStringSubmatch(response, -1)
		out := make([]string, len(match))

		for i := range out {
			url := match[i][1]
			if IsUrlValid(url) {
				out[i] = url
			}
		}
		return out
	case "www.mangahere.cc":

		doc := soup.HTMLParse(response)
		links := doc.FindAll("option")
		pageLinks := make([]string, 1)
		imgLinks := make([]string, 1)

		for _, link := range links {
			newLink := fmt.Sprintf("http://%s", link.Attrs()["value"][2:])
			if !CheckValueInSlice(newLink, pageLinks) {
				pageLinks = append(pageLinks, newLink)
			}
		}

		for _, link := range pageLinks {
			if link != "" {
				resp, err := soup.Get(link)

				if err != nil {
					log.Error("Something went wrong: ", err)
				}

				match := re.FindAllStringSubmatch(resp, -1)
				for _, f := range match {
					if IsUrlValid(f[1]) {
						imgLinks = append(imgLinks, f[1])
					}
				}
			}
		}
		return imgLinks
	default:
		links := make([]string, 1)
		return links
	}
}

// NewComic create a Comic instance
func NewComic(name, issueNumber, imageRegex, source string) *Comic {
	return &Comic{
		Name:        name,
		IssueNumber: issueNumber,
		ImageRegex:  imageRegex,
		Source:      source,
	}
}

// MakeComic create the pdf file
func (c *Comic) MakeComic(url, source string, splittedUrl []string) {
	log.Debug("Source Detected the download will start soon")
	response, soupErr := soup.Get(url)

	if soupErr != nil {
		log.Error("Ops something went wrong: ", soupErr)
	}
	// retrieve links
	imgLinks := imagesLink(source, response, c.ImageRegex)
	// setup the pdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	// setup the progress bar
	bar := progressbar.New(len(imgLinks))
	// this will show up the progress bar since the beginning
	bar.RenderBlank()
	// for each link get the image
	// to add to the pdf file
	for i, link := range imgLinks {
		rsp, err := http.Get(link)
		if link != "" {
			if err == nil {
				// add a new PDF page
				pdf.AddPage()
				// retrieve the image format from the response header (jpeg, png...)
				tp := pdf.ImageTypeFromMime(rsp.Header["Content-Type"][0])
				// The image is directly added to the pdf without being saved to the disk
				pdf.RegisterImageReader(link, tp, rsp.Body)
				// Here we set the image position on the pdf page
				pdf.Image(link, 0, 0, 210, 0, false, "", 0, "")
				// increase the progressbar
				bar.Add(i)
			} else {
				pdf.SetError(err)
			}
		}
	}
	// Set progressbar to maximum
	bar.Finish()
	// this will create the path where the file will be added
	dir, dirErr := filepath.Abs(fmt.Sprintf("%s/%s/%s/%s/", filepath.Dir(os.Args[0]), "comics", c.Source, c.Name))

	if dirErr != nil {
		log.Error("There was an error: ", dirErr)
	}
	// create all the directory and the pdf file
	mkdirErr := os.MkdirAll(dir, os.ModePerm)
	if mkdirErr != nil {
		log.Error("There was an error while creating the needed folders: ", mkdirErr)
	}

	// Save the pdf file
	err := pdf.OutputFileAndClose(fmt.Sprintf("%s/%s.pdf", dir, c.IssueNumber))
	if err != nil {
		log.Error("There was an error while making the PDF: ", err)
	}

	fmt.Printf("\n")
	log.Info("Download Completed")
}
