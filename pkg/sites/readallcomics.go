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
	urlParts := strings.Split(lastPart, "-")

	// Handle simple case with no hyphens
	if len(urlParts) <= 1 {
		return r.parseSimpleFormat(lastPart)
	}

	// Find potential issue number indices
	issueIndices := r.findIssueIndices(urlParts)

	// Determine the best split point
	splitIndex := r.determineBestSplitIndex(urlParts, issueIndices)

	// Extract name and issue number based on split index
	if splitIndex > 0 {
		return r.extractInfoWithSplitIndex(urlParts, splitIndex)
	}

	// Handle year suffix pattern (e.g., "name-issue-year")
	if r.hasYearSuffix(urlParts) {
		return r.parseYearSuffixFormat(urlParts)
	}

	// Default fallback: last part is issue number
	return r.parseDefaultFormat(urlParts)
}

// parseSimpleFormat handles URLs with no hyphens in the last part
func (r *Readallcomics) parseSimpleFormat(lastPart string) (string, string) {
	title := strings.ReplaceAll(lastPart, "-", " ")
	splittedTitle := strings.Split(title, " ")
	name := strings.Join(splittedTitle[:len(splittedTitle)-2], " ")
	issueNumber := splittedTitle[len(splittedTitle)-1]
	return name, issueNumber
}

// findIssueIndices identifies parts of the URL that could be issue numbers
func (r *Readallcomics) findIssueIndices(urlParts []string) []int {
	var issueIndices []int

	for i, part := range urlParts {
		if r.isThreeDigitIssue(part) ||
			r.isVersionNumber(part) ||
			r.isShortNumeric(part) ||
			r.isYear(part) ||
			r.hasEmbeddedNumber(part) {
			issueIndices = append(issueIndices, i)
		}
	}

	return issueIndices
}

// isThreeDigitIssue checks if a part is a 3-digit numeric issue number
func (r *Readallcomics) isThreeDigitIssue(part string) bool {
	return len(part) == 3 && isNumeric(part)
}

// isVersionNumber checks if a part is a version number (e.g., "v2")
func (r *Readallcomics) isVersionNumber(part string) bool {
	return len(part) >= 2 &&
		strings.HasPrefix(strings.ToLower(part), "v") &&
		isNumeric(part[1:])
}

// isShortNumeric checks if a part is a short numeric value (1-2 digits)
func (r *Readallcomics) isShortNumeric(part string) bool {
	return len(part) <= 2 && isNumeric(part)
}

// isYear checks if a part represents a year (1900-2099)
func (r *Readallcomics) isYear(part string) bool {
	return len(part) == 4 && part >= "1900" && part <= "2099"
}

// hasEmbeddedNumber checks if a part contains embedded numbers (e.g., "029something")
func (r *Readallcomics) hasEmbeddedNumber(part string) bool {
	if !strings.Contains(part, "029") {
		return false
	}

	for j := 0; j <= len(part)-3; j++ {
		if j+3 <= len(part) && isNumeric(part[j:j+3]) {
			return true
		}
	}
	return false
}

// determineBestSplitIndex finds the optimal point to split name from issue info
func (r *Readallcomics) determineBestSplitIndex(urlParts []string, issueIndices []int) int {
	if len(issueIndices) == 0 {
		return -1
	}

	// Look for high-priority patterns first
	for _, idx := range issueIndices {
		part := urlParts[idx]
		if r.isThreeDigitIssue(part) ||
			r.isVersionNumber(part) ||
			r.hasEmbeddedNumber(part) {
			return idx
		}
	}

	// Use the first available index as fallback
	return issueIndices[0]
}

// extractInfoWithSplitIndex extracts name and issue using the determined split point
func (r *Readallcomics) extractInfoWithSplitIndex(urlParts []string, splitIndex int) (string, string) {
	name := strings.Join(urlParts[:splitIndex], "-")
	name = strings.ReplaceAll(name, "-", " ")

	issueNumber := r.extractIssueNumber(urlParts, splitIndex)

	return name, issueNumber
}

// extractIssueNumber extracts the issue number from the URL parts starting at splitIndex
func (r *Readallcomics) extractIssueNumber(urlParts []string, splitIndex int) string {
	issueNumberPart := urlParts[splitIndex]

	// Handle complex issue number extraction for embedded numbers
	if strings.Contains(issueNumberPart, "029") ||
		(len(issueNumberPart) > 3 && !isNumeric(issueNumberPart)) {

		if extractedNumber := r.extractEmbeddedNumber(issueNumberPart); extractedNumber != "" {
			// For embedded numbers, don't add year suffix - just return the extracted number
			return extractedNumber
		}

		// If we couldn't extract an embedded number, check for year suffix
		if yearSuffix := r.findYearInRemainingParts(urlParts, splitIndex); yearSuffix != "" {
			return issueNumberPart + "-" + yearSuffix
		}
		return issueNumberPart
	}

	// Simple case: join all parts from split index onward
	return strings.Join(urlParts[splitIndex:], "-")
}

// extractEmbeddedNumber finds numeric patterns within a complex part
func (r *Readallcomics) extractEmbeddedNumber(part string) string {
	// Try 3-digit patterns first
	for i := 0; i <= len(part)-3; i++ {
		if i+3 <= len(part) && isNumeric(part[i:i+3]) {
			return part[i : i+3]
		}
	}

	// Try shorter patterns
	for i := 0; i <= len(part)-1; i++ {
		for length := 2; length >= 1; length-- {
			if i+length <= len(part) && isNumeric(part[i:i+length]) {
				return part[i : i+length]
			}
		}
	}

	return ""
}

// findYearInRemainingParts looks for year information in parts after the split index
func (r *Readallcomics) findYearInRemainingParts(urlParts []string, splitIndex int) string {
	if splitIndex < len(urlParts)-1 {
		remainingParts := urlParts[splitIndex+1:]
		for _, part := range remainingParts {
			if r.isYear(part) {
				return part
			}
		}
	}
	return ""
}

// hasYearSuffix checks if the URL has a year as the last component
func (r *Readallcomics) hasYearSuffix(urlParts []string) bool {
	if len(urlParts) < 2 {
		return false
	}
	lastPart := urlParts[len(urlParts)-1]
	return r.isYear(lastPart)
}

// parseYearSuffixFormat handles URLs ending with year (e.g., "name-issue-2023")
func (r *Readallcomics) parseYearSuffixFormat(urlParts []string) (string, string) {
	lastPart := urlParts[len(urlParts)-1]
	secondLastPart := urlParts[len(urlParts)-2]

	name := strings.Join(urlParts[:len(urlParts)-2], "-")
	issueNumber := secondLastPart + "-" + lastPart
	name = strings.ReplaceAll(name, "-", " ")

	return name, issueNumber
}

// parseDefaultFormat handles the fallback case where last part is the issue number
func (r *Readallcomics) parseDefaultFormat(urlParts []string) (string, string) {
	name := strings.Join(urlParts[:len(urlParts)-1], "-")
	issueNumber := urlParts[len(urlParts)-1]
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
