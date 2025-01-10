package timeline

import (
	"grpc-test/protogen/post"
	timelinegrpc "grpc-test/protogen/timeline"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{}
}

type GrpcServer struct {
	timelinegrpc.UnimplementedTimelineServiceServer
}

func (s *GrpcServer) GetPosts(req *timelinegrpc.TimelineRequest, stream timelinegrpc.TimelineService_GetPostsServer) error {
	post1 := &timelinegrpc.TimelineResponse{
		EventType: timelinegrpc.Event_POST_CREATED,
		Posts: []*post.Post{
			dummyPost(req.GetId()),
		},
	}
	if err := stream.Send(post1); err != nil {
		return err
	}

	post2 := &timelinegrpc.TimelineResponse{
		EventType: timelinegrpc.Event_POST_CREATED,
		Posts: []*post.Post{
			dummyPost(req.GetId()),
		},
	}
	if err := stream.Send(post2); err != nil {
		return err
	}

	return nil
}

func dummyPost(id string) *post.Post {
	return &post.Post{
		UserId:    id,
		Id:        "piyo",
		Text:      "new test",
		CreatedAt: timestamppb.New(time.Now()),
	}
}
