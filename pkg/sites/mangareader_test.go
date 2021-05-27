package sites

//import (
//	"testing"
//
//	"github.com/Girbons/comics-downloader/internal/logger"
//	"github.com/Girbons/comics-downloader/pkg/config"
//	"github.com/Girbons/comics-downloader/pkg/core"
//	"github.com/stretchr/testify/assert"
//)
//
//const (
//	URL    = "https://www.mangareader.net/naruto/1/"
//	SOURCE = "www.mangareader.net"
//)
//
//func TestMangareaderGetInfo(t *testing.T) {
//	mr := new(Mangareader)
//	name, issueNumber := mr.GetInfo(URL)
//
//	assert.Equal(t, "naruto", name)
//	assert.Equal(t, "1", issueNumber)
//}
//
//func TestRetrieveMangareaderImageLinks(t *testing.T) {
//	opt :=
//		&config.Options{
//			URL:    URL,
//			Source: SOURCE,
//			All:    false,
//			Last:   false,
//			Debug:  false,
//			Logger: logger.NewLogger(false, make(chan string)),
//		}
//	mr := NewMangareader(opt)
//
//	comic := new(core.Comic)
//	comic.URLSource = URL
//	comic.Name = "naruto"
//	comic.IssueNumber = "1"
//	comic.Source = SOURCE
//
//	links, err := mr.retrieveImageLinks(comic)
//
//	assert.Equal(t, 53, len(links))
//	assert.Nil(t, err)
//}
//
//func TestSetupMangareader(t *testing.T) {
//	opt :=
//		&config.Options{
//			URL:    URL,
//			Source: SOURCE,
//			All:    false,
//			Last:   false,
//			Debug:  false,
//			Logger: logger.NewLogger(false, make(chan string)),
//		}
//	mr := NewMangareader(opt)
//
//	comic := new(core.Comic)
//	comic.Name = "naruto"
//	comic.IssueNumber = "1"
//	comic.URLSource = URL
//	comic.Source = SOURCE
//
//	err := mr.Initialize(comic)
//
//	assert.Nil(t, err)
//	assert.Equal(t, 53, len(comic.Links))
//}
//
//func TestMangareaderRetrieveIssueLinks(t *testing.T) {
//	opt :=
//		&config.Options{
//			URL:    "https://www.mangareader.net/naruto",
//			All:    false,
//			Last:   false,
//			Debug:  false,
//			Logger: logger.NewLogger(false, make(chan string)),
//		}
//	mr := NewMangareader(opt)
//	issues, err := mr.RetrieveIssueLinks()
//
//	assert.Nil(t, err)
//	assert.Equal(t, 700, len(issues))
//}
//
//func TestMangareaderRetrieveLastIssueLink(t *testing.T) {
//	opt :=
//		&config.Options{
//			URL:    "https://www.mangareader.net/naruto",
//			All:    false,
//			Last:   true,
//			Debug:  false,
//			Logger: logger.NewLogger(false, make(chan string)),
//		}
//	mr := NewMangareader(opt)
//	issue, err := mr.retrieveLastIssue("https://www.mangareader.net/naruto")
//
//	assert.Nil(t, err)
//	assert.Equal(t, "https://www.mangareader.net/naruto/700", issue)
//}
