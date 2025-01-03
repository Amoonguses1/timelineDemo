package timeline

import (
	"grpc-test/protogen/post"
	"grpc-test/protogen/timeline"
	"log"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func BasicPosts() {
	p := post.Post{
		UserId:    "hoge",
		Id:        "piyo",
		Text:      "new test",
		CreatedAt: timestamppb.New(time.Now()),
	}

	log.Println(&p)
}

func BasicTimelineResponse() {
	p1 := post.Post{
		UserId:    "hoge",
		Id:        "piyo",
		Text:      "new test",
		CreatedAt: timestamppb.New(time.Now()),
	}

	p2 := post.Post{
		UserId:    "hoge",
		Id:        "piyo",
		Text:      "new test",
		CreatedAt: timestamppb.New(time.Now()),
	}

	resp := timeline.TimelineResponse{
		EventType: timeline.Event_POSTS_DELETED,
		Posts:     []*post.Post{&p1, &p2},
	}

	log.Println(&resp)
}
