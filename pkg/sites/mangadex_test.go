package sites

//import (
//"testing"

//"github.com/Girbons/comics-downloader/internal/logger"
//"github.com/Girbons/comics-downloader/pkg/config"
//"github.com/Girbons/comics-downloader/pkg/core"
//"github.com/stretchr/testify/assert"
//)

//const testMangadexBase string = "mangadex.org"
//const testMangadexURL string = "https://" + testMangadexBase + "/"

//func TestMangadexGetInfo(t *testing.T) {
//opt := &config.Options{
//Url:     testMangadexURL + "chapter/155061/1",
//Country: "",
//Source:  testMangadexBase,
//Debug:   false,
//Logger:  logger.NewLogger(false, make(chan string)),
//}
//md := NewMangadex(opt)
//name, issueNumber := md.GetInfo(testMangadexURL + "chapter/155061/1")

//assert.Equal(t, "Naruto", name)
//assert.Equal(t, "Vol 60 Chapter 575, A Will of Stone", issueNumber)
//}

//func TestMangadexSetup(t *testing.T) {
//opt := &config.Options{
//Url:     testMangadexURL + "chapter/155061/1",
//Country: "",
//Source:  testMangadexBase,
//Debug:   false,
//Logger:  logger.NewLogger(false, make(chan string)),
//}
//comic := new(core.Comic)
//comic.URLSource = testMangadexURL + "chapter/155061/1"

//md := NewMangadex(opt)
//err := md.Initialize(comic)

//assert.Nil(t, err)
//assert.Equal(t, 14, len(comic.Links))
//}

//func TestMangadexRetrieveIssueLinks(t *testing.T) {
//opt := &config.Options{
//Url:     testMangadexURL + "chapter/155061/",
//Country: "",
//Source:  testMangadexBase,
//Last:    false,
//All:     false,
//Debug:   false,
//Logger:  logger.NewLogger(false, make(chan string)),
//}
//md := NewMangadex(opt)
//urls, err := md.RetrieveIssueLinks()
//assert.Nil(t, err)
//assert.Equal(t, 1, len(urls))
//}

//func TestMangadexRetrieveIssueLinksAllChapter(t *testing.T) {
//opt := &config.Options{
//Url:     testMangadexURL + "title/5/naruto/",
//Country: "gb",
//Source:  testMangadexBase,
//Last:    false,
//All:     true,
//Debug:   false,
//Logger:  logger.NewLogger(false, make(chan string)),
//}
//md := NewMangadex(opt)
//urls, err := md.RetrieveIssueLinks()
//assert.Nil(t, err)
//assert.Len(t, urls, 713)
//}

//func TestMangadexRetrieveIssueLinksLastChapter(t *testing.T) {
//opt := &config.Options{
//Url:     testMangadexURL + "title/5/naruto/",
//Country: "gb",
//Source:  testMangadexBase,
//Last:    true,
//All:     false,
//Debug:   false,
//Logger:  logger.NewLogger(false, make(chan string)),
//}
//md := NewMangadex(opt)
//urls, err := md.RetrieveIssueLinks()
//assert.Nil(t, err)
//assert.Len(t, urls, 1)
//}

//func TestMangadexUnsupportedURL(t *testing.T) {
//opt := &config.Options{
//Url:     testMangadexURL,
//Country: "",
//Source:  testMangadexBase,
//Last:    false,
//All:     false,
//Debug:   false,
//Logger:  logger.NewLogger(false, make(chan string)),
//}
//md := NewMangadex(opt)
//_, err := md.RetrieveIssueLinks()
//assert.EqualError(t, err, "URL not supported")

//md.options.Url = testMangadexURL + "test/0/"
//_, err = md.RetrieveIssueLinks()
//assert.EqualError(t, err, "URL not supported")
//}

//func TestMangadexNoManga(t *testing.T) {
//opt := &config.Options{
//Url:     testMangadexURL + "title/0/",
//Country: "",
//Source:  testMangadexBase,
//Last:    false,
//All:     false,
//Debug:   false,
//Logger:  logger.NewLogger(false, make(chan string)),
//}
//md := NewMangadex(opt)
//_, err := md.RetrieveIssueLinks()
//assert.Error(t, err)
//assert.Contains(t, err.Error(), "could not get manga 0")
//}

//func TestMangadexNoChapters(t *testing.T) {
//opt := &config.Options{
//Url:     testMangadexURL + "title/5/naruto/",
//Country: "xyz",
//Source:  testMangadexBase,
//Last:    false,
//All:     true,
//Debug:   false,
//Logger:  logger.NewLogger(false, make(chan string)),
//}
//md := NewMangadex(opt)
//_, err := md.RetrieveIssueLinks()
//assert.EqualError(t, err, "no chapters found")
//}
