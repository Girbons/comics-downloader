package core

import (
	"testing"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestGetDescriptionForLanguage_MatchingLanguage(t *testing.T) {
	opts := &config.Options{Country: "en"}
	c := &ComicIssue{
		SeriesMetadata: &SeriesMetadata{
			Description: map[string]string{
				"en": "English description",
				"jp": "Japanese description",
			},
		},
	}

	desc := c.getDescriptionForLanguage(opts, "en")
	require.Equal(t, "English description", desc)
}

func TestGetDescriptionForLanguage_FallbackSingleEntry(t *testing.T) {
	opts := &config.Options{Country: "fr"}
	c := &ComicIssue{
		SeriesMetadata: &SeriesMetadata{
			Description: map[string]string{
				"es": "Spanish description",
			},
		},
	}

	desc := c.getDescriptionForLanguage(opts, "fr")
	require.Equal(t, "Spanish description", desc)
}

func TestGetDescriptionForLanguage_Empty(t *testing.T) {
	opts := &config.Options{Country: "en"}
	c := &ComicIssue{
		SeriesMetadata: &SeriesMetadata{
			Description: map[string]string{},
		},
	}

	desc := c.getDescriptionForLanguage(opts, "en")
	require.Equal(t, "", desc)
}
