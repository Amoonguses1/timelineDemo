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
	"google.golang.org/protobuf/proto"
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
	// measure request size
	reqBytes, err := proto.Marshal(req)
	if err != nil {
		log.Printf("Failed to marshal request: %v", err)
		return
	}
	reqSize := len(reqBytes)
	fmt.Println("reqsize", reqSize)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	timestamp := time.Now().UTC().Format("15:04:05.000")
	fileio.WriteNewText("gRPCBenchLogs.txt", fmt.Sprintf("send, %s, %s", userID.String()[:7], timestamp))
	stream, err := client.GetPosts(ctx, req)
	if err != nil {
		log.Fatalf("could not get posts: %v", err)
	}

	// notification for connection established
	connected <- struct{}{}

	var totalRespSize int64

	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Printf("stream end or error: %v", err)
			break
		}
		if resp.EventType != timeline.Event_INITIAL_ACCESS {
			// measure response size
			respBytes, err := proto.Marshal(resp)
			if err != nil {
				log.Printf("Failed to marshal response: %v", err)
				continue
			}
			respSize := len(respBytes)
			totalRespSize += int64(respSize)

			timestamp := time.Now().UTC().Format("15:04:05.000")
			fileio.WriteNewText("gRPCBenchLogs.txt", fmt.Sprintf("comes, %s, %s", userID.String()[:7], timestamp))
		}
		if end(resp.Posts) {
			break
		}
		// log.Println("Response:", resp.Posts)
	}
	fmt.Println("total resp: ", totalRespSize)
}

func end(posts []*post.Post) bool {
	for _, post := range posts {
		if strings.Contains(post.Text, "end") {
			return true
		}
	}
	return false
}
