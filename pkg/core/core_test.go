package core

import (
	"bytes"
	"encoding/base64"
	"image"
	_ "image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/pkg/config"
	httpclient "github.com/Girbons/comics-downloader/pkg/http"
	"github.com/stretchr/testify/require"
)

var samplePNG = func() []byte {
	data, _ := base64.StdEncoding.DecodeString(
		"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAADElEQVR4nGP4//8/AAX+Av4N70a4AAAAAElFTkSuQmCC",
	)
	return data
}()

func newImageServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(samplePNG)
	}))
}

func newTestOptions(t *testing.T, server *httptest.Server) *config.Options {
	t.Helper()

	client := httpclient.NewComicClient(
		httpclient.WithHTTPClient(server.Client()),
		httpclient.WithRetry(0, 0),
	)

	return &config.Options{
		OutputFolder:      t.TempDir(),
		CreateDefaultPath: true,
		Debug:             false,
		Logger:            logger.NewLogger(false, nil),
		Client:            client,
		IssueFolderName:   "issue-",
	}
}

func buildLinks(server *httptest.Server, count int) []string {
	links := make([]string, count)
	for i := 0; i < count; i++ {
		links[i] = server.URL + "/img.png"
	}
	return links
}

func TestDownloadImagesCreatesFiles(t *testing.T) {
	server := newImageServer()
	defer server.Close()

	opts := newTestOptions(t, server)

	comic := &Comic{
		Name:         "foo",
		Source:       "test-source",
		IssueNumber:  "1",
		ImagesFormat: "png",
		Links:        buildLinks(server, 3),
	}

	result, err := comic.DownloadImages(opts)
	require.NoError(t, err)
	require.Len(t, result.FilePaths, 3)

	for _, file := range result.FilePaths {
		data, err := readFile(file)
		require.NoError(t, err)
		require.NotEmpty(t, data)
		_, format, err := image.Decode(bytes.NewReader(data))
		require.NoError(t, err)
		require.Equal(t, "png", format)
	}
}

func TestMakeComicPDF(t *testing.T) {
	server := newImageServer()
	defer server.Close()

	opts := newTestOptions(t, server)

	comic := &Comic{
		Name:         "foo",
		Source:       "test-source",
		IssueNumber:  "1",
		Format:       PDF,
		ImagesFormat: "png",
		Links:        buildLinks(server, 2),
	}

	require.NoError(t, comic.MakeComic(opts))

	output := filepath.Join(opts.OutputFolder, "comics", comic.Source, comic.Name, "foo-1.pdf")
	require.FileExists(t, output)
}

func TestMakeComicEPUB(t *testing.T) {
	server := newImageServer()
	defer server.Close()

	opts := newTestOptions(t, server)

	comic := &Comic{
		Name:         "bar",
		Source:       "test-source",
		IssueNumber:  "42",
		Author:       "Author",
		Format:       EPUB,
		ImagesFormat: "png",
		Links:        buildLinks(server, 2),
	}

	require.NoError(t, comic.MakeComic(opts))

	output := filepath.Join(opts.OutputFolder, "comics", comic.Source, comic.Name, "bar-42.epub")
	require.FileExists(t, output)
}

func TestMakeComicCBZ(t *testing.T) {
	server := newImageServer()
	defer server.Close()

	opts := newTestOptions(t, server)

	comic := &Comic{
		Name:         "baz",
		Source:       "test-source",
		IssueNumber:  "7",
		Format:       CBZ,
		ImagesFormat: "png",
		Links:        buildLinks(server, 2),
	}

	require.NoError(t, comic.MakeComic(opts))

	output := filepath.Join(opts.OutputFolder, "comics", comic.Source, comic.Name, "baz-7.cbz")
	require.FileExists(t, output)
}

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}
