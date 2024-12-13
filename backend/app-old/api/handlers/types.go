package handlers

// createPostRequestBody is the type of the "CreatePost"
// endpoint request body.
type createPostRequestBody struct {
	UserID string `json:"user_id,omitempty"`
	Text   string `json:"text"`
}

const (
	TimelineAccessed = "TimelineAccessed"
	PollingRequest   = "PollingRequest"
)
