package entities

const (
	TimelineAccessed = "TimelineAccessed"
	PostCreated      = "PostCreated"
	PostDeleted      = "PostDeleted"
)

type TimelineEvent struct {
	EventType string
	Posts     []*Post
	ImagePath string
}
