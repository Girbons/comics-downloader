package sites

import (
	"errors"
	"testing"

	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubSite struct {
	issues []string
	comics map[string]*core.ComicIssue
}

func (s *stubSite) Initialize(comic *core.ComicIssue) error {
	if stub, ok := s.comics[comic.Source.URL]; ok {
		*comic = *stub
		return nil
	}
	return errors.New("missing comic")
}

func (s *stubSite) GetInfo(url string) (string, string, error) {
	if stub, ok := s.comics[url]; ok {
		return stub.Name, stub.IssueNumber, nil
	}
	return "", "", errors.New("missing comic")
}

func (s *stubSite) RetrieveIssueLinks() ([]string, error) {
	return s.issues, nil
}

func TestInitializeCollectionFiltersIssues(t *testing.T) {
	options := &config.Options{
		SourceName:         "test-source",
		OutputFormat:       "pdf",
		OutputImagesFormat: "png",
		IssuesRange:        "1-2",
		All:                true,
		Logger:             logger.NewLogger(false, nil),
	}

	site := &stubSite{
		issues: []string{"url-1", "url-2", "url-3"},
		comics: map[string]*core.ComicIssue{
			"url-1": {Name: "series", IssueNumber: "issue-1", Source: &core.ComicSource{Name: "test-source", URL: "url-1"}},
			"url-2": {Name: "series", IssueNumber: "issue-2", Source: &core.ComicSource{Name: "test-source", URL: "url-2"}},
			"url-3": {Name: "series", IssueNumber: "issue-3", Source: &core.ComicSource{Name: "test-source", URL: "url-3"}},
		},
	}

	collection, err := initializeCollection(site.issues, options, site)
	require.NoError(t, err)
	require.Len(t, collection, 2)
	require.Equal(t, "issue-1", collection[0].IssueNumber)
	require.Equal(t, "issue-2", collection[1].IssueNumber)
}

func TestLoadComicFromSourceUnknown(t *testing.T) {
	options := &config.Options{SourceName: "unknown", Logger: logger.NewLogger(false, nil)}
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

func TestVolumeIssuesRange(t *testing.T) {
	tt := []struct {
		name        string
		input       string
		start       float64
		end         float64
		returnValue bool
	}{
		// Volume 4, Issue 78-99 range tests (user specifies: -range=4.78-4.99)
		{"v4-078 in range", "v4-078-2016", 4.78, 4.99, false},
		{"v4-099 in range", "v4-099-2016", 4.78, 4.99, false},
		{"v4-077 out of range (too low)", "v4-077-2016", 4.78, 4.99, true},
		{"v4-100 out of range (too high)", "v4-100-2016", 4.78, 4.99, true},
		{"v3-078 wrong volume", "v3-078-2016", 4.78, 4.99, true},
		{"v5-078 wrong volume", "v5-078-2016", 4.78, 4.99, true},

		// Volume 2, Issue 1-50 range tests (user specifies: -range=2.01-2.50)
		{"v2-001 in range", "v2-001-1989", 2.01, 2.50, false},
		{"v2-025 in range", "v2-025-1989", 2.01, 2.50, false},
		{"v2-050 in range", "v2-050-1989", 2.01, 2.50, false},
		{"v2-051 out of range", "v2-051-1989", 2.01, 2.50, true},

		// Edge cases with different formats
		{"v4-078 without year", "v4-078", 4.78, 4.99, false},
		{"v4_078 with underscore", "v4_078", 4.78, 4.99, false},
		{"v10-005 two-digit volume", "v10-005", 10.05, 10.10, false},

		// Backwards compatibility with simple numeric issues
		{"078 simple format", "078", 78, 99, false},
		{"100 simple format out of range", "100", 78, 99, true},
		{"issue-1 with prefix", "issue-1", 1, 3, false},
		{"issue-5 with prefix out of range", "issue-5", 1, 3, true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.returnValue, notInIssuesRange(tc.input, tc.start, tc.end),
				"Issue %s with range %.2f-%.2f", tc.input, tc.start, tc.end)
		})
	}
}

func TestExtractIssueNumberForRange(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		expected float64
	}{
		// Volume and issue format
		{"v4-078-2016", "v4-078-2016", 4.78},
		{"v4-099-2016", "v4-099-2016", 4.99},
		{"v2-075-1989", "v2-075-1989", 2.75},
		{"v4-078 no year", "v4-078", 4.78},
		{"v4_078 underscore", "v4_078", 4.78},
		{"v10-005 two-digit volume", "v10-005", 10.05},

		// Simple numeric format
		{"078", "078", 78},
		{"99", "99", 99},
		{"1", "1", 1},

		// Decimal format
		{"20.5", "20.5", 20.5},
		{"3.14", "3.14", 3.14},

		// With prefixes
		{"issue-1", "issue-1", 1},
		{"issue-123", "issue-123", 123},

		// Edge cases
		{"empty", "", 0},
		{"no numbers", "abc", 0},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, extractIssueNumberForRange(tc.input),
				"Issue %s should extract to %.2f", tc.input, tc.expected)
		})
	}
}

