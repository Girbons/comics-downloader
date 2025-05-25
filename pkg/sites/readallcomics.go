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

	splittedByHyphen := strings.Split(lastPart, "-")

	if len(splittedByHyphen) <= 1 {
		title := strings.ReplaceAll(lastPart, "-", " ")
		splittedTitle := strings.Split(title, " ")
		name := strings.Join(splittedTitle[:len(splittedTitle)-2], " ")
		issueNumber := splittedTitle[len(splittedTitle)-1]
		return name, issueNumber
	}

	var issueIndices []int

	for i, part := range splittedByHyphen {
		if len(part) == 3 && isNumeric(part) {
			issueIndices = append(issueIndices, i)
		} else if len(part) >= 2 && strings.HasPrefix(strings.ToLower(part), "v") && isNumeric(part[1:]) {
			issueIndices = append(issueIndices, i)
		} else if len(part) <= 2 && isNumeric(part) {
			issueIndices = append(issueIndices, i)
		} else if len(part) == 4 && part >= "1900" && part <= "2099" {
			issueIndices = append(issueIndices, i)
		} else if strings.Contains(part, "029") {
			for j := 0; j <= len(part)-3; j++ {
				if j+3 <= len(part) && isNumeric(part[j:j+3]) {
					issueIndices = append(issueIndices, i)
					break
				}
			}
		}
	}

	splitIndex := -1

	if len(issueIndices) > 0 {
		for _, idx := range issueIndices {
			part := splittedByHyphen[idx]
			if (len(part) == 3 && isNumeric(part)) ||
				(len(part) >= 2 && strings.HasPrefix(strings.ToLower(part), "v") && isNumeric(part[1:])) ||
				strings.Contains(part, "029") {
				splitIndex = idx
				break
			}
		}

		if splitIndex == -1 {
			splitIndex = issueIndices[0]
		}
	}

	if splitIndex > 0 {
		name := strings.Join(splittedByHyphen[:splitIndex], "-")

		issueNumberPart := splittedByHyphen[splitIndex]
		var issueNumber string

		if strings.Contains(issueNumberPart, "029") || (len(issueNumberPart) > 3 && !isNumeric(issueNumberPart)) {
			found := false
			for i := 0; i <= len(issueNumberPart)-3; i++ {
				if i+3 <= len(issueNumberPart) && isNumeric(issueNumberPart[i:i+3]) {
					issueNumber = issueNumberPart[i : i+3]
					found = true
					break
				}
			}
			if !found {
				for i := 0; i <= len(issueNumberPart)-1; i++ {
					for length := 2; length >= 1; length-- {
						if i+length <= len(issueNumberPart) && isNumeric(issueNumberPart[i:i+length]) {
							issueNumber = issueNumberPart[i : i+length]
							found = true
							break
						}
					}
					if found {
						break
					}
				}
			}
			if !found {
				issueNumber = issueNumberPart
			}

			if splitIndex < len(splittedByHyphen)-1 {
				remainingParts := splittedByHyphen[splitIndex+1:]
				for _, part := range remainingParts {
					if len(part) == 4 && part >= "1900" && part <= "2099" {
						issueNumber += "-" + part
						break
					}
				}
			}
		} else {
			issueNumber = strings.Join(splittedByHyphen[splitIndex:], "-")
		}

		name = strings.ReplaceAll(name, "-", " ")

		return name, issueNumber
	}

	if len(splittedByHyphen) >= 2 {
		lastPart := splittedByHyphen[len(splittedByHyphen)-1]
		secondLastPart := splittedByHyphen[len(splittedByHyphen)-2]

		if len(lastPart) == 4 && lastPart >= "1900" && lastPart <= "2099" {
			name := strings.Join(splittedByHyphen[:len(splittedByHyphen)-2], "-")
			issueNumber := secondLastPart + "-" + lastPart

			name = strings.ReplaceAll(name, "-", " ")

			return name, issueNumber
		}
	}

	name := strings.Join(splittedByHyphen[:len(splittedByHyphen)-1], "-")
	issueNumber := splittedByHyphen[len(splittedByHyphen)-1]

	name = strings.ReplaceAll(name, "-", " ")

	return name, issueNumber
}

// isNumeric checks if a string contains only digits
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// Initialize prepare the comic instance with links and images.
func (r *Readallcomics) Initialize(comic *core.Comic) error {
	links, err := r.retrieveImageLinks(comic)
	comic.Links = links

	return err
}
