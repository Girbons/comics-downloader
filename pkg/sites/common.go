package sites

import (
	"strings"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
	"github.com/dlclark/regexp2"
)

// mangakakalot.com and manganato.com functions

func MangaKakalotGetInfo(domain string, url string) (string, string) {
	// get chapter name
	res, err := soup.Get(url)
	if err != nil {
		return "", ""
	}
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
	name := f.Text()
	name, err = regexp2.MustCompile("(Vol\\.[0-9]{1,3} )?(Chapter [0-9]{1,3}(\\.[0-9])?) ?: ", 0).Replace(name, "", 0, 1)
	if err != nil {
		return "", ""
	}
	// parse number from url
	parts := util.TrimAndSplitURL(url)
	issueNumber := strings.Split(parts[len(parts)-1], "-")[1]
	return name, issueNumber
}

func MangaKakalotInitialize(comic *core.Comic) error {
	res, err := soup.Get(comic.URLSource)
	if err != nil {
		return err
	}
	doc := soup.HTMLParse(res)
	f := doc.Find("div", "class", "container-chapter-reader")
	var links []string
	for _, img := range f.FindAll("img") {
		links = append(links, img.Attrs()["src"])
	}
	comic.Links = links
	return nil
}

func MangaKakalotRetrieveIssueLinks(domain string, url string) ([]string, error) {
	res, err := soup.Get(url)
	if err != nil {
		panic(err)
	}
	// chapter page link
	if strings.Contains(url, "/chapter") {
		return []string{url}, nil
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
