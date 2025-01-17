package main

import (
	"context"
	"log"
	"time"

	timeline "testGrpcTestClient/grpc/protogen/timeline"

	"google.golang.org/grpc"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.NewClient("localhost:8081", opts...)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	client := timeline.NewTimelineServiceClient(conn)

	req := &timeline.TimelineRequest{Id: "12345"}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.GetPosts(ctx, req)
	if err != nil {
		log.Fatalf("could not get posts: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Printf("stream end or error: %v", err)
			break
		}
		log.Printf("Response: %v", resp)
	}
}
