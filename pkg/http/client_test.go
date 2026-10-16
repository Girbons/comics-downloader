package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareRequestMangakakalot(t *testing.T) {
	cc := NewComicClient()
	link := "http://mangakakalot.com"
	source := "mangakakalot.com"
	req, err := cc.PrepareRequest(link, source)

	assert.Equal(t, req.Header["Referer"], []string{link})
	assert.Nil(t, err)
}

func TestPrepareRequest(t *testing.T) {
	cc := NewComicClient()
	link := "http://foo.com"
	req, err := cc.PrepareRequest(link, "foo")

	assert.Equal(t, len(req.Header["Referer"]), 0)
	assert.Nil(t, err)
}

func TestGetReferer(t *testing.T) {
	cases := []struct {
		source      string
		wantReferer bool
	}{
		{"www.mangatown.com", true},
		{"mangatown.com", true},
		{"mangakakalot.com", true},
		{"manganato.com", true},
		{"readcomiconline.li", false},
		{"foo", false},
	}

	for _, tc := range cases {
		t.Run(tc.source, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				want := ""
				if tc.wantReferer {
					want = "http://" + r.Host + r.URL.RequestURI()
				}
				if got := r.Header.Get("Referer"); got != want {
					t.Errorf("Referer = %q, want %q", got, want)
					w.WriteHeader(http.StatusForbidden)
					return
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			// Image requests can use a different host from the comic source.
			response, err := NewComicClient().Get(server.URL+"/images/page.jpg", tc.source)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			assert.Equal(t, http.StatusOK, response.StatusCode)
		})
	}
}

func TestPrepareRequestInvalidURL(t *testing.T) {
	for _, source := range []string{"www.mangatown.com", "mangatown.com"} {
		t.Run(source, func(t *testing.T) {
			req, err := NewComicClient().PrepareRequest("http://example.com/%zz", source)
			assert.Error(t, err)
			assert.Nil(t, req)
		})
	}
}

func TestGetInvalidURL(t *testing.T) {
	for _, source := range []string{"www.mangatown.com", "mangatown.com"} {
		t.Run(source, func(t *testing.T) {
			response, err := NewComicClient().Get("http://example.com/%zz", source)
			assert.Error(t, err)
			assert.Nil(t, response)
		})
	}
}
