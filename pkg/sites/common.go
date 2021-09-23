package sites

import (
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
	switch domain {
	case "mangakakalot.com":
		f = f.Find("p")
		items := f.FindAll("span", "itemprop", "itemListElement")
		f = items[len(items)-1]
		f = f.Find("a").Find("span")
	case "manganato.com":
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
	var fillPart string
	if len(parts) == 6 {
		fillPart = parts[3] + "/" + parts[4]
	} else {
		fillPart = parts[3]
	}
	exp := regexp2.MustCompile("(?<=https://(mangakakalot|manganato|readmanganato).com/"+fillPart+"/chapter(_|-))[0-9]{1,3}(_end)?(\\.[0-9])?", 0)
	match, err := exp.FindStringMatch(url)
	if match == nil || err != nil {
		return "", ""
	}
	issueNumber := match.String()
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
	switch domain {
	case "mangakakalot.com":
		return "row"
	case "manganato.com":
		return "a-h"
	}
	return ""
}

func chapterListClassName(domain string) string {
	switch domain {
	case "mangakakalot.com":
		return "chapter-list"
	case "manganato.com":
		return "panel-story-chapter-list"
	}
	return ""
}

func chapterListItemElementName(domain string) string {
	switch domain {
	case "mangakakalot.com":
		return "div"
	case "manganato.com":
		return "li"
	}
	return ""
}

func breadcrumbClassName(domain string) string {
	switch domain {
	case "mangakakalot.com":
		return "breadcrumb"
	case "manganato.com":
		return "panel-breadcrumb"
	}
	return ""
}
