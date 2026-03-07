package core

import (
	"strings"
	"time"

	"github.com/Girbons/comics-downloader/pkg/config"
)

type AgeRating string

// Based on https://anansi-project.github.io/docs/comicinfo/schemas/v2.1
const (
	AgeRatingUnrated        AgeRating = "Unrated"
	AgeRatingAO18           AgeRating = "Adults Only 18+"
	AgeRatingEarlyChildhood AgeRating = "Early Childhood"
	AgeRatingEveryone       AgeRating = "Everyone"
	AgeRatingEveryone10     AgeRating = "Everyone 10+"
	AgeRatingG              AgeRating = "G"
	AgeRatingKidsToAdults   AgeRating = "Kids to Adults"
	AgeRatingM              AgeRating = "M"
	AgeRatingMA15           AgeRating = "MA15+"
	AgeRatingMature         AgeRating = "Mature 17+"
	AgeRatingPG             AgeRating = "PG"
	AgeRatingR18            AgeRating = "R18+"
	AgeRatingRatingPending  AgeRating = "Rating Pending"
	AgeRatingTeen           AgeRating = "Teen"
	AgeRatingX18            AgeRating = "X18+"
)

type ComicSource struct {
	Name string
	URL  string // URL of the comic/manga issue
}

type CreatorRole string

type SeriesCreator struct {
	Name string
	Role CreatorRole
}

type SeriesMetadata struct {
	Title          string            // Series title, should be in the native language of the comic/manga when possible
	LocalizedTitle map[string]string // Map of language code to localized title, e.g. {"en": "One Piece", "jp": "ワンピース"}
	Description    map[string]string // Map of language code to description, e.g. {"en": "A story about pirates...", "jp": "海賊の物語..."}
	Creators       []SeriesCreator

	IsManga *bool // True if it's a manga, false if it's a comic
	IsRTL   *bool // True if the comic/manga is read right-to-left, false if left-to-right. Only relevant for manga, but some comics may also be RTL.

	CommunityRating *float64 // Average rating from the community, 0-5
	AgeRating       *AgeRating
	Tags            []string // ninja or school life
	Genres          []string // eg Science-Fiction or Shonen

	WebLinks []string // Official website, social media, etc.
}

// ComicIssue struct contains all the informations about a comic
type ComicIssue struct {
	Name string // Issue name/title

	IssueNumber string
	Volume      *string
	LanguageISO *string // IETF language tag
	ReleaseDate *time.Time

	ImageLinks   []string
	OutputFormat ComicOutputFormat
	ImagesFormat string

	Source         *ComicSource
	SeriesMetadata *SeriesMetadata
}

const (
	CreatorRoleUnknown     CreatorRole = "Unknown"
	CreatorRoleWriter      CreatorRole = "Writer"
	CreatorRolePenciller   CreatorRole = "Penciller"
	CreatorRoleInker       CreatorRole = "Inker"
	CreatorRoleColorist    CreatorRole = "Colorist"
	CreatorRoleLetterer    CreatorRole = "Letterer"
	CreatorRoleCoverArtist CreatorRole = "CoverArtist"
	CreatorRoleEditor      CreatorRole = "Editor"
	CreatorRoleTranslator  CreatorRole = "Translator"
)

// func SourceAuthorRoleToSeriesAuthorRole(sourceRole string) string {
// 	sourceRole = strings.ToLower(strings.TrimSpace(sourceRole))
// 	switch sourceRole {
// 	case "writer", "author":
// 		return CreatorRoleWriter
// 	case "penciller":
// 		return CreatorRolePenciller
// 	case "inker":
// 		return CreatorRoleInker
// 	case "colorist":
// 		return CreatorRoleColorist
// 	case "letterer":
// 		return CreatorRoleLetterer
// 	case "cover artist", "coverartist", "cover-artist":
// 		return CreatorRoleCoverArtist
// 	case "editor":
// 		return CreatorRoleEditor
// 	case "translator":
// 		return CreatorRoleTranslator
// 	default:
// 		return CreatorRoleUnknown
// 	}
// }

func (c *ComicIssue) getLocalizedDescription(options *config.Options) string {
	if len(c.SeriesMetadata.Description) <= 0 {
		return ""
	}

	for lang, desc := range c.SeriesMetadata.Description {
		if options.Country == "" || options.Country == strings.ToLower(lang) {
			return desc
		}
	}

	// if no description matching the country option is found, set the description to the first one in the map
	for _, desc := range c.SeriesMetadata.Description {
		return desc
	}

	return ""
}

func (c *ComicIssue) getLocalizedTitle(options *config.Options) string {
	if len(c.SeriesMetadata.LocalizedTitle) <= 0 {
		// fmt.Println("No localized titles")
		return c.SeriesMetadata.Title
	}

	// for key, value := range c.SeriesMetadata.LocalizedTitle {
	// 	options.Logger.Infof("%s: %s\n", key, value)
	// }

	for lang, title := range c.SeriesMetadata.LocalizedTitle {
		if options.Country == "" || options.Country == strings.ToLower(lang) {
			// options.Logger.Infof("Found desired title \"%s\" in the lang %s\n", title, lang)
			return title
		}
	}

	// if no title matching the country option is found, set the title to the first one in the map
	for _, title := range c.SeriesMetadata.LocalizedTitle {
		return title
	}

	return c.SeriesMetadata.Title
}
