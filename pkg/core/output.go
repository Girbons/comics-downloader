package core

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Girbons/comics-downloader/internal/version"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/beevik/etree"
)

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

// makeComicInfoXML generates a ComicInfo.xml file for the given comic issue and saves it to the output directory. It returns the path to the generated ComicInfo.xml file.
// Based on the ComicInfo.xml https://anansi-project.github.io/docs/comicinfo/schemas/v2.1
func (comic *ComicIssue) makeComicInfoXML(options *config.Options, images *DownloadResult) (string, error) {
	outputDir, err := util.ImagesPathSetup(options.CreateDefaultPath, options.OutputFolder, comic.Source.Name, comic.Name, options.IssueFolderName, comic.IssueNumber)
	if err != nil {
		return "", err
	}

	comicInfoPath := filepath.Join(outputDir, "ComicInfo.xml")
	options.Logger.Infof("ComicInfo.xml path: %s", comicInfoPath)

	fo, err := os.Create(comicInfoPath)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)

	comicInfo := doc.CreateElement("ComicInfo")

	comicInfo.CreateAttr("xmlns:xsi", "http://www.w3.org/2001/XMLSchema-instance")
	comicInfo.CreateAttr("xsi:noNamespaceSchemaLocation", "https://github.com/anansi-project/comicinfo/raw/db8e1d84132f97403b226f2e12aaec1342c2a223/drafts/v2.1/ComicInfo.xsd")
	comicInfo.CreateElement("Notes").SetText(fmt.Sprintf("Tagged by comics-downloader version %s using info from %s at %s", version.Tag, comic.Source.Name, time.Now().Format(time.RFC3339)))

	comicInfo.CreateElement("Series").SetText(comic.SeriesMetadata.Title)
	comicInfo.CreateElement("Title").SetText(comic.Name)

	localizedTitle := comic.getLocalizedTitle(options)
	comicInfo.CreateElement("LocalizedSeries").SetText(localizedTitle)

	comicDescription := comic.getLocalizedDescription(options)
	comicInfo.CreateElement("Summary").SetText(comicDescription)

	comicInfo.CreateElement("Number").SetText(comic.IssueNumber)
	if comic.Volume != nil {
		comicInfo.CreateElement("Volume").SetText(*comic.Volume)
	}

	if comic.ReleaseDate != nil {
		comicInfo.CreateElement("Year").SetText(fmt.Sprintf("%d", comic.ReleaseDate.Year()))
		comicInfo.CreateElement("Month").SetText(fmt.Sprintf("%02d", comic.ReleaseDate.Month()))
		comicInfo.CreateElement("Day").SetText(fmt.Sprintf("%02d", comic.ReleaseDate.Day()))
	}
	if comic.LanguageISO != nil {
		comicInfo.CreateElement("LanguageISO").SetText(*comic.LanguageISO)
	}

	if comic.SeriesMetadata.IsManga != nil {
		// TODO: consider adding unknown value for manga field if IsManga is nil?
		mangaTag := comicInfo.CreateElement("Manga")
		if *comic.SeriesMetadata.IsManga {
			if comic.SeriesMetadata.IsRTL != nil && *comic.SeriesMetadata.IsRTL {
				mangaTag.SetText("YesAndRightToLeft")
			} else {
				mangaTag.SetText("Yes")
			}
		} else {
			mangaTag.SetText("No")
		}
	}
	if comic.SeriesMetadata.AgeRating != nil {
		comicInfo.CreateElement("AgeRating").SetText(string(*comic.SeriesMetadata.AgeRating))
	}
	if comic.SeriesMetadata.CommunityRating != nil {
		comicInfo.CreateElement("CommunityRating").SetText(fmt.Sprintf("%.2f", *comic.SeriesMetadata.CommunityRating))
	}
	if len(comic.SeriesMetadata.Tags) > 0 {
		comicInfo.CreateElement("Tags").SetText(strings.Join(comic.SeriesMetadata.Tags, ","))
	}
	if len(comic.SeriesMetadata.Genres) > 0 {
		comicInfo.CreateElement("Genres").SetText(strings.Join(comic.SeriesMetadata.Genres, ","))
	}
	if len(comic.SeriesMetadata.WebLinks) > 0 {
		var cleanedWebLinks []string
		for _, link := range comic.SeriesMetadata.WebLinks {
			// the links must be URL-encoded as spaces are the separator
			cleanedWebLinks = append(cleanedWebLinks, url.QueryEscape(link))
		}
		comicInfo.CreateElement("WebLinks").SetText(strings.Join(cleanedWebLinks, " "))
	}
	if comic.SeriesMetadata.AgeRating != nil {
		comicInfo.CreateElement("AgeRating").SetText(string(*comic.SeriesMetadata.AgeRating))
	}
	if len(comic.SeriesMetadata.Creators) > 0 {
		var writers []string
		var pencillers []string
		var inkers []string
		var colorists []string
		var letterers []string
		var coverArtists []string
		var editors []string
		var translators []string

		for _, creator := range comic.SeriesMetadata.Creators {
			switch creator.Role {
			case CreatorRoleWriter:
				writers = append(writers, creator.Name)
			case CreatorRolePenciller:
				pencillers = append(pencillers, creator.Name)
			case CreatorRoleInker:
				inkers = append(inkers, creator.Name)
			case CreatorRoleColorist:
				colorists = append(colorists, creator.Name)
			case CreatorRoleLetterer:
				letterers = append(letterers, creator.Name)
			case CreatorRoleCoverArtist:
				coverArtists = append(coverArtists, creator.Name)
			case CreatorRoleEditor:
				editors = append(editors, creator.Name)
			case CreatorRoleTranslator:
				translators = append(translators, creator.Name)
			}
		}
		if len(writers) > 0 {
			comicInfo.CreateElement("Writer").SetText(strings.Join(writers, ","))
		}
		if len(pencillers) > 0 {
			comicInfo.CreateElement("Penciller").SetText(strings.Join(pencillers, ","))
		}
		if len(inkers) > 0 {
			comicInfo.CreateElement("Inker").SetText(strings.Join(inkers, ","))
		}
		if len(colorists) > 0 {
			comicInfo.CreateElement("Colorist").SetText(strings.Join(colorists, ","))
		}
		if len(letterers) > 0 {
			comicInfo.CreateElement("Letterer").SetText(strings.Join(letterers, ","))
		}
		if len(coverArtists) > 0 {
			comicInfo.CreateElement("CoverArtist").SetText(strings.Join(coverArtists, ","))
		}
		if len(editors) > 0 {
			comicInfo.CreateElement("Editor").SetText(strings.Join(editors, ","))
		}
		if len(translators) > 0 {
			comicInfo.CreateElement("Translator").SetText(strings.Join(translators, ","))
		}
	}
	comicInfo.CreateElement("PageCount").SetText(fmt.Sprintf("%d", len(images.FilePaths)))

	doc.Indent(2)
	_, err = doc.WriteTo(fo)
	if err != nil {
		return "", err
	}

	return comicInfoPath, nil
}
