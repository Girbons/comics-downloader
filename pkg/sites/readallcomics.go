package sites

import (
	"fmt"
	"strings"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
)

const DefaultUrl string = "https://readallcomics.com"

// Readallcomics  represents a Readallcomics instance.
type Readallcomics struct {
	options *config.Options
}

// NewReadallcomics returns a new Readallcomics instance.
func NewReadallcomics(options *config.Options) *Readallcomics {
	return &Readallcomics{
		options: options,
	}
}

func (r *Readallcomics) retrieveImageLinks(comic *core.Comic) ([]string, error) {
	var links []string

	response, err := soup.Get(comic.URLSource)

	if err != nil {
		return links, err
	}

	document := soup.HTMLParse(response)

	images := document.FindAll("img")
	for _, img := range images {
		url := img.Attrs()["src"]
		if util.IsURLValid(url) {
			links = append(links, url)
		}
	}

	if r.options.Debug {
		r.options.Logger.Debug(fmt.Sprintf("Image links found: %s", strings.Join(links, " ")))
	}

	return links, err
}

// Retrieve issues links from main comic page or from comic issue.
func (r *Readallcomics) getIssues(url string) ([]string, error) {
	var links []string

	response, err := soup.Get(url)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(response)

	if strings.Contains(url, "category") {
		chapters := doc.Find("ul", "class", "list-story").FindAll("a")
		for _, chapter := range chapters {
			issueUrl := chapter.Attrs()["href"]
			if util.IsURLValid(issueUrl) {
				links = append(links, issueUrl)
			}

		}
	} else {
		chapters := doc.Find("select", "id", "selectbox").FindAll("option")
		for _, chapter := range chapters {
			issueUrl := chapter.Attrs()["value"]
			if util.IsURLValid(issueUrl) {
				links = append(links, issueUrl)
			}
		}
	}

	if r.options.Debug {
		r.options.Logger.Debug(fmt.Sprintf("Image Links found: %s", strings.Join(links, " ")))
	}

	return links, nil
}

// RetrieveIssueLinks retrieves the links to all the issue.
func (r *Readallcomics) RetrieveIssueLinks() ([]string, error) {
	url := r.options.URL

	if r.options.All || r.options.Last {

		if (r.options.All || r.options.Last) && !strings.Contains(url, "category") {
			// override url to retrieve chapters
			comicName := strings.Join(strings.Split(url, "/")[3:], "")
			url = fmt.Sprintf("%s/category/%s", DefaultUrl, comicName)
		}

		chapters, err := r.getIssues(url)
		if err != nil {
			return nil, err
		}

		if r.options.Last {
			return []string{chapters[len(chapters)-1]}, nil
		}

		if r.options.All && len(chapters) == 1 {
			return []string{chapters[0]}, err
		}
		return chapters, err
	}

	return []string{url}, nil
}

// GetInfo extracts the comic info from the given URL.
func (r *Readallcomics) GetInfo(url string) (string, string) {
	parts := util.TrimAndSplitURL(url)

	lastPart := parts[len(parts)-1]
	title := strings.Replace(lastPart, "-", " ", -1)
	splittedTitle := strings.Split(title, " ")

	name := strings.Join(splittedTitle[:len(splittedTitle)-2], " ")
	issueNumber := splittedTitle[len(splittedTitle)-1] // get the issue number
	return name, issueNumber
}

// Initialize prepare the comic instance with links and images.
func (r *Readallcomics) Initialize(comic *core.Comic) error {
	links, err := r.retrieveImageLinks(comic)
	comic.Links = links

	return err
}
