package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	timeline "testGrpcTestClient/grpc/protogen/timeline"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	userID, _ := uuid.NewRandom()
	var text string
	for i := 1; i < 6; i++ {
		text = fmt.Sprintf("%d times post", i)
		createPosts(userID, text)
	}

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

	stream, err := client.GetPosts(ctx, req)
	if err != nil {
		log.Fatalf("could not get posts: %v", err)
	}
	userID2, _ := uuid.NewRandom()
	createPosts(userID2, "test")

	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Printf("stream end or error: %v", err)
			break
		}
		log.Printf("Response: %v", resp)
	}
}

type createPostRequestBody struct {
	UserID string `json:"user_id,omitempty"`
	Text   string `json:"text"`
}

func createPosts(userID uuid.UUID, text string) {
	// create request body
	requestBody := createPostRequestBody{
		UserID: userID.String(),
		Text:   text,
	}

	// encode body into json
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// send post request
	resp, err := http.Post("http://localhost:80/api/posts", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// display status code
	fmt.Printf("Response status: %s\n", resp.Status)
}
