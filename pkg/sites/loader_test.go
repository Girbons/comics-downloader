package sites

import (
	"errors"
	"testing"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/require"
)

type stubSite struct {
	issues []string
	comics map[string]*core.Comic
}

func (s *stubSite) Initialize(comic *core.Comic) error {
	if stub, ok := s.comics[comic.URLSource]; ok {
		*comic = *stub
		return nil
	}
	return errors.New("missing comic")
}

func (s *stubSite) GetInfo(url string) (string, string) {
	if stub, ok := s.comics[url]; ok {
		return stub.Name, stub.IssueNumber
	}
	return "", ""
}

func (s *stubSite) RetrieveIssueLinks() ([]string, error) {
	return s.issues, nil
}

func TestInitializeCollectionFiltersIssues(t *testing.T) {
	options := &config.Options{
		Source:       "test-source",
		Format:       "pdf",
		ImagesFormat: "png",
		IssuesRange:  "1-2",
		All:          true,
	}

	site := &stubSite{
		issues: []string{"url-1", "url-2", "url-3"},
		comics: map[string]*core.Comic{
			"url-1": {Name: "series", IssueNumber: "issue-1", URLSource: "url-1"},
			"url-2": {Name: "series", IssueNumber: "issue-2", URLSource: "url-2"},
			"url-3": {Name: "series", IssueNumber: "issue-3", URLSource: "url-3"},
		},
	}

	collection, err := initializeCollection(site.issues, options, site)
	require.NoError(t, err)
	require.Len(t, collection, 2)
	require.Equal(t, "issue-1", collection[0].IssueNumber)
	require.Equal(t, "issue-2", collection[1].IssueNumber)
}

func TestLoadComicFromSourceUnknown(t *testing.T) {
	options := &config.Options{Source: "unknown"}
	collection, err := LoadComicFromSource(options)
	require.Error(t, err)
	require.Empty(t, collection)
}

func TestNotInIssuesRange(t *testing.T) {
	testCases := []struct {
		issue string
		start float64
		end   float64
		skip  bool
	}{
		{"1", 1, 2, false},
		{"3", 1, 2, true},
		{"2.5", 2, 3, false},
		{"abc", 1, 2, true},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.skip, notInIssuesRange(tc.issue, tc.start, tc.end))
	}
}
