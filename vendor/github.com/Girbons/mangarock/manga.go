package mangarock

// Manga Struct contains all the informations returned by mangarock about a manga
type Manga struct {
	Author         string            `json:"author"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Thumbnail      string            `json:"thumbnail"`
	Cover          string            `json:"cover"`
	OID            string            `json:"oid"`
	MsID           int               `json:"msid"`
	Authors        []*Author         `json:"authors"`
	Alias          []string          `json:"alias"`
	Artworks       []string          `json:"artworks"`
	Characters     []*Character      `json:"characters"`
	Chapters       []*Chapter        `json:"chapters"`
	Categories     []int             `json:"categories"`
	Completed      bool              `json:"completed"`
	RichCategories []*RichCategory   `json:"rich_categories"`
	MrsSeries      int               `json:"mrs_series"`
	Extra          map[string]string `json:"extra"`
	LastUpdate     int               `json:"last_update"`
	Direction      int               `json:"direction"`
	Rank           int               `json:"rank"`
	Mid            int               `json:"mid"`
	TotalChapters  int               `json:"total_chapters"`
	Removed        bool              `json:"removed"`
}
