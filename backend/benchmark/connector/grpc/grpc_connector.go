package grpcconnector

import (
	"benchmark/connector"
	"benchmark/connector/grpc/protogen/post"
	"benchmark/connector/grpc/protogen/timeline"
	fileio "benchmark/fileIO"
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCConnector struct {
}

func NewGRPCConnector() connector.ConnectorInterface {
	return &GRPCConnector{}
}

func (c *GRPCConnector) Connect(userID uuid.UUID, wg *sync.WaitGroup, connected chan<- struct{}) {
	defer wg.Done()

	// set up connection
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient("localhost:8081", opts...)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	client := timeline.NewTimelineServiceClient(conn)

	req := &timeline.TimelineRequest{Id: userID.String()}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fileio.WriteNewText("gRPCBenchLogs.txt", fmt.Sprintf("request send\n%s: %v\n", userID.String()[:7], time.Now()))
	stream, err := client.GetPosts(ctx, req)
	if err != nil {
		log.Fatalf("could not get posts: %v", err)
	}

	// notification for connection established
	connected <- struct{}{}

	for {
		resp, err := stream.Recv()
		fileio.WriteNewText("gRPCBenchLogs.txt", fmt.Sprintf("response received\n%s: %v\n", userID.String()[:7], time.Now()))
		if err != nil {
			log.Printf("stream end or error: %v", err)
			break
		}
		if end(resp.Posts) {
			break
		}
		log.Println("Response:", resp.Posts)
	}
}

func end(posts []*post.Post) bool {
	for _, post := range posts {
		if strings.Contains(post.Text, "end") {
			return true
		}
	}
	return false
}
