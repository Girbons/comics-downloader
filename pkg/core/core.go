package core

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"io"
	"net/http"
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

// DownloadResult captures the outcome of downloading a comic's images.
type DownloadResult struct {
	Dir       string
	FilePaths []string // Absolute paths to the downloaded image files
}

func ensureClient(options *config.Options) *httpclient.ComicClient {
	if options.Client == nil {
		options.Client = httpclient.NewComicClient()
	}
	return options.Client
}

// makeEPUB creates the epub file.
func (comic *ComicIssue) makeEPUB(options *config.Options, images *DownloadResult) error {
	isCoverSet := false
	imgTag := `<img src="%s" alt="Comic Page" />`

	localizedTitle := comic.getLocalizedTitle(options)
	e := epub.NewEpub(localizedTitle)

	if len(comic.SeriesMetadata.Creators) > 0 {
		setAuthor := false
		writerRole := string(CreatorRoleWriter)
		for _, creator := range comic.SeriesMetadata.Creators {
			if strings.EqualFold(string(creator.Role), writerRole) {
				e.SetAuthor(creator.Name)
				setAuthor = true
				break
			}
		}

		if !setAuthor {
			// if no creator with role "Writer" is found, set the author to the first creator in the list
			author := comic.SeriesMetadata.Creators[0].Name
			e.SetAuthor(author)
		}
	}

	if comic.LanguageISO != nil {
		e.SetLang(*comic.LanguageISO)
	}

	comicDescription := comic.getLocalizedDescription(options)
	if comicDescription != "" {
		e.SetDescription(comicDescription)
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

	outputFilePath, err := comic.GetOutputFilePath(options)
	if err != nil {
		return err
	}

	if err = e.Write(outputFilePath); err != nil {
		return err
	}

	if options.Logger != nil {
		options.Logger.Infof("%s %s", strings.ToUpper(comic.OutputFormat.String()), DefaultMessage)
	}
	return nil
}

// makePDF create the pdf file.
func (comic *ComicIssue) makePDF(options *config.Options, images *DownloadResult) error {
	var mmWd, mmHt float64
	const px2mm = 0.2645833333

	pdf := gofpdf.New("P", "mm", "A4", "")

	imageOptions := gofpdf.ImageOptions{ImageType: util.ImageType(comic.OutputImagesFormat).String(), ReadDpi: true, AllowNegativePosition: false}
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
				if closeErr := img.Close(); closeErr != nil && options.Logger != nil {
					options.Logger.Errorf("failed to close image %s: %v", fileName, closeErr)
				}
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

	filePath, err := comic.GetOutputFilePath(options)
	if err != nil {
		return err
	}

	if err = pdf.OutputFileAndClose(filePath); err != nil {
		return err
	}

	if options.Logger != nil {
		options.Logger.Infof("%s %s", strings.ToUpper(comic.OutputFormat.String()), DefaultMessage)
	}
	return nil
}

// makeCBRZ will create the CBR/CBZ.
func (comic *ComicIssue) makeCBRZ(options *config.Options, images *DownloadResult) error {
	dir, err := comic.GetOutputDir(options)
	if err != nil {
		return err
	}

	comicinfoXMLPath, err := comic.makeComicInfoXML(options, images)
	if err != nil {
		return err
	}

	newName, err := comic.GetOutputFilePath(options)
	if err != nil {
		return err
	}

	// check if file already exists to avoid creating the archive again
	// this is a final sanity check
	if _, statErr := os.Stat(newName); statErr == nil {
		if options.Logger != nil {
			options.Logger.Infof("Skipping %s because it already exists: %s", strings.ToUpper(comic.OutputFormat.String()), newName)
		}
		return nil
	} else if !os.IsNotExist(statErr) {
		return statErr
	}

	out, err := os.CreateTemp(dir, fmt.Sprintf("%s-*.zip", comic.IssueNumber))
	if err != nil {
		return err
	}
	zipArchiveName := out.Name()
	defer func() {
		if out != nil {
			if closeErr := out.Close(); closeErr != nil && options.Logger != nil {
				options.Logger.Errorf("failed to close archive %s: %v", zipArchiveName, closeErr)
			}
		}
		if zipArchiveName != "" {
			if removeErr := os.Remove(zipArchiveName); removeErr != nil && !os.IsNotExist(removeErr) && options.Logger != nil {
				options.Logger.Errorf("failed to cleanup temp archive %s: %v", zipArchiveName, removeErr)
			}
		}
	}()

	fileMap := make(map[string]string)
	for _, filePath := range images.FilePaths {
		fileMap[filePath] = path.Base(filePath)
	}
	fileMap[comicinfoXMLPath] = "ComicInfo.xml"

	archiveFiles, err := archives.FilesFromDisk(context.Background(), nil, fileMap)
	if err != nil {
		return err
	}

	format := archives.Zip{}
	if err = format.Archive(context.Background(), out, archiveFiles); err != nil {
		return err
	}

	if err = out.Close(); err != nil {
		return err
	}
	out = nil

	if err = os.Rename(zipArchiveName, newName); err != nil {
		return err
	}

	if options.Logger != nil {
		options.Logger.Infof("%s %s", strings.ToUpper(comic.OutputFormat.String()), DefaultMessage)
	}
	return nil
}

// DownloadImages will download the comic/manga images.
func (comic *ComicIssue) DownloadImages(options *config.Options) (*DownloadResult, error) {
	if len(comic.ImageLinks) == 0 {
		return nil, fmt.Errorf("download failed, no links found for: %s", comic.Source.URL)
	}

	client := ensureClient(options)

	dir, err := comic.GetImagesOutputDir(options)
	if err != nil {
		return nil, err
	}

	existing, err := readExistingImages(dir)
	if err == nil && len(existing) == len(comic.ImageLinks) && len(existing) > 0 {
		return &DownloadResult{Dir: dir, FilePaths: existing}, nil
	}

	if err := os.RemoveAll(dir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	var progress *progressbar.ProgressBar
	if !options.Debug {
		progress = progressbar.NewOptions(len(comic.ImageLinks), progressbar.OptionSetRenderBlankState(true), progressbar.OptionSetDescription(fmt.Sprintf("#%s", comic.IssueNumber)))
	}

	outputFormat := util.ImageType(comic.OutputImagesFormat)

	type downloadJob struct {
		index int
		link  string
	}

	jobs := make([]downloadJob, 0, len(comic.ImageLinks))
	for idx, link := range comic.ImageLinks {
		if strings.TrimSpace(link) == "" {
			continue
		}
		jobs = append(jobs, downloadJob{index: idx, link: link})
	}

	results := make([]string, len(comic.ImageLinks))
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)
	sem := semaphore.NewWeighted(int64(runtime.NumCPU()))
	var mu sync.Mutex
	const sniffLimit = 256

	for _, job := range jobs {
		job := job
		if err := sem.Acquire(ctx, 1); err != nil {
			return nil, err
		}
		group.Go(func() error {
			defer sem.Release(1)
			defer func() {
				if progress != nil {
					err := progress.Add(1)
					if err != nil && options.Logger != nil {
						options.Logger.Error(err.Error())
					}
				} else if options.Logger != nil {
					options.Logger.Infof("Downloaded image %d/%d", job.index+1, len(comic.ImageLinks))
				}
			}()

			reqCtx, cancelReq := context.WithTimeout(ctx, 30*time.Second)
			defer cancelReq()

			request, err := client.PrepareRequest(job.link, comic.Source.Name)
			if err != nil {
				return err
			}
			request = request.WithContext(reqCtx)

			response, err := client.DoRaw(request)
			if err != nil {
				return err
			}
			defer func() {
				if closeErr := response.Body.Close(); closeErr != nil {
					if options.Logger != nil {
						options.Logger.Errorf("failed to close response body for %s: %v", job.link, closeErr)
					}
				}
			}()

			if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
				if options.Logger != nil {
					options.Logger.Errorf("There was an error while downloading image number: %d - comic issue: %s (status code: %d)", job.index, comic.IssueNumber, response.StatusCode)
				}
				return nil
			}

			data, err := io.ReadAll(response.Body)
			if err != nil {
				if options.Logger != nil {
					options.Logger.Errorf("Failed reading image number: %d - comic issue: %s (%v)", job.index, comic.IssueNumber, err)
				}
				return nil
			}

			if len(data) == 0 {
				if options.Logger != nil {
					options.Logger.Errorf("Image number: %d - comic issue: %s returned an empty response body", job.index, comic.IssueNumber)
				}
				return nil
			}

			contentType := strings.ToLower(strings.TrimSpace(response.Header.Get("Content-Type")))
			if contentType == "" {
				sniffLen := len(data)
				if sniffLen > sniffLimit {
					sniffLen = sniffLimit
				}
				contentType = strings.ToLower(http.DetectContentType(data[:sniffLen]))
			}

			isWebp := strings.HasSuffix(strings.ToLower(job.link), ".webp") || strings.Contains(contentType, "image/webp")
			var inputImageFormat util.ImageFormat
			if isWebp {
				inputImageFormat = util.ImgFormatWEBP
			} else {
				inputImageFormat = util.ImageType(contentType)
			}
			if options.Logger != nil {
				// if the content type is present but does not indicate an image
				if contentType != "" && !strings.HasPrefix(contentType, "image/") {
					reportLen := len(data)
					if reportLen > sniffLimit {
						reportLen = sniffLimit
					}
					snippet := base64.StdEncoding.EncodeToString(data[:reportLen])
					options.Logger.Errorf("Unexpected content type '%s' while downloading image number: %d - url: %s (bytes=%d, snippet_base64=%s)", contentType, job.index, job.link, len(data), snippet)
				} else if inputImageFormat == util.ImgFormatUnknown {
					options.Logger.Warningf("Could not determine image format for image number: %d - comic issue: %s, content type: '%s', url: %s", job.index, comic.IssueNumber, contentType, job.link)
				}
			}

			fileName := fmt.Sprintf("%04d-image.%s", job.index, outputFormat)
			targetPath := filepath.Join(dir, fileName)
			imgFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}

			reader := bytes.NewReader(data)
			if err := util.SaveImage(options.Logger, imgFile, reader, outputFormat, inputImageFormat); err != nil {
				if options.Logger != nil {
					reportLen := len(data)
					if reportLen > sniffLimit {
						reportLen = sniffLimit
					}
					snippet := base64.StdEncoding.EncodeToString(data[:reportLen])
					options.Logger.Errorf("There was an error while downloading image number: %d - comic issue: %s (%v) (content-type=%s bytes=%d snippet_base64=%s)", job.index, comic.IssueNumber, err, contentType, len(data), snippet)
				}
				if closeErr := imgFile.Close(); closeErr != nil && options.Logger != nil {
					options.Logger.Errorf("failed to close image file %s: %v", targetPath, closeErr)
				}
				if removeErr := os.Remove(targetPath); removeErr != nil && options.Logger != nil {
					options.Logger.Errorf("failed to remove incomplete image %s: %v", targetPath, removeErr)
				}
			} else {
				if closeErr := imgFile.Close(); closeErr != nil && options.Logger != nil {
					options.Logger.Errorf("failed to close image file %s: %v", targetPath, closeErr)
				}
				mu.Lock()
				results[job.index] = targetPath
				mu.Unlock()
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

func (comic *ComicIssue) GetIssueNumAndVolume() string {
	if comic.Volume == nil {
		return fmt.Sprintf("c%s", comic.IssueNumber)
	}

	return fmt.Sprintf("v%s c%s", *comic.Volume, comic.IssueNumber)
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
func (comic *ComicIssue) MakeComic(options *config.Options) error {
	result, err := comic.DownloadImages(options)
	if err != nil {
		return err
	}
	defer func() {
		if err := os.RemoveAll(result.Dir); err != nil && options.Logger != nil {
			options.Logger.Errorf("failed to remove temporary directory %s: %v", result.Dir, err)
		}
	}()

	switch comic.OutputFormat {
	case EPUB:
		return comic.makeEPUB(options, result)
	case CBR, CBZ:
		return comic.makeCBRZ(options, result)
	default:
		return comic.makePDF(options, result)
	}
}
