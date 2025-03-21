package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

	imageReq := &timeline.ImageRequest{FileNames: []string{"Go-Logo_Aqua.png"}}
	imageCtx, imageCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer imageCancel()
	imageStream, err := client.GetImages(imageCtx, imageReq)
	if err != nil {
		log.Fatalf("Error calling GetImages: %v", err)
	}

	outFile, err := os.Create("received_Go-Logo_Aqua.png")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	for {
		resp, err := imageStream.Recv()
		if err != nil {
			log.Printf("Stream closed: %v", err)
			break
		}

		_, err = outFile.Write(resp.Chunk.Data)
		if err != nil {
			log.Fatalf("Failed to write chunk: %v", err)
		}
		log.Printf("Write chunk: %d\n", resp.Chunk.Position)
	}

	log.Println("File received successfully: received_Go-Logo_Aqua.png")
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
