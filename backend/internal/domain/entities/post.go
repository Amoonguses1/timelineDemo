package entities

// Post represents an entry of `posts` table.
type Post struct {
	UserID string `json:"user_id"`
	Text   string `json:"text"`
}
