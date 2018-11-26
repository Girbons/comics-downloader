package sites

//import (
//"image"
//"net/http"

//"github.com/BakeRolls/mri"
//"github.com/Girbons/comics-downloader/pkg/core"
//"github.com/Girbons/mangarock"
//log "github.com/sirupsen/logrus"
//)

//func retrieveMangaRockImages(links []string) []*image.Image {
//var (
//format   string
//response *http.Response
//err      error
//)

//for _, link := range links {
//response, err = http.Get(link)

//if err != nil {
//log.Error(err)
//}
//defer response.Body.Close()

//decodedImage, decodeErr := mri.Decode(response.Body)

//if decodeErr != nil {
//log.Error(decodeErr)
//}

//// images = append(images, decodedImage)
//}
//return images
//}

//func findChapterName(chapterID string, chapters []*mangarock.Chapter) (string, bool) {
//for _, chapter := range chapters {
//if chapter.Oid == chapterID {
//return chapter.Name, true
//}
//}
//return "", false
//}

//func SetupMangaRock(c *core.Comic, splittedUrl []string) {
//var (
//chapterID string
//chapter   string
//name      string
//series    string
//found     bool

//client *mangarock.Client
//info   *mangarock.MangaRockInfo
//pages  *mangarock.MangaRockPages
//images []*image.Image

//infoErr  error
//pagesErr error
//)

//series = splittedUrl[4]
//chapterID = splittedUrl[6]

//client = mangarock.NewClient()
//// get info about the manga
//info, infoErr = client.Info(series)
//if infoErr != nil {
//log.Error(infoErr)
//}
//// retrieve pages
//pages, pagesErr = client.Pages(chapterID)
//if pagesErr != nil {
//log.Error(pagesErr)
//}

//name = info.Data.Name
//chapter, found = findChapterName(chapterID, info.Data.Chapters)

//images = retrieveMangaRockImages(pages.Data)

//if !found {
//log.Info("Chapter not found")
//chapter = chapterID
//}

//c.SetInfo(name, chapter, "")
//// c.SetImages(images)
//}
