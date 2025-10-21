package core

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"

	"github.com/Girbons/comics-downloader/pkg/config"
	httpclient "github.com/Girbons/comics-downloader/pkg/http"
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

// DownloadResult captures the outcome of downloading a comic's images.
type DownloadResult struct {
	Dir       string
	FilePaths []string
}

func ensureClient(options *config.Options) *httpclient.ComicClient {
	if options.Client == nil {
		options.Client = httpclient.NewComicClient()
	}
	return options.Client
}

// makeEPUB creates the epub file.
func (comic *Comic) makeEPUB(options *config.Options, images *DownloadResult) error {
	isCoverSet := false
	imgTag := `<img src="%s" alt="Cover Image" />`
	e := epub.NewEpub(comic.IssueNumber)
	e.SetTitle(fmt.Sprintf("%s-%s", comic.Name, comic.IssueNumber))

	if comic.Author != "" {
		e.SetAuthor(comic.Author)
	}

	for _, file := range images.FilePaths {
		imgpath, err := e.AddImage(file, "")
		if err != nil && options.Logger != nil {
			options.Logger.Error(err.Error())
			continue
		}
		if !isCoverSet {
			isCoverSet = true
			e.SetCover(imgpath, "")
			continue
		}
		if _, err := e.AddSection(fmt.Sprintf(imgTag, imgpath), "", "", ""); err != nil && options.Logger != nil {
			options.Logger.Error(err.Error())
		}
	}

	dir, err := util.PathSetup(options.CreateDefaultPath, options.OutputFolder, comic.Source, comic.Name)
	if err != nil {
		return err
	}

	if err = e.Write(util.GetPathToFile(dir, comic.Name, comic.IssueNumber, comic.Format, options.IssueNumberNameOnly)); err != nil {
		return err
	}

	if options.Logger != nil {
		options.Logger.Infof("%s %s", strings.ToUpper(comic.Format), DefaultMessage)
	}
	return nil
}

// makePDF create the pdf file.
func (comic *Comic) makePDF(options *config.Options, images *DownloadResult) error {
	var mmWd, mmHt float64
	const px2mm = 0.2645833333

	pdf := gofpdf.New("P", "mm", "A4", "")

	imageOptions := gofpdf.ImageOptions{ImageType: util.ImageType(comic.ImagesFormat), ReadDpi: true, AllowNegativePosition: false}
	for _, fileName := range images.FilePaths {
		mmWd = 210.0
		mmHt = 297.0

		if !options.ForceAspect {
			img, err := os.Open(fileName)
			if err != nil {
				if options.Logger != nil {
					options.Logger.Error(err.Error())
				}
			} else {
				im, _, err := image.DecodeConfig(img)
				img.Close()
				if err != nil {
					if options.Logger != nil {
						options.Logger.Error(err.Error())
					}
				} else {
					mmWd = px2mm * float64(im.Width)
					mmHt = px2mm * float64(im.Height)
				}
			}
		}
		pdf.AddPageFormat("P", gofpdf.SizeType{Wd: mmWd, Ht: mmHt})

		data, err := os.ReadFile(fileName)
		if err != nil {
			return err
		}
		content := bytes.NewReader(data)
		pdf.RegisterImageOptionsReader(path.Base(fileName), imageOptions, content)
		pdf.ImageOptions(path.Base(fileName), 0, 0, mmWd, mmHt, false, imageOptions, 0, "")
	}

	dir, err := util.PathSetup(options.CreateDefaultPath, options.OutputFolder, comic.Source, comic.Name)
	if err != nil {
		return err
	}

	filePath := util.GetPathToFile(dir, comic.Name, comic.IssueNumber, comic.Format, options.IssueNumberNameOnly)
	if err = pdf.OutputFileAndClose(filePath); err != nil {
		return err
	}

	if options.Logger != nil {
		options.Logger.Infof("%s %s", strings.ToUpper(comic.Format), DefaultMessage)
	}
	return nil
}

// makeCBRZ will create the CBR/CBZ.
func (comic *Comic) makeCBRZ(options *config.Options, images *DownloadResult) error {
	dir, err := util.PathSetup(options.CreateDefaultPath, options.OutputFolder, comic.Source, comic.Name)
	if err != nil {
		return err
	}

	zipArchiveName := filepath.Join(dir, fmt.Sprintf("%s.zip", comic.IssueNumber))
	newName := util.GetPathToFile(dir, comic.Name, comic.IssueNumber, comic.Format, options.IssueNumberNameOnly)

	out, err := os.Create(zipArchiveName)
	if err != nil {
		return err
	}
	defer out.Close()

	fileMap := make(map[string]string)
	for _, filePath := range images.FilePaths {
		fileMap[filePath] = path.Base(filePath)
	}

	archiveFiles, err := archives.FilesFromDisk(context.Background(), nil, fileMap)
	if err != nil {
		return err
	}

	format := archives.Zip{}
	if err = format.Archive(context.Background(), out, archiveFiles); err != nil {
		return err
	}

	if err = os.Rename(zipArchiveName, newName); err != nil {
		return err
	}

	if options.Logger != nil {
		options.Logger.Infof("%s %s", strings.ToUpper(comic.Format), DefaultMessage)
	}
	return nil
}

