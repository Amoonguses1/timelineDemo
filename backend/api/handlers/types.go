package handlers

// createPostRequestBody is the type of the "CreatePost"
// endpoint request body.
type createPostRequestBody struct {
	UserID string `json:"user_id,omitempty"`
	Text   string `json:"text"`
}

type PollingEventType int

const (
	_ PollingEventType = iota
	TimelineAccessed
	PollingRequest
)

// longPollingTimelineRequestBody is the type of the "LongPollingTimeline"
// endpoint request body.
type longPollingTimelineRequestBody struct {
	PollingEventType `json:"polling_event_type"`
}
