package server

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
	"timelineDemo/grpc/protogen/image"
	"timelineDemo/grpc/protogen/post"
	timelinegrpc "timelineDemo/grpc/protogen/timeline"
	"timelineDemo/internal/app/usecases"
	"timelineDemo/internal/domain/entities"
	fileio "timelineDemo/internal/infrastructure/fileIO"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const baseDir = "./images/"

func NewGrpcServer(u usecases.GetUserAndFolloweePostsUsecaseInterface, mu *sync.Mutex, usersChan *map[uuid.UUID]chan entities.TimelineEvent, isBench bool) *GrpcServer {
	return &GrpcServer{u: u, mu: mu, usersChan: usersChan, isBench: isBench}
}

type GrpcServer struct {
	timelinegrpc.UnimplementedTimelineServiceServer

	u         usecases.GetUserAndFolloweePostsUsecaseInterface
	mu        *sync.Mutex
	usersChan *map[uuid.UUID]chan entities.TimelineEvent
	isBench   bool
}

func (s *GrpcServer) GetPosts(req *timelinegrpc.TimelineRequest, stream timelinegrpc.TimelineService_GetPostsServer) error {
	userID, err := uuid.Parse(req.GetId())
	if err != nil {
		err = status.Error(codes.InvalidArgument, "failed to parse user id")
		return err
	}
	if s.isBench {
		timestamp := time.Now().Format("15:04:05.000")
		fileio.WriteNewText("gRPCBenchLogs.txt", fmt.Sprintf("comes, %s, %s", userID.String()[:7], timestamp))
	}

	posts, err := s.u.GetUserAndFolloweePosts(userID)
	if err != nil {
		err = status.Error(codes.Internal, "failed to get posts")
		return err
	}

	s.mu.Lock()
	if _, exists := (*s.usersChan)[userID]; !exists {
		(*s.usersChan)[userID] = make(chan entities.TimelineEvent, 1)
	}
	userChan := (*s.usersChan)[userID]
	s.mu.Unlock()
	ctx := stream.Context()

	response, err := convertToTimelineResponse(entities.TimelineEvent{EventType: entities.TimelineAccessed, Posts: posts})
	if err != nil {
		err = status.Error(codes.Internal, "failed to convert posts")
		return err
	}
	if err = stream.Send(response); err != nil {
		return err
	}

	for {
		select {
		case event := <-userChan:
			log.Println("event comes in")
			response, err := convertToTimelineResponse(event)
			if err != nil {
				err = status.Error(codes.Internal, "failed to convert posts")
				return err
			}
			if s.isBench {
				timestamp := time.Now().Format("15:04:05.000")
				fileio.WriteNewText("gRPCBenchLogs.txt", fmt.Sprintf("send, %s, %s", userID.String()[:7], timestamp))
			}
			if err = stream.Send(response); err != nil {
				return err
			}
		case <-ctx.Done():
			s.mu.Lock()
			delete(*s.usersChan, userID)
			s.mu.Unlock()
			return nil
		case <-time.After(time.Second * 10):
			fmt.Println("timed out waiting for messages")
			return nil
		}
	}
}

func convertToTimelineResponse(event entities.TimelineEvent) (*timelinegrpc.TimelineResponse, error) {
	var eventType timelinegrpc.Event
	switch event.EventType {
	case entities.TimelineAccessed:
		eventType = timelinegrpc.Event_INITIAL_ACCESS
	case entities.PostCreated:
		eventType = timelinegrpc.Event_POST_CREATED
	case entities.PostDeleted:
		eventType = timelinegrpc.Event_POSTS_DELETED
	default:
		return nil, fmt.Errorf("unknown type")
	}

	var posts []*post.Post
	for _, p := range event.Posts {
		posts = append(posts, &post.Post{
			UserId:    p.UserID.String(),
			Id:        p.ID.String(),
			Text:      p.Text,
			CreatedAt: timestamppb.New(p.CreatedAt),
		})
	}

	return &timelinegrpc.TimelineResponse{
		EventType: eventType,
		Posts:     posts,
		ImagePath: event.ImagePath,
	}, nil
}

func (s *GrpcServer) GetImages(req *timelinegrpc.ImageRequest, stream timelinegrpc.TimelineService_GetImagesServer) error {
	userID, err := uuid.Parse(req.GetId())
	if err != nil {
		err = status.Error(codes.InvalidArgument, "failed to parse user id")
		return err
	}

	if s.isBench {
		timestamp := time.Now().Format("15:04:05.000")
		fileio.WriteNewText("gRPCBenchLogsImage.txt", fmt.Sprintf("comes, %s, %s", userID.String()[:7], timestamp))
	}

	fileNames := req.GetFileNames()
	buf := make([]byte, 1024)

	for _, fileName := range fileNames {
		filePath := filepath.Join(baseDir, filepath.Clean(fileName))
		log.Println(filePath)

		imgFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer imgFile.Close()

		if s.isBench {
			timestamp := time.Now().Format("15:04:05.000")
			fileio.WriteNewText("gRPCBenchLogsImage.txt", fmt.Sprintf("send, %s, %s", userID.String()[:7], timestamp))
		}

		for {
			pos, err := imgFile.Read(buf)
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return err
			}

			if err := stream.Send(&timelinegrpc.ImageResponse{
				FileName: fileName,
				Chunk: &image.Chunk{
					Data:     buf[:pos],
					Position: int64(pos),
				},
			}); err != nil {
				return err
			}
		}

	}

	return nil
}
