package core

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/BakeRolls/mri"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/util"
	epub "github.com/bmaupin/go-epub"
	"github.com/jung-kurt/gofpdf"
	"github.com/mholt/archiver"
	"github.com/schollz/progressbar"
	log "github.com/sirupsen/logrus"
)

// DEFAULT_MESSAGE for correctly saved file
const DEFAULT_MESSAGE = "file correctly saved"

// manga output format supported
const (
	CBR  = "cbr"
	CBZ  = "cbz"
	EPUB = "epub"
	PDF  = "pdf"
)

// Comic struct contains all the informations about a comic
type Comic struct {
	Author      string
	Name        string
	IssueNumber string
	Source      string
	URLSource   string
	Links       []string
	Options     map[string]string
	Format      string

	Config *config.ComicConfig
}

// SetAuthor sets the comic author
func (comic *Comic) SetAuthor(author string) {
	comic.Author = author
}

// SetName sets the comic name
func (comic *Comic) SetName(name string) {
	comic.Name = name
}

// SetIssueNumber sets the comic issue number
func (comic *Comic) SetIssueNumber(issueNumber string) {
	comic.IssueNumber = issueNumber
}

// SetURLSource sets the URL Source
func (comic *Comic) SetURLSource(URLSource string) {
	comic.URLSource = URLSource
}

// SetSource sets the source without the http prefix
func (comic *Comic) SetSource(source string) {
	comic.Source = source
}

// SetLinks sets the image links retrieved for a manga
func (comic *Comic) SetImageLinks(links []string) {
	comic.Links = links
}

// SetConfiguration
func (comic *Comic) SetConfig(config *config.ComicConfig) {
	comic.Config = config
}

// getConfigValues will try to get some
func (comic *Comic) getConfigValues() {
	// retrieve values from config file if exist
	if comic.Format == "" {
		if configFormat := comic.Config.DefaultOutputFormat; configFormat != "" {
			comic.Format = configFormat
		}
	}
}

// SetFormat sets the comic output format
func (comic *Comic) SetFormat(format string) {
	switch strings.ToLower(format) {
	case EPUB:
		comic.Format = EPUB
	case CBR:
		comic.Format = CBR
	case CBZ:
		comic.Format = CBZ
	default:
		comic.Format = PDF
	}
}

// SetInfo will sets the name, issueNumber
func (comic *Comic) SetInfo(name, issueNumber string) {
	comic.Name = name
	comic.IssueNumber = issueNumber
}

// SplitURL return the url splitted by "/"
func (comic *Comic) SplitURL() []string {
	return strings.Split(comic.URLSource, "/")
}

// SetOptions set options to the current comic
func (comic *Comic) SetOptions(options map[string]string) {
	comic.Options = options
}

// generateFileName will return the path where the file should be saved
func (comic *Comic) generateFileName(dir string) string {
	return fmt.Sprintf("%s/%s.%s", dir, comic.IssueNumber, comic.Format)
}

// RetrieveImageFromResponse will return the image byte and its type
func (comic *Comic) retrieveImageFromResponse(response *http.Response) (io.Reader, string) {
	var (
		content io.Reader
		tp      string
	)

	switch comic.Source {
	case "mangarock.com":
		// mangarock image needs to be decoded first
		img, decErr := mri.Decode(response.Body)
		if decErr != nil {
			log.Error(decErr)
		}

		imgData := new(bytes.Buffer)
		if err := util.ConvertTo8BitPNG(img, imgData); err != nil {
			log.Error(err)
		}

		content = imgData
		tp = "png"
	default:
		content = response.Body
		tp = util.ImageType(response.Header["Content-Type"][0])
	}

	return content, tp

}

// makeEPUB create the epub file
func (comic *Comic) makeEPUB() {
	var err error
	// used to check if the epub cover already exists
	isCoverSet := false
	// used to add the image in the epub section
	imgTag := `<img src="%s" alt="Cover Image" />`
	// setup a new Epub instance
	e := epub.NewEpub(comic.IssueNumber)
	// set Epub title
	e.SetTitle(fmt.Sprintf("%s-%s", comic.Name, comic.IssueNumber))
	// check if the author exists for this comic
	if comic.Author != "" {
		e.SetAuthor(comic.Author)
	}
	// in order to create an epub we'll need to download all the images so we create a tempdir for that
	tempDir, err := ioutil.TempDir("", "comics-images")
	if err != nil {
		log.WithFields(log.Fields{
			"source": comic.Source,
		}).Fatal(err)
	}
	defer os.RemoveAll(tempDir) // clean up

	if err = os.Chdir(tempDir); err != nil {
		log.Fatal(err)
	}
	// setup the progress bar
	bar := progressbar.NewOptions(len(comic.Links), progressbar.OptionSetRenderBlankState(true))

	for _, link := range comic.Links {
		if link != "" {
			rsp, err := http.Get(link)

			if err != nil {
				log.Error(err)
			}

			defer rsp.Body.Close()
			// retrieve the image from the response
			content, tp := comic.retrieveImageFromResponse(rsp)
			// create a tempfile to store the image
			tmpfile, err := ioutil.TempFile(tempDir, fmt.Sprintf("image.*.%s", tp))
			defer os.Remove(tmpfile.Name()) // clean up

			if err != nil {
				log.Fatal(err)
			}

			if _, err = io.Copy(tmpfile, content); err != nil {
				log.Error(err)
			}
			// add the image to the epub will return a path
			imgpath, err := e.AddImage(tmpfile.Name(), "")
			if err != nil {
				log.Error(err)
			}
			// if the cover is not set we'll use the first image
			// otherwise the image will be added as a section
			if !isCoverSet {
				isCoverSet = true
				e.SetCover(imgpath, "")
			} else {
				_, err = e.AddSection(fmt.Sprintf(imgTag, imgpath), "", "", "")
				if err != nil {
					log.Error(err)
				}
			}
		}
		if barErr := bar.Add(1); barErr != nil {
			log.Error(barErr)
		}
	}
	if err = os.Chdir(tempDir); err != nil {
		log.Fatal(err)
	}
	// get the PathSetup where the file should be saved
	// e.g. /www.mangarock.com/comic-name/
	dir, err := util.PathSetup(comic.Source, comic.Name)
	if err != nil {
		log.Fatal(err)
	}

	if err = e.Write(comic.generateFileName(dir)); err != nil {
		log.Fatal(err)
	} else {
		log.Info(fmt.Sprintf("%s %s", strings.ToUpper(comic.Format), DEFAULT_MESSAGE))
	}
}

