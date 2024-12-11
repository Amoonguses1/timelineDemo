package entities

type EventType int

const (
	TimelineAccessed = "TimelineAccessed"
	PostCreated      = "PostCreated"
	PostDeleted      = "PostDeleted"
)

type TimelineEvent struct {
	EventType string
	Posts     []*Post
}