func TestLoadComicFromSourceWithRegistry(t *testing.T) {
	// Save original registry to restore it after the test
	originalRegistry := make(map[string]SupportedSite)
	for k, v := range SupportedSites {
		originalRegistry[k] = v
	}
	defer func() {
		SupportedSites = originalRegistry
	}()

	// Create a test site implementation
	testSite := &stubSite{
		issues: []string{"url-1", "url-2"},
		comics: map[string]*core.ComicIssue{
			"url-1": {Name: "test-series", IssueNumber: "1", Source: &core.ComicSource{Name: "test-site", URL: "url-1"}},
			"url-2": {Name: "test-series", IssueNumber: "2", Source: &core.ComicSource{Name: "test-site", URL: "url-2"}},
		},
	}

	// Register the test site in the registry
	SupportedSites["test-site"] = SupportedSite{
		IsEnabled: true,
		Loader: func(opts *config.Options) BaseSite {
			return testSite
		},
	}

	options := &config.Options{
		SourceName:         "test-site",
		URL:                "http://test-site.com",
		OutputFormat:       "pdf",
		OutputImagesFormat: "png",
		Logger:             logger.NewLogger(false, nil),
	}

	collection, err := LoadComicFromSource(options)
	require.NoError(t, err)
	require.Len(t, collection, 2)
	assert.Equal(t, "test-series", collection[0].Name)
	assert.Equal(t, "1", collection[0].IssueNumber)
	assert.Equal(t, "test-series", collection[1].Name)
	assert.Equal(t, "2", collection[1].IssueNumber)
}

func TestLoadComicFromSourceDisabledSite(t *testing.T) {
	// Save original registry
	originalRegistry := make(map[string]SupportedSite)
	for k, v := range SupportedSites {
		originalRegistry[k] = v
	}
	defer func() {
		SupportedSites = originalRegistry
	}()

	// Register a disabled test site
	testSite := &stubSite{
		issues: []string{"url-1"},
		comics: map[string]*core.ComicIssue{
			"url-1": {Name: "test-series", IssueNumber: "1", Source: &core.ComicSource{Name: "disabled-test-site", URL: "url-1"}},
		},
	}

	SupportedSites["disabled-test-site"] = SupportedSite{
		IsEnabled: false,
		Loader: func(opts *config.Options) BaseSite {
			return testSite
		},
	}

	options := &config.Options{
		SourceName: "disabled-test-site",
		URL:        "http://disabled-test-site.com",
		Logger:     logger.NewLogger(false, nil),
	}

	collection, err := LoadComicFromSource(options)
	require.Error(t, err)
	require.Empty(t, collection)
	assert.Contains(t, err.Error(), "disabled")
}

func TestLoadComicFromSourcePartialMatch(t *testing.T) {
	// Save original registry
	originalRegistry := make(map[string]SupportedSite)
	for k, v := range SupportedSites {
		originalRegistry[k] = v
	}
	defer func() {
		SupportedSites = originalRegistry
	}()

	testSite := &stubSite{
		issues: []string{"url-1"},
		comics: map[string]*core.ComicIssue{
			"url-1": {Name: "my-comic", IssueNumber: "42", Source: &core.ComicSource{Name: "mysite.com", URL: "url-1"}},
		},
	}

	SupportedSites["mysite"] = SupportedSite{
		IsEnabled: true,
		Loader: func(opts *config.Options) BaseSite {
			return testSite
		},
	}

	// Test with full domain name to verify partial matching works
	options := &config.Options{
		SourceName:         "mysite.com",
		URL:                "http://mysite.com/comic",
		OutputFormat:       "pdf",
		OutputImagesFormat: "png",
		Logger:             logger.NewLogger(false, nil),
	}

	collection, err := LoadComicFromSource(options)
	require.NoError(t, err)
	require.Len(t, collection, 1)
	assert.Equal(t, "my-comic", collection[0].Name)
	assert.Equal(t, "42", collection[0].IssueNumber)
}

func TestLoadComicFromSourceUnsupportedSite(t *testing.T) {
	// Save original registry
	originalRegistry := make(map[string]SupportedSite)
	for k, v := range SupportedSites {
		originalRegistry[k] = v
	}
	defer func() {
		SupportedSites = originalRegistry
	}()

	// Clear registry to ensure no sites are registered
	SupportedSites = make(map[string]SupportedSite)

	options := &config.Options{
		SourceName: "unsupported-site.com",
		URL:        "http://unsupported-site.com/comic",
		Logger:     logger.NewLogger(false, nil),
	}

	collection, err := LoadComicFromSource(options)
	require.Error(t, err)
	require.Empty(t, collection)
	assert.Contains(t, err.Error(), "unknown")
}
