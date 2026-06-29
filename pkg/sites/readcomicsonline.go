package sites

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
)

// readComicsOnlineBaseURL is the canonical root of the site.
const readComicsOnlineBaseURL = "https://readcomicsonline.ru"

// ReadComicsOnline represents a readcomicsonline.ru scraper instance.
type ReadComicsOnline struct {
	options *config.Options
	baseURL string // overridable for tests
}

// NewReadComicsOnline returns a ReadComicsOnline instance.
func NewReadComicsOnline(options *config.Options) *ReadComicsOnline {
	return &ReadComicsOnline{
		options: options,
		baseURL: readComicsOnlineBaseURL,
	}
}

// pageEntry mirrors one element of the JS `var pages = [...]` array.
type pageEntry struct {
	PageImage string `json:"page_image"`
}

// retrieveImageLinks fetches the issue page and extracts all image URLs.
//
// Primary strategy  – collect `data-src` attributes on <img class="img-responsive">
// elements (lazy-load pattern used by the site).
//
// Fallback strategy – if the primary pass yields nothing, parse the inline JS
// `var pages = [{"page_image":"01.jpg",...}]` variable and reconstruct absolute
// URLs from the known path pattern:
//
//	{baseURL}/uploads/manga/{slug}/chapters/{issue}/{filename}
func (r *ReadComicsOnline) retrieveImageLinks(comic *core.Comic) ([]string, error) {
	response, err := soup.Get(comic.URLSource)
	if err != nil {
		return nil, fmt.Errorf("readcomicsonline: fetch %s: %w", comic.URLSource, err)
	}

	links := r.extractDataSrcLinks(response)

	if len(links) == 0 {
		links, err = r.extractFromPagesVar(response, comic.URLSource)
		if err != nil {
			return nil, err
		}
	}

	if r.options.Debug && r.options.Logger != nil {
		r.options.Logger.Debugf("readcomicsonline: found %d image links for %s", len(links), comic.URLSource)
	}

	return links, nil
}

// extractDataSrcLinks parses the HTML and collects trimmed `data-src` values
// from every <img class="img-responsive"> element.
func (r *ReadComicsOnline) extractDataSrcLinks(html string) []string {
	var links []string
	doc := soup.HTMLParse(html)

	for _, img := range doc.FindAll("img", "class", "img-responsive") {
		attrs := img.Attrs()
		src, ok := attrs["data-src"]
		if !ok {
			// some pages serve the real src directly when JS is disabled
			src, ok = attrs["src"]
			if !ok {
				continue
			}
		}
		src = strings.TrimSpace(src)
		if util.IsURLValid(src) && !util.IsValueInSlice(src, links) {
			links = append(links, src)
		}
	}

	return links
}

// pagesVarRe matches: var pages = [...];
var pagesVarRe = regexp.MustCompile(`var\s+pages\s*=\s*(\[.*?\]);`)

// extractFromPagesVar parses the inline JS `var pages` array and reconstructs
// absolute image URLs from the slug and issue number embedded in the issue URL.
func (r *ReadComicsOnline) extractFromPagesVar(html, issueURL string) ([]string, error) {
	m := pagesVarRe.FindStringSubmatch(html)
	if m == nil {
		return nil, nil // not an error – page may just have no images yet
	}

	var entries []pageEntry
	if err := json.Unmarshal([]byte(m[1]), &entries); err != nil {
		return nil, fmt.Errorf("readcomicsonline: parse pages var: %w", err)
	}

	// derive slug and issue number from the URL
	// URL shape: {base}/comic/{slug}/{issue}
	parts := util.TrimAndSplitURL(issueURL)
	if len(parts) < 2 {
		return nil, fmt.Errorf("readcomicsonline: cannot derive slug/issue from %s", issueURL)
	}
	slug := parts[len(parts)-2]
	issue := parts[len(parts)-1]

	base := fmt.Sprintf("%s/uploads/manga/%s/chapters/%s/", r.baseURL, slug, issue)

	var links []string
	for _, e := range entries {
		if e.PageImage == "" {
			continue
		}
		url := base + e.PageImage
		if !util.IsValueInSlice(url, links) {
			links = append(links, url)
		}
	}

	return links, nil
}

