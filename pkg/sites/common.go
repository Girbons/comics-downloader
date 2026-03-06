package sites

import (
	"context"
	"strings"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
	"github.com/dlclark/regexp2"
)

// mangakakalot.com and manganato.com functions

func mangaKakalotRequestContext(options *config.Options) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), options.RequestTimeout)
}

func MangaKakalotGetInfo(options *config.Options, domain string, url string) (name, issueNumber string, err error) {
	ctx, cancel := mangaKakalotRequestContext(options)
	defer cancel()

	res, err := fetchHTML(ctx, options.Client, url)
	if err != nil {
		return "", "", err
	}

	// get chapter name
	doc := soup.HTMLParse(res)
	f := doc.Find("div", "class", breadcrumbClassName(domain))
	switch {
	case strings.Contains(domain, "mangakakalot"):
		f = f.Find("p")
		items := f.FindAll("span", "itemprop", "itemListElement")
		f = items[len(items)-1]
		f = f.Find("a").Find("span")
	case strings.Contains(domain, "manganato"):
		items := f.FindAll("a", "class", "a-h")
		f = items[len(items)-1]
	}
	name = f.Text()
	name, err = regexp2.MustCompile("(Vol\\.[0-9]{1,3} )?(Chapter [0-9]{1,3}(\\.[0-9])?) ?: ", 0).Replace(name, "", 0, 1)
	if err != nil {
		return "", "", err
	}
	// parse number from url
	parts := util.TrimAndSplitURL(url)
	issueNumber = strings.Split(parts[len(parts)-1], "-")[1]
	return name, issueNumber, nil
}

func MangaKakalotInitialize(options *config.Options, comic *core.ComicIssue) error {
	ctx, cancel := mangaKakalotRequestContext(options)
	defer cancel()

	res, err := fetchHTML(ctx, options.Client, comic.Source.URL)
	if err != nil {
		return err
	}

	doc := soup.HTMLParse(res)
	f := doc.Find("div", "class", "container-chapter-reader")
	var links []string
	for _, img := range f.FindAll("img") {
		links = append(links, img.Attrs()["src"])
	}
	comic.ImageLinks = links
	return nil
}

func MangaKakalotRetrieveIssueLinks(options *config.Options, domain string, url string) ([]string, error) {
	// if chapter page, skip fetching and parsing the list page
	if strings.Contains(url, "/chapter") {
		return []string{url}, nil
	}

	ctx, cancel := mangaKakalotRequestContext(options)
	defer cancel()

	res, err := fetchHTML(ctx, options.Client, url)
	if err != nil {
		return nil, err
	}

	// manga page link
	doc := soup.HTMLParse(res)
	f := doc.Find("div", "class", chapterListClassName(domain))
	var urls []string
	for _, row := range f.FindAll(chapterListItemElementName(domain), "class", rowClassName(domain)) {
		urls = append(urls, row.Find("a").Attrs()["href"])
	}
	return urls, nil
}

// mangakakalot.com and (read)manganato.com domain-specific css class or element names

func rowClassName(domain string) string {
	switch {
	case strings.Contains(domain, "mangakakalot"):
		return "row"
	case strings.Contains(domain, "manganato"):
		return "a-h"
	}
	return ""
}

func chapterListClassName(domain string) string {
	switch {
	case strings.Contains(domain, "mangakakalot"):
		return "chapter-list"
	case strings.Contains(domain, "manganato"):
		return "panel-story-chapter-list"
	}
	return ""
}

func chapterListItemElementName(domain string) string {
	switch {
	case strings.Contains(domain, "mangakakalot"):
		return "div"
	case strings.Contains(domain, "manganato"):
		return "li"
	}
	return ""
}

func breadcrumbClassName(domain string) string {
	switch {
	case strings.Contains(domain, "mangakakalot"):
		return "breadcrumb"
	case strings.Contains(domain, "manganato"):
		return "panel-breadcrumb"
	}
	return ""
}
