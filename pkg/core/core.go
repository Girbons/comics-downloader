package core

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/util"
	epub "github.com/bmaupin/go-epub"
	"github.com/jung-kurt/gofpdf"
	"github.com/mholt/archives"
	"github.com/schollz/progressbar/v2"
)

// DefaultMessage for correctly saved file
const DefaultMessage = "file correctly saved"

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

	files, err := os.ReadDir(imagesPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		// add the image to the epub will return a path
		imgpath, err := e.AddImage(fmt.Sprintf("%s/%s", imagesPath, file.Name()), "")
		if err != nil {
			options.Logger.Error(err.Error())
		}
		// if the cover is not set use the first image
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
	dir, err := util.PathSetup(options.CreateDefaultPath, options.OutputFolder, comic.Source, comic.Name)
	if err != nil {
		return err
	}

	if err = e.Write(util.GetPathToFile(dir, comic.Name, comic.IssueNumber, comic.Format, options.IssueNumberNameOnly)); err != nil {
		return err
	}

	options.Logger.Info(fmt.Sprintf("%s %s", strings.ToUpper(comic.Format), DefaultMessage))
	return err
}

// makePDF create the pdf file
func (comic *Comic) makePDF(options *config.Options) error {
	var err error
	var mmWd, mmHt float64
	const px2mm = 0.2645833333

	pdf := gofpdf.New("P", "mm", "A4", "")

	imagesPath, err := comic.DownloadImages(options)
	if err != nil {
		return err
	}

	defer os.RemoveAll(imagesPath)

	files, err := os.ReadDir(imagesPath)
	if err != nil {
		return err
	}

	imageOptions := gofpdf.ImageOptions{ImageType: util.ImageType(comic.ImagesFormat), ReadDpi: true, AllowNegativePosition: false}
	for _, file := range files {
		mmWd = 210.0
		mmHt = 297.0
		fileName := fmt.Sprintf("%s/%s", imagesPath, file.Name())

		if !options.ForceAspect {
			img, err := os.Open(fileName)
			if err != nil {
				options.Logger.Error(err.Error())
			}

			defer img.Close()

			im, _, err := image.DecodeConfig(img)
			if err != nil {
				options.Logger.Error(err.Error())
			} else {
				mmWd = px2mm * float64(im.Width)
				mmHt = px2mm * float64(im.Height)
			}
		}
		pdf.AddPageFormat("P", gofpdf.SizeType{Wd: mmWd, Ht: mmHt})

		data, err := os.ReadFile(fileName)
		if err != nil {
			return err
		}
		content := bytes.NewReader(data)
		pdf.RegisterImageOptionsReader(file.Name(), imageOptions, content)
		pdf.ImageOptions(file.Name(), 0, 0, mmWd, mmHt, false, imageOptions, 0, "")
	}

	dir, err := util.PathSetup(options.CreateDefaultPath, options.OutputFolder, comic.Source, comic.Name)
	if err != nil {
		return err
	}

	// Save the pdf file
	filePath := util.GetPathToFile(dir, comic.Name, comic.IssueNumber, comic.Format, options.IssueNumberNameOnly)
	if err = pdf.OutputFileAndClose(filePath); err != nil {
		return err
	}

	options.Logger.Info(fmt.Sprintf("%s %s", strings.ToUpper(comic.Format), DefaultMessage))
	return err
}

// makeCBRZ will create the CBR/CBZ
func (comic *Comic) makeCBRZ(options *config.Options) error {
	var filesToAdd []string
	var err error

	imagesPath, err := comic.DownloadImages(options)
	if err != nil {
		return err
	}
	defer os.RemoveAll(imagesPath)

	files, err := os.ReadDir(imagesPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		filesToAdd = append(filesToAdd, fmt.Sprintf("%s/%s", imagesPath, file.Name()))
	}

	// e.g. /www.mangarock.com/comic-name/
	dir, err := util.PathSetup(options.CreateDefaultPath, options.OutputFolder, comic.Source, comic.Name)
	if err != nil {
		return err
	}

	// the archive must be created as `.zip` then change the extension to `.cbr` or `.cbz`.
	zipArchiveName := fmt.Sprintf("%s/%s.zip", dir, comic.IssueNumber)
	newName := util.GetPathToFile(dir, comic.Name, comic.IssueNumber, comic.Format, options.IssueNumberNameOnly)

	// Create output file
	out, err := os.Create(zipArchiveName)
	if err != nil {
		return err
	}
	defer out.Close()
	// Sort files to ensure consistent ordering
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	// Map files on disk to their paths in the archive
	fileMap := make(map[string]string)
	for _, filePath := range filesToAdd {
		// Use just the filename (no path) in the archive
		fileName := path.Base(filePath)
		fileMap[filePath] = fileName
	}

	// Get files from disk using the archives helper
	archiveFiles, err := archives.FilesFromDisk(context.Background(), nil, fileMap)
	if err != nil {
		return err
	}

	// Create ZIP format (no compression needed for CBZ)
	format := archives.Zip{}

	// Create the archive
	err = format.Archive(context.Background(), out, archiveFiles)
	if err != nil {
		return err
	}

	if err = os.Rename(zipArchiveName, newName); err != nil {
		return err
	}

	options.Logger.Info(fmt.Sprintf("%s %s", strings.ToUpper(comic.Format), DefaultMessage))
	return nil
}

// DownloadImages will download the comic/manga images
func (comic *Comic) DownloadImages(options *config.Options) (string, error) {
	if len(comic.Links) == 0 {
		return "", fmt.Errorf("Download failed, no links found for: %s", comic.URLSource)
	}

	var dir string
	var err error

	dir, err = util.ImagesPathSetup(options.CreateDefaultPath, options.OutputFolder, comic.Source, comic.Name, options.IssueFolderName, comic.IssueNumber)
	if err != nil {
		return dir, err
	}

	files, err := os.ReadDir(dir)
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

	g := new(errgroup.Group)

	maxWorkers := int64(runtime.NumCPU())
	sem := semaphore.NewWeighted(maxWorkers)
	ctx := context.Background()

	defer sem.Acquire(ctx, maxWorkers)
	for i, link := range comic.Links {
		link, i := link, i
		sem.Acquire(ctx, 1)

		if link == "" {
			continue
		}

		g.Go(func() error {
			defer sem.Release(1)
			rsp, err := options.Client.Get(link, comic.Source)
			if err != nil {
				return err
			}
			defer rsp.Body.Close()

			imgName := fmt.Sprintf("%04d-image.%s", i, format)
			imgFile, err := os.Create(imgName)
			if err != nil {
				return err
			}
			defer imgFile.Close()

			isWebp := strings.HasSuffix(link, ".webp")
			err = util.SaveImage(imgFile, rsp.Body, format, isWebp)
			if err != nil {
				msgError := fmt.Sprintf("There was an error while downloading image number: %d - comic issue: %s", i, comic.IssueNumber)
				options.Logger.Error(msgError)
				os.Remove(imgName)
			}

			if barErr := bar.Add(1); barErr != nil {
				options.Logger.Error(barErr.Error())
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return dir, err
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
