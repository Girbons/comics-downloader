package sites

import (
	"fmt"
	"strings"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
)

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

func (r *Readallcomics) isSingleIssue(url string) bool {
	response, _ := soup.Get(url)

	doc := soup.HTMLParse(response)
	chapters := doc.Find("select", "id", "selectbox").FindAll("option")

	return len(chapters) == 1
}

func (r *Readallcomics) retrieveLastIssue(url string) (string, error) {
	response, err := soup.Get(url)

	doc := soup.HTMLParse(response)
	chapters := doc.Find("select", "id", "selectbox").FindAll("option")

	return chapters[0].Attrs()["value"], err
}

// RetrieveIssueLinks retrieves the links to all the issue.
func (r *Readallcomics) RetrieveIssueLinks() ([]string, error) {
	url := r.options.URL

	if r.options.Last {
		lastIssue, err := r.retrieveLastIssue(url)
		return []string{lastIssue}, err
	}

	if r.options.All && r.isSingleIssue(url) {
		return []string{url}, nil
	}

	response, err := soup.Get(url)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(response)
	chapters := doc.Find("select", "id", "selectbox").FindAll("option")

	var links []string

	for _, chapter := range chapters {
		url := chapter.Attrs()["value"]
		if util.IsURLValid(url) {
			links = append(links, url)
		}
	}

	return links, err
}

// GetInfo extracts the comic info from the given URL.
func (r *Readallcomics) GetInfo(url string) (string, string) {
	parts := util.TrimAndSplitURL(url)

	lastPart := parts[len(parts)-1]
	title := strings.Replace(lastPart, "-", " ", -1)
	splittedTitle := strings.Split(title, " ")

	name := strings.Join(splittedTitle[:len(splittedTitle)-2], " ")
	issueNumber := splittedTitle[len(splittedTitle)-2] // get the issue number
	return name, issueNumber
}

// Initialize prepare the comic instance with links and images.
func (r *Readallcomics) Initialize(comic *core.Comic) error {
	links, err := r.retrieveImageLinks(comic)
	comic.Links = links

	return err
}
