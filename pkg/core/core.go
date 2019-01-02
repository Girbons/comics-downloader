package core

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BakeRolls/mri"
	"github.com/Girbons/comics-downloader/pkg/util"
	epub "github.com/bmaupin/go-epub"
	"github.com/jung-kurt/gofpdf"
	"github.com/schollz/progressbar"
	log "github.com/sirupsen/logrus"
)

// Comic struct all the information ab
type Comic struct {
	Author      string
	Name        string
	IssueNumber string
	ImageRegex  string
	Source      string
	URLSource   string
	Links       []string
	Options     map[string]string
}

// SetAuthor sets the comic author
func (c *Comic) SetAuthor(author string) {
	c.Author = author
}

// SetName sets the comic name
func (c *Comic) SetName(name string) {
	c.Name = name
}

// SetIssueNumber sets the comic issue number
func (c *Comic) SetIssueNumber(issueNumber string) {
	c.IssueNumber = issueNumber
}

// SetImageRegex sets the imagerex used to retrieve images link
func (c *Comic) SetImageRegex(imageRegex string) {
	c.ImageRegex = imageRegex
}

// SetURLSource sets the URL Source
func (c *Comic) SetURLSource(URLSource string) {
	c.URLSource = URLSource
}

// SetSource sets the source without the http prefix
func (c *Comic) SetSource(source string) {
	c.Source = source
}

// SetLinks sets the image links retrieve for a manga
func (c *Comic) SetImageLinks(links []string) {
	c.Links = links
}

// SetInfo will sets the name, issueNumber and imageRegex
func (c *Comic) SetInfo(name, issueNumber, imageRegex string) {
	c.Name = name
	c.IssueNumber = issueNumber
	c.ImageRegex = imageRegex
}

// SplitURL return the url splitted by "/"
func (c *Comic) SplitURL() []string {
	return strings.Split(c.URLSource, "/")
}

// SetOptions set options to the current comic
func (c *Comic) SetOptions(options map[string]string) {
	c.Options = options
}

// MakeEPUB create the epub file
func (c *Comic) makeEPUB() {
	var (
		content io.Reader
		tp      string
	)

	isCoverSet := false
	imgTag := `<img src="%s" alt="Cover Image" />`

	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	e := epub.NewEpub(c.IssueNumber)
	e.SetTitle(fmt.Sprintf("%s-%s", c.Name, c.IssueNumber))

	if c.Author != "" {
		e.SetAuthor(c.Author)
	}

	tempDir, err := ioutil.TempDir("", "comics-images")
	os.Chdir(tempDir)

	// setup the progress bar
	bar := progressbar.New(len(c.Links))
	// this will show up the progress bar since the beginning
	bar.RenderBlank()

	for i, link := range c.Links {
		if link != "" {
			rsp, err := http.Get(link)
			if err == nil {
				defer rsp.Body.Close()

				switch c.Source {
				case "mangarock.com":
					// mangarock image needs to be decoded first
					img, decErr := mri.Decode(rsp.Body)
					if decErr != nil {
						log.Error("[Mangarock] Image decode failed", decErr)
					}
					imgData := new(bytes.Buffer)
					util.ConvertTo8BitPNG(img, imgData)
					content = imgData
					tp = "png"
				default:
					content = rsp.Body
					tp = util.ImageTypeFromMime(rsp.Header["Content-Type"][0])
				}

				tmpfile, err := ioutil.TempFile(tempDir, fmt.Sprintf("image.*.%s", tp))
				if err != nil {
					log.Fatal(err)
				}
				io.Copy(tmpfile, content)

				imgpath, err := e.AddImage(tmpfile.Name(), "")
				if err != nil {
					log.Error("Can't add image: ", err)
				}

				if !isCoverSet {
					isCoverSet = true
					e.SetCover(imgpath, "")
				} else {
					_, err := e.AddSection(fmt.Sprintf(imgTag, imgpath), "", "", "")
					if err != nil {
						log.Error("Can't add section: ", err)
					}
				}

				bar.Add(i)

			} else {
				log.Error("Something went wrong with url: ", link, err)
			}
		}
	}
	os.Chdir(currentDir)
	// Set progressbar to its maximum
	bar.Finish()

	dir, dirErr := filepath.Abs(fmt.Sprintf("%s/%s/%s/%s/", filepath.Dir(os.Args[0]), "comics", c.Source, c.Name))
	if dirErr != nil {
		log.Error("There was an error: ", dirErr)
	}
	// create directories
	mkdirErr := os.MkdirAll(dir, os.ModePerm)
	if mkdirErr != nil {
		log.Error("There was an error while creating the needed folders: ", mkdirErr)
	}

	err = e.Write(fmt.Sprintf("%s/%s.epub", dir, c.IssueNumber))
	if err != nil {
		log.Error("There was an error creating the epub file: ", err)
	} else {
		log.Info("EPUB correctly saved")
	}
}

// MakeComic create the pdf file
func (c *Comic) makePDF() {
	var (
		tp      string
		content io.Reader
	)
	log.Debug("Image Download Started")
	// setup the pdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	// setup the progress bar
	bar := progressbar.New(len(c.Links))
	// this will show up the progress bar since the beginning
	bar.RenderBlank()
	// for each link get the image to add to the pdf file
	for i, link := range c.Links {
		if link != "" {
			rsp, err := http.Get(link)
			if err == nil {
				defer rsp.Body.Close()
				// add a new PDF page
				pdf.AddPage()
				switch c.Source {
				case "mangarock.com":
					// mangarock image needs to be decoded first
					// then converted to a `png` since `gofpdf` does not support webp format yet
					img, decErr := mri.Decode(rsp.Body)
					if decErr != nil {
						log.Error("[Mangarock] Image decode failed", decErr)
					}
					imgData := new(bytes.Buffer)
					util.ConvertTo8BitPNG(img, imgData)
					tp = "png"
					content = imgData
				default:
					// retrieve the image format from the response header (jpeg, png...)
					tp = pdf.ImageTypeFromMime(rsp.Header["Content-Type"][0])
					content = rsp.Body
				}
				// The image is directly added to the pdf without being saved to the disk
				pdf.RegisterImageOptionsReader(link, gofpdf.ImageOptions{tp, false}, content)
				// Here we set the image position on the pdf page
				pdf.Image(link, 0, 0, 210, 0, false, tp, 0, "")
				// increase the progressbar
				bar.Add(i)
			} else {
				log.Error("Something went wrong with url: ", link, err)
				pdf.SetError(err)
			}
		}
	}
	// Set progressbar to its maximum
	bar.Finish()
	// this will create the path where the file will be added
	dir, dirErr := filepath.Abs(fmt.Sprintf("%s/%s/%s/%s/", filepath.Dir(os.Args[0]), "comics", c.Source, c.Name))

	if dirErr != nil {
		log.Error("There was an error: ", dirErr)
	}
	// create directories and pdf file
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
	if pdf.Ok() {
		log.Info("pdf file correctly saved")
	}
}

func (c *Comic) MakeComic(format string) {
	switch strings.ToLower(format) {
	case "epub":
		c.makeEPUB()
	default:
		c.makePDF()
	}
}
