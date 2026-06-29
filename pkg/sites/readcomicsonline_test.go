package sites

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/require"
)

const (
	rcouIssuePath = "/comic/my-comic/2"
	rcouListPath  = "/comic/my-comic"
)

// newReadComicsOnlineServer builds a test HTTP server that mimics the
// structure returned by readcomicsonline.ru.
func newReadComicsOnlineServer() *httptest.Server {
	// Issue page: images in data-src attributes (lazy-load pattern).
	issueHTML := `
		<html><body>
			<ul class="chapters">
				<li><a href="/comic/my-comic/2">Issue 2</a></li>
				<li><a href="/comic/my-comic/1">Issue 1</a></li>
			</ul>
			<div id="all">
				<img class="img-responsive" src="data:image/gif;base64,R0=" data-src=" https://readcomicsonline.ru/uploads/manga/my-comic/chapters/2/01.jpg "/>
				<img class="img-responsive" src="data:image/gif;base64,R0=" data-src=" https://readcomicsonline.ru/uploads/manga/my-comic/chapters/2/02.jpg "/>
			</div>
		</body></html>`

	// Issue page with no data-src: fallback to var pages JS variable.
	issueNoDataSrcHTML := `
		<html><body>
			<ul class="chapters">
				<li><a href="/comic/my-comic/2">Issue 2</a></li>
				<li><a href="/comic/my-comic/1">Issue 1</a></li>
			</ul>
			<div id="all">
				<img class="img-responsive" src="data:image/gif;base64,R0="/>
			</div>
			<script>
				var pages = [{"page_image":"01.jpg","page_slug":1,"external":0},{"page_image":"02.jpg","page_slug":2,"external":0}];
			</script>
		</body></html>`

	// Comic listing page.
	listHTML := `
		<html><body>
			<ul class="chapters">
				<li><h5><a href="/comic/my-comic/2">Issue 2</a></h5></li>
				<li><h5><a href="/comic/my-comic/1">Issue 1</a></h5></li>
			</ul>
		</body></html>`

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case rcouIssuePath:
			_, _ = fmt.Fprint(w, issueHTML)
		case "/comic/my-comic/1":
			_, _ = fmt.Fprint(w, issueNoDataSrcHTML)
		case rcouListPath:
			_, _ = fmt.Fprint(w, listHTML)
		default:
			http.NotFound(w, r)
		}
	}))
}

func newRCOUScraper(server *httptest.Server, url string) *ReadComicsOnline {
	opts := &config.Options{
		URL:    url,
		Logger: logger.NewLogger(false, nil),
	}
	s := NewReadComicsOnline(opts)
	s.baseURL = server.URL // redirect absolute URL construction to test server
	return s
}

// TestReadComicsOnline_Initialize_DataSrc verifies the primary (data-src) path.
func TestReadComicsOnline_Initialize_DataSrc(t *testing.T) {
	server := newReadComicsOnlineServer()
	defer server.Close()

	s := newRCOUScraper(server, server.URL+rcouIssuePath)
	comic := &core.Comic{URLSource: server.URL + rcouIssuePath}
	require.NoError(t, s.Initialize(comic))
	require.Equal(t, []string{
		"https://readcomicsonline.ru/uploads/manga/my-comic/chapters/2/01.jpg",
		"https://readcomicsonline.ru/uploads/manga/my-comic/chapters/2/02.jpg",
	}, comic.Links)
}

// TestReadComicsOnline_Initialize_PagesVarFallback verifies the JS-pages fallback.
func TestReadComicsOnline_Initialize_PagesVarFallback(t *testing.T) {
	server := newReadComicsOnlineServer()
	defer server.Close()

	s := newRCOUScraper(server, server.URL+"/comic/my-comic/1")
	comic := &core.Comic{URLSource: server.URL + "/comic/my-comic/1"}
	require.NoError(t, s.Initialize(comic))
	require.Equal(t, []string{
		server.URL + "/uploads/manga/my-comic/chapters/1/01.jpg",
		server.URL + "/uploads/manga/my-comic/chapters/1/02.jpg",
	}, comic.Links)
}

// TestReadComicsOnline_RetrieveIssueLinks_Single verifies a single-issue URL.
func TestReadComicsOnline_RetrieveIssueLinks_Single(t *testing.T) {
	server := newReadComicsOnlineServer()
	defer server.Close()

	s := newRCOUScraper(server, server.URL+rcouIssuePath)
	issues, err := s.RetrieveIssueLinks()
	require.NoError(t, err)
	require.Equal(t, []string{server.URL + rcouIssuePath}, issues)
}

// TestReadComicsOnline_RetrieveIssueLinks_All verifies --all from a listing page.
func TestReadComicsOnline_RetrieveIssueLinks_All(t *testing.T) {
	server := newReadComicsOnlineServer()
	defer server.Close()

	s := newRCOUScraper(server, server.URL+rcouListPath)
	issues, err := s.RetrieveIssueLinks()
	require.NoError(t, err)
	require.Equal(t, []string{
		server.URL + "/comic/my-comic/2",
		server.URL + "/comic/my-comic/1",
	}, issues)
}

// TestReadComicsOnline_RetrieveIssueLinks_Last verifies --last from a listing page.
func TestReadComicsOnline_RetrieveIssueLinks_Last(t *testing.T) {
	server := newReadComicsOnlineServer()
	defer server.Close()

	opts := &config.Options{
		URL:    server.URL + rcouListPath,
		Last:   true,
		Logger: logger.NewLogger(false, nil),
	}
	s := NewReadComicsOnline(opts)
	s.baseURL = server.URL
	issues, err := s.RetrieveIssueLinks()
	require.NoError(t, err)
	require.Equal(t, []string{server.URL + "/comic/my-comic/2"}, issues)
}

// TestReadComicsOnline_RetrieveIssueLinks_AllFromIssue verifies --all when starting
// from a specific issue URL (should fetch the parent listing).
func TestReadComicsOnline_RetrieveIssueLinks_AllFromIssue(t *testing.T) {
	server := newReadComicsOnlineServer()
	defer server.Close()

	opts := &config.Options{
		URL:    server.URL + rcouIssuePath,
		All:    true,
		Logger: logger.NewLogger(false, nil),
	}
	s := NewReadComicsOnline(opts)
	s.baseURL = server.URL
	issues, err := s.RetrieveIssueLinks()
	require.NoError(t, err)
	require.Equal(t, []string{
		server.URL + "/comic/my-comic/2",
		server.URL + "/comic/my-comic/1",
	}, issues)
}

// TestReadComicsOnline_GetInfo validates name/issue extraction.
func TestReadComicsOnline_GetInfo(t *testing.T) {
	s := NewReadComicsOnline(&config.Options{})
	cases := []struct {
		url, name, issue string
	}{
		{"https://readcomicsonline.ru/comic/briar-nights-terror-2025/1", "briar-nights-terror-2025", "1"},
		{"https://readcomicsonline.ru/comic/some-comic/42", "some-comic", "42"},
	}
	for _, tc := range cases {
		name, issue := s.GetInfo(tc.url)
		require.Equal(t, tc.name, name, "url=%s", tc.url)
		require.Equal(t, tc.issue, issue, "url=%s", tc.url)
	}
}