// DownloadImages will download the comic/manga images.
func (comic *Comic) DownloadImages(options *config.Options) (*DownloadResult, error) {
	if len(comic.Links) == 0 {
		return nil, fmt.Errorf("download failed, no links found for: %s", comic.URLSource)
	}

	client := ensureClient(options)

	dir, err := util.ImagesPathSetup(options.CreateDefaultPath, options.OutputFolder, comic.Source, comic.Name, options.IssueFolderName, comic.IssueNumber)
	if err != nil {
		return nil, err
	}

	existing, err := readExistingImages(dir)
	if err == nil && len(existing) == len(comic.Links) && len(existing) > 0 {
		return &DownloadResult{Dir: dir, FilePaths: existing}, nil
	}

	if err := os.RemoveAll(dir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	progress := progressbar.NewOptions(len(comic.Links), progressbar.OptionSetRenderBlankState(true))
	format := util.ImageType(comic.ImagesFormat)

	requestDelay := options.RequestDelay
	requestJitter := options.RequestDelayJitter
	if requestDelay < 0 {
		requestDelay = 0
	}
	if requestJitter < 0 {
		requestJitter = 0
	}
	if requestDelay == 0 && requestJitter == 0 {
		requestDelay = config.DefaultRequestDelay
		requestJitter = config.DefaultRequestDelayJitter
	}

	type downloadJob struct {
		index int
		link  string
	}

	jobs := make([]downloadJob, 0, len(comic.Links))
	for idx, link := range comic.Links {
		if strings.TrimSpace(link) == "" {
			continue
		}
		jobs = append(jobs, downloadJob{index: idx, link: link})
	}

	results := make([]string, len(comic.Links))
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)
	sem := semaphore.NewWeighted(int64(runtime.NumCPU()))
	var mu sync.Mutex
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var rngMu sync.Mutex

	for _, job := range jobs {
		job := job
		if err := sem.Acquire(ctx, 1); err != nil {
			return nil, err
		}
		group.Go(func() error {
			defer sem.Release(1)

			reqCtx, cancelReq := context.WithTimeout(ctx, 30*time.Second)
			defer cancelReq()

			sleepDuration := requestDelay
			if requestJitter > 0 {
				rngMu.Lock()
				extra := time.Duration(rng.Int63n(int64(requestJitter)))
				rngMu.Unlock()
				sleepDuration += extra
			}
			if sleepDuration > 0 {
				time.Sleep(sleepDuration)
			}

			request, err := client.PrepareRequest(job.link, comic.Source)
			if err != nil {
				return err
			}
			request = request.WithContext(reqCtx)

			response, err := client.Do(request)
			if err != nil {
				return err
			}
			defer response.Body.Close()

			fileName := fmt.Sprintf("%04d-image.%s", job.index, format)
			targetPath := filepath.Join(dir, fileName)
			imgFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}

			isWebp := strings.HasSuffix(strings.ToLower(job.link), ".webp")
			if err := util.SaveImage(imgFile, response.Body, format, isWebp); err != nil {
				if options.Logger != nil {
					options.Logger.Errorf("There was an error while downloading image number: %d - comic issue: %s (%v)", job.index, comic.IssueNumber, err)
				}
				imgFile.Close()
				_ = os.Remove(targetPath)
			} else {
				imgFile.Close()
				mu.Lock()
				results[job.index] = targetPath
				mu.Unlock()
			}

			if progressErr := progress.Add(1); progressErr != nil && options.Logger != nil {
				options.Logger.Error(progressErr.Error())
			}
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}

	paths := filterEmpty(results)
	sort.Strings(paths)
	return &DownloadResult{Dir: dir, FilePaths: paths}, nil
}

func readExistingImages(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files = append(files, filepath.Join(dir, entry.Name()))
	}
	sort.Strings(files)
	return files, nil
}

func filterEmpty(items []string) []string {
	var filtered []string
	for _, item := range items {
		if item != "" {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// MakeComic will create the file based on the output format selected.
func (comic *Comic) MakeComic(options *config.Options) error {
	result, err := comic.DownloadImages(options)
	if err != nil {
		return err
	}
	defer os.RemoveAll(result.Dir)

	switch comic.Format {
	case EPUB:
		return comic.makeEPUB(options, result)
	case CBR, CBZ:
		return comic.makeCBRZ(options, result)
	default:
		return comic.makePDF(options, result)
	}
}