// makePDF create the pdf file
func (comic *Comic) makePDF() {
	var err error
	// setup the pdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	// setup the progress bar
	bar := progressbar.NewOptions(len(comic.Links), progressbar.OptionSetRenderBlankState(true))
	// for each link get the image to add to the pdf file
	for _, link := range comic.Links {
		if link != "" {
			rsp, err := http.Get(link)
			if err == nil {
				defer rsp.Body.Close()
				// add a new PDF page
				pdf.AddPage()
				content, tp := comic.retrieveImageFromResponse(rsp)
				// The image is directly added to the pdf without being saved to the disk
				imageOptions := gofpdf.ImageOptions{ImageType: tp, ReadDpi: false, AllowNegativePosition: true}
				pdf.RegisterImageOptionsReader(link, imageOptions, content)
				// set the image position on the pdf page
				pdf.Image(link, 0, 0, 210, 0, false, tp, 0, "")
				// increase the progressbar
			} else {
				log.Error(err)
				pdf.SetError(err)
			}
		}
		if barErr := bar.Add(1); barErr != nil {
			log.Error(barErr)
		}
	}
	// get the PathSetup where the file should be saved
	// e.g. /www.mangarock.com/comic-name/
	dir, err := util.PathSetup(comic.Source, comic.Name)
	if err != nil {
		log.Fatal(err)
	}

	// Save the pdf file
	if err = pdf.OutputFileAndClose(comic.generateFileName(dir)); err != nil {
		log.Fatal(err)
	}

	if pdf.Ok() {
		log.Info(fmt.Sprintf("%s %s", strings.ToUpper(comic.Format), DEFAULT_MESSAGE))
	}
}

// makeCBRZ will create the CBR/CBZ
func (comic *Comic) makeCBRZ() {
	var filesToAdd []string
	var err error
	// setup a new Epub instance
	archive := archiver.NewZip()
	// in order to create the archive we'll need to download all the images
	tempDir, err := ioutil.TempDir("", "comics-images")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(tempDir) // clean up
	if err = os.Chdir(tempDir); err != nil {
		log.Fatal(err)
	}
	// setup the progress bar
	bar := progressbar.NewOptions(len(comic.Links), progressbar.OptionSetRenderBlankState(true))

	for i, link := range comic.Links {
		if link != "" {
			rsp, err := http.Get(link)
			if err == nil {
				defer rsp.Body.Close()
				// retrieve the image from the response
				content, tp := comic.retrieveImageFromResponse(rsp)
				// create a tempfile to store the image
				tmpfile, err := ioutil.TempFile(tempDir, fmt.Sprintf("%d-image.*.%s", i, tp))
				defer os.Remove(tmpfile.Name()) // clean up

				if err != nil {
					log.Fatal(err)
				}

				if _, err = io.Copy(tmpfile, content); err != nil {
					log.Error(err)
				}

				filesToAdd = append(filesToAdd, tmpfile.Name())

			} else {
				log.Error(err)
			}
		}
		if barErr := bar.Add(1); barErr != nil {
			log.Error(barErr)
		}
	}

	if err = os.Chdir(tempDir); err != nil {
		log.Fatal(err)
	}
	// e.g. /www.mangarock.com/comic-name/
	dir, err := util.PathSetup(comic.Source, comic.Name)
	if err != nil {
		log.Fatal(err)
	}
	// the archive must be created as .zip
	// then we can change the extension to .cbr or .cbz
	zipArchiveName := fmt.Sprintf("%s/%s.zip", dir, comic.IssueNumber)
	newName := fmt.Sprintf("%s/%s.%s", dir, comic.IssueNumber, comic.Format)
	if err = archive.Archive(filesToAdd, zipArchiveName); err != nil {
		log.Fatal(err)
	} else {
		if err := os.Rename(zipArchiveName, newName); err != nil {
			log.Fatal(err)
		}
		log.Info(fmt.Sprintf("%s %s", strings.ToUpper(comic.Format), DEFAULT_MESSAGE))
	}
}

// MakeComic will create the file based on the output format selected.
func (comic *Comic) MakeComic() {
	if comic.Config != nil {
		comic.getConfigValues()
	}

	switch comic.Format {
	case EPUB:
		comic.makeEPUB()
	case CBR, CBZ:
		comic.makeCBRZ()
	default:
		comic.makePDF()
	}
}
