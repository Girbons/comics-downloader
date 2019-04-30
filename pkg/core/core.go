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

// getConfigValues will try to get some
func (comic *Comic) readConfigValues() {
	// retrieve values from config file
	if comic.Format == "" {
		if configFormat := comic.Config.DefaultOutputFormat; configFormat != "" {
			comic.Format = configFormat
		}
	}
}

// generateFileName will return the path where the file should be saved
func (comic *Comic) generateFileName(dir string) string {
	return fmt.Sprintf("%s/%s.%s", dir, comic.IssueNumber, comic.Format)
}

// RetrieveImageFromResponse will return the image byte and its type
func (comic *Comic) retrieveImageFromResponse(response *http.Response) (io.Reader, string, error) {
	var (
		content io.Reader
		tp      string
		err     error
	)

	switch comic.Source {
	case "mangarock.com":
		// mangarock image needs to be decoded first
		img, decErr := mri.Decode(response.Body)
		if decErr != nil {
			return content, tp, decErr
		}

		imgData := new(bytes.Buffer)
		if err := util.ConvertTo8BitPNG(img, imgData); err != nil {
			return content, tp, err
		}

		content = imgData
		tp = "png"
	default:
		content = response.Body
		tp = util.ImageType(response.Header["Content-Type"][0])
	}

	return content, tp, err
}

// makeEPUB create the epub file
func (comic *Comic) makeEPUB() error {
	var err error

	currentDir, err := util.CurrentDir()
	if err != nil {
		return err
	}
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
		return err
	}
	defer os.RemoveAll(tempDir) // clean up

	if err = os.Chdir(tempDir); err != nil {
		return err
	}
	// setup the progress bar
	bar := progressbar.NewOptions(len(comic.Links), progressbar.OptionSetRenderBlankState(true))

	for _, link := range comic.Links {
		if link != "" {
			rsp, err := http.Get(link)
			if err != nil {
				return err
			}

			defer rsp.Body.Close()
			// retrieve the image from the response
			content, tp, err := comic.retrieveImageFromResponse(rsp)
			if err != nil {
				return err
			}
			// create a tempfile to store the image
			tmpfile, err := ioutil.TempFile(tempDir, fmt.Sprintf("image.*.%s", tp))

			if err != nil {
				return err
			}
			defer os.Remove(tmpfile.Name()) // clean up

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

	if err = os.Chdir(currentDir); err != nil {
		return err
	}

	// get the PathSetup where the file should be saved
	// e.g. /www.mangarock.com/comic-name/
	dir, err := util.PathSetup(comic.Source, comic.Name)
	if err != nil {
		return err
	}

	if err = e.Write(comic.generateFileName(dir)); err != nil {
		return err
	}

	log.Info(fmt.Sprintf("%s %s", strings.ToUpper(comic.Format), DEFAULT_MESSAGE))

	return err
}

// makePDF create the pdf file
func (comic *Comic) makePDF() error {
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
				content, tp, err := comic.retrieveImageFromResponse(rsp)
				if err != nil {
					return err
				}
				// The image is directly added to the pdf without being saved to the disk
				imageOptions := gofpdf.ImageOptions{ImageType: tp, ReadDpi: false, AllowNegativePosition: true}
				pdf.RegisterImageOptionsReader(link, imageOptions, content)
				// set the image position on the pdf page
				pdf.Image(link, 0, 0, 210, 0, false, tp, 0, "")
				// increase the progressbar
			} else {
				return err
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
		return err
	}

	// Save the pdf file
	if err = pdf.OutputFileAndClose(comic.generateFileName(dir)); err != nil {
		return err
	}

	log.Info(fmt.Sprintf("%s %s", strings.ToUpper(comic.Format), DEFAULT_MESSAGE))
	return err
}

// makeCBRZ will create the CBR/CBZ
func (comic *Comic) makeCBRZ() error {
	var filesToAdd []string
	var err error

	currentDir, err := util.CurrentDir()
	if err != nil {
		return err
	}

	// setup a new Epub instance
	archive := archiver.NewZip()
	// in order to create the archive we'll need to download all the images
	tempDir, err := ioutil.TempDir("", "comics-images")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir) // clean up

	if err = os.Chdir(tempDir); err != nil {
		return err
	}
	// setup the progress bar
	bar := progressbar.NewOptions(len(comic.Links), progressbar.OptionSetRenderBlankState(true))

	for i, link := range comic.Links {
		if link != "" {
			rsp, err := http.Get(link)
			if err != nil {
				return err
			}

			defer rsp.Body.Close()
			// retrieve the image from the response
			content, tp, err := comic.retrieveImageFromResponse(rsp)
			if err != nil {
				return err
			}
			// create a tempfile to store the image
			tmpfile, err := ioutil.TempFile(tempDir, fmt.Sprintf("%d-image.*.%s", i, tp))
			defer os.Remove(tmpfile.Name()) // clean up

			if err != nil {
				return err
			}

			if _, err = io.Copy(tmpfile, content); err != nil {
				return err
			}

			filesToAdd = append(filesToAdd, tmpfile.Name())
		}
		if barErr := bar.Add(1); barErr != nil {
			log.Error(barErr)
		}
	}

	if err = os.Chdir(currentDir); err != nil {
		return err
	}
	// e.g. /www.mangarock.com/comic-name/
	dir, err := util.PathSetup(comic.Source, comic.Name)
	if err != nil {
		return err
	}
	// the archive must be created as .zip
	// then we can change the extension to .cbr or .cbz
	zipArchiveName := fmt.Sprintf("%s/%s.zip", dir, comic.IssueNumber)
	newName := fmt.Sprintf("%s/%s.%s", dir, comic.IssueNumber, comic.Format)

	if err = archive.Archive(filesToAdd, zipArchiveName); err != nil {
		return err
	}

	if err = os.Rename(zipArchiveName, newName); err != nil {
		return err
	}

	log.Info(fmt.Sprintf("%s %s", strings.ToUpper(comic.Format), DEFAULT_MESSAGE))
	return err
}

// MakeComic will create the file based on the output format selected.
func (comic *Comic) MakeComic() error {
	var err error

	if comic.Config != nil {
		comic.readConfigValues()
	}

	switch comic.Format {
	case EPUB:
		err = comic.makeEPUB()
	case CBR, CBZ:
		err = comic.makeCBRZ()
	default:
		err = comic.makePDF()
	}

	return err
}