// isSingleIssue returns true when the URL points to a specific issue number,
// e.g. /comic/some-slug/1  (has a numeric final segment).
func (r *ReadComicsOnline) isSingleIssue(url string) bool {
	parts := util.TrimAndSplitURL(url)
	if len(parts) < 2 {
		return false
	}
	last := parts[len(parts)-1]
	return isNumeric(last)
}

// retrieveIssueListFromComicPage fetches the comic's landing page and collects
// all issue URLs from the <ul class="chapters"> block.
func (r *ReadComicsOnline) retrieveIssueListFromComicPage(url string) ([]string, error) {
	response, err := soup.Get(url)
	if err != nil {
		return nil, fmt.Errorf("readcomicsonline: fetch %s: %w", url, err)
	}

	doc := soup.HTMLParse(response)
	chapterList := doc.Find("ul", "class", "chapters")
	if chapterList.Error != nil {
		return nil, fmt.Errorf("readcomicsonline: chapter list not found on %s", url)
	}

	var links []string
	for _, a := range chapterList.FindAll("a") {
		href, ok := a.Attrs()["href"]
		if !ok {
			continue
		}
		href = strings.TrimSpace(href)
		if href == "" {
			continue
		}
		if !strings.HasPrefix(href, "http") {
			href = r.baseURL + href
		}
		// strip query strings
		href = strings.Split(href, "?")[0]
		if util.IsURLValid(href) && !util.IsValueInSlice(href, links) {
			links = append(links, href)
		}
	}

	return links, nil
}

// comicPageFromIssueURL derives the comic landing-page URL from an issue URL.
// e.g. https://readcomicsonline.ru/comic/slug/3  →  https://readcomicsonline.ru/comic/slug
func comicPageFromIssueURL(issueURL string) string {
	parts := util.TrimAndSplitURL(issueURL)
	if len(parts) < 2 {
		return issueURL
	}
	return strings.Join(parts[:len(parts)-1], "/")
}

// RetrieveIssueLinks returns a slice of issue URLs to download.
func (r *ReadComicsOnline) RetrieveIssueLinks() ([]string, error) {
	url := r.options.URL

	if !r.isSingleIssue(url) {
		// URL is a comic listing page
		issues, err := r.retrieveIssueListFromComicPage(url)
		if err != nil {
			return nil, err
		}

		if r.options.Last && len(issues) > 0 {
			return []string{issues[0]}, nil
		}

		return issues, nil
	}

	// URL is already a single issue
	if r.options.All {
		comicPage := comicPageFromIssueURL(url)
		issues, err := r.retrieveIssueListFromComicPage(comicPage)
		if err != nil {
			return nil, err
		}
		if r.options.Last && len(issues) > 0 {
			return []string{issues[0]}, nil
		}
		return issues, nil
	}

	return []string{url}, nil
}

// GetInfo extracts the comic name and issue number from a URL.
// URL shape: https://readcomicsonline.ru/comic/{slug}/{issue}
func (r *ReadComicsOnline) GetInfo(url string) (string, string) {
	parts := util.TrimAndSplitURL(url)
	// minimum: ["https:", "", "readcomicsonline.ru", "comic", slug, issue]
	//           0         1   2                      3        4    5
	if len(parts) < 6 {
		return "", ""
	}
	name := parts[4]
	issueNumber := parts[5]
	return name, issueNumber
}

// Initialize populates comic.Links for the given issue.
func (r *ReadComicsOnline) Initialize(comic *core.Comic) error {
	links, err := r.retrieveImageLinks(comic)
	comic.Links = links
	return err
}
