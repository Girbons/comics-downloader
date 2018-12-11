package mangarock

// Character contains the information related to a Manga character
type Character struct {
	Name      string `json:"name"`
	OID       string `json:"oid"`
	Thumbnail string `json:"thumbnail"`
}
