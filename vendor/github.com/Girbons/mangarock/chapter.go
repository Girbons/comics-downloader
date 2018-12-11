package mangarock

// Chapter contains the information related to a Manga chapter
type Chapter struct {
	CID       int    `json:"cid"`
	Name      string `json:"name"`
	OID       string `json:"oid"`
	Order     int    `json:"order"`
	UpdatedAt int    `json:"updated_at"`
}
