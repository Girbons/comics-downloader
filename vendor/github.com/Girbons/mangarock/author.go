package mangarock

// Author contains the info related to a manga author
type Author struct {
	Name      string `json:"name"`
	OID       string `json:"oid"`
	Role      string `json:"role"`
	Thumbnail string `json:"thumbnail"`
}
