package sites

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
)

func retrieveImageLinks(comic *core.Comic) ([]string, error) {
	var (
		pageLinks []string
		imgLinks  []string
	)
	// compile the image regex
	re := regexp.MustCompile(util.IMAGEREGEX)

	response, err := soup.Get(comic.URLSource)

	if err != nil {
		return nil, err
	}

	document := soup.HTMLParse(response)
	links := document.Find("div", "class", "cp-pager-list").Find("div").Find("span").FindAll("a")

	var pages []int
	for _, val := range links {
		index, _ := strconv.Atoi(val.Attrs()["data-page"])
		pages = append(pages, index)
	}

	lastPageIndex := util.FindMaxValueInSlice(pages)

	for i := 1; i <= lastPageIndex; i++ {
		newLink := fmt.Sprintf("http://www.mangahere.cc/manga/%s/%s/%s/%d.html", comic.Name, comic.SplitURL()[6], comic.IssueNumber, i)
		if !util.IsValueInSlice(newLink, pageLinks) {
			pageLinks = append(pageLinks, newLink)
		}
	}

	for _, link := range pageLinks {
		if link != "" {
			imgResponse, imgResponseError := soup.Get(link)

			if imgResponseError != nil {
				return nil, imgResponseError
			}

			match := re.FindAllStringSubmatch(imgResponse, -1)
			for _, f := range match {
				if util.IsURLValid(f[1]) {
					imgLinks = append(imgLinks, f[1])
				}
			}
		}
	}
	return imgLinks, err
}

// SetupMangaHere will initialize the comic based
// on mangahere.cc
func Initialize(comic *core.Comic) error {
	comic.Name = comic.SplitURL()[4]
	comic.IssueNumber = comic.SplitURL()[5]
	links, err := retrieveImageLinks(comic)
	comic.Links = links

	return err
}
