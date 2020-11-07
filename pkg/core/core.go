package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/util"
	epub "github.com/bmaupin/go-epub"
	"github.com/jung-kurt/gofpdf"
	"github.com/mholt/archiver"
	"github.com/schollz/progressbar/v2"
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
	Author       string
	Name         string
	IssueNumber  string
	Source       string
	URLSource    string
	Links        []string
	Format       string
	ImagesFormat string
}

// makeEPUB create the epub file
func (comic *Comic) makeEPUB(options *config.Options) error {
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

	imagesPath, err := comic.DownloadImages(options)
	if err != nil {
		return err
	}
	defer os.RemoveAll(imagesPath)

	files, err := ioutil.ReadDir(imagesPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		// add the image to the epub will return a path
		imgpath, err := e.AddImage(fmt.Sprintf("%s/%s", imagesPath, file.Name()), "")
		if err != nil {
			options.Logger.Error(err.Error())
		}
		// if the cover is not set we'll use the first image
		// otherwise the image will be added as a section
		if !isCoverSet {
			isCoverSet = true
			e.SetCover(imgpath, "")
		} else {
			_, err = e.AddSection(fmt.Sprintf(imgTag, imgpath), "", "", "")
			if err != nil {
				options.Logger.Error(err.Error())
			}
		}
	}

	if err = os.Chdir(currentDir); err != nil {
		return err
	}

	// get the PathSetup where the file should be saved
	// e.g. /www.mangarock.com/comic-name/
	dir, err := util.PathSetup(options.OutputFolder, comic.Source, comic.Name)
	if err != nil {
		return err
	}

	if err = e.Write(util.GenerateFileName(dir, comic.Name, comic.IssueNumber, comic.Format)); err != nil {
		return err
	}

	options.Logger.Info(fmt.Sprintf("%s %s", strings.ToUpper(comic.Format), DEFAULT_MESSAGE))
	return err
}

// makePDF create the pdf file
func (comic *Comic) makePDF(options *config.Options) error {
	var err error
	// setup the pdf
	pdf := gofpdf.New("P", "mm", "A4", "")

	imagesPath, err := comic.DownloadImages(options)
	if err != nil {
		return err
	}

	defer os.RemoveAll(imagesPath)

	files, err := ioutil.ReadDir(imagesPath)
	if err != nil {
		return err
	}

	// for each link get the image to add to the pdf file
	for _, file := range files {
		// add a new PDF page
		pdf.AddPage()
		imageOptions := gofpdf.ImageOptions{ImageType: util.ImageType(comic.ImagesFormat), ReadDpi: true, AllowNegativePosition: false}
		fileName := fmt.Sprintf("%s/%s", imagesPath, file.Name())
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			return err
		}
		content := bytes.NewReader(data)
		pdf.RegisterImageOptionsReader(file.Name(), imageOptions, content)
		// set the image position on the pdf page
		pdf.Image(file.Name(), 0, 0, 210, 297, false, comic.ImagesFormat, 0, "")
	}
	// get the PathSetup where the file should be saved
	// e.g. /www.mangarock.com/comic-name/
	dir, err := util.PathSetup(options.OutputFolder, comic.Source, comic.Name)
	if err != nil {
		return err
	}

	// Save the pdf file
	if err = pdf.OutputFileAndClose(util.GenerateFileName(dir, comic.Name, comic.IssueNumber, comic.Format)); err != nil {
		return err
	}

	options.Logger.Info(fmt.Sprintf("%s %s", strings.ToUpper(comic.Format), DEFAULT_MESSAGE))
	return err
}

// makeCBRZ will create the CBR/CBZ
func (comic *Comic) makeCBRZ(options *config.Options) error {
	var filesToAdd []string
	var err error

	// setup a new Epub instance
	archive := archiver.NewZip()

	imagesPath, err := comic.DownloadImages(options)
	if err != nil {
		return err
	}
	defer os.RemoveAll(imagesPath)

	files, err := ioutil.ReadDir(imagesPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		filesToAdd = append(filesToAdd, fmt.Sprintf("%s/%s", imagesPath, file.Name()))
	}

	// e.g. /www.mangarock.com/comic-name/
	dir, err := util.PathSetup(options.OutputFolder, comic.Source, comic.Name)
	if err != nil {
		return err
	}
	// the archive must be created as .zip
	// then we can change the extension to .cbr or .cbz
	zipArchiveName := fmt.Sprintf("%s/%s.zip", dir, comic.IssueNumber)
	newName := util.GenerateFileName(dir, comic.Name, comic.IssueNumber, comic.Format)

	if err = archive.Archive(filesToAdd, zipArchiveName); err != nil {
		return err
	}

	if err = os.Rename(zipArchiveName, newName); err != nil {
		return err
	}

	options.Logger.Info(fmt.Sprintf("%s %s", strings.ToUpper(comic.Format), DEFAULT_MESSAGE))
	return err
}

// DownloadImages will download the comic/manga images
func (comic *Comic) DownloadImages(options *config.Options) (string, error) {
	var dir string
	var err error

	dir, err = util.ImagesPathSetup(options.OutputFolder, comic.Source, comic.Name, comic.IssueNumber)
	if err != nil {
		return dir, err
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return dir, err
	}

	if !util.DirectoryOrFileDoesNotExist(dir) && len(files) == len(comic.Links) {
		return dir, err
	}

	format := util.ImageType(comic.ImagesFormat)

	currentDir, err := util.CurrentDir()
	if err != nil {
		return dir, err
	}

	// setup the progress bar
	bar := progressbar.NewOptions(len(comic.Links), progressbar.OptionSetRenderBlankState(true))

	err = os.Chdir(dir)
	if err != nil {
		return dir, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    11,
			IdleConnTimeout: 30 * time.Second,
		},
	}

	for i, link := range comic.Links {
		if link != "" {
			rsp, err := client.Get(link)
			if err != nil {
				return dir, err
			}
			defer rsp.Body.Close()

			imgFile, err := os.Create(fmt.Sprintf("%04d-image.%s", i, format))
			if err != nil {
				return dir, err
			}
			defer imgFile.Close()

			err = util.SaveImage(imgFile, rsp.Body, format)
			if err != nil {
				return dir, err
			}
		}

		if barErr := bar.Add(1); barErr != nil {
			options.Logger.Error(barErr.Error())
		}
	}

	err = os.Chdir(currentDir)
	if err != nil {
		return dir, err
	}

	return dir, err
}

// MakeComic will create the file based on the output format selected.
func (comic *Comic) MakeComic(options *config.Options) error {
	var err error

	switch comic.Format {
	case EPUB:
		err = comic.makeEPUB(options)
	case CBR, CBZ:
		err = comic.makeCBRZ(options)
	default:
		err = comic.makePDF(options)
	}

	return err
}
