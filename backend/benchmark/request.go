package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type CreatePostRequestBody struct {
	UserID string `json:"user_id,omitempty"`
	Text   string `json:"text"`
}

func EndRequest() {
	createPost(uuid.New(), "end")
}

func BenchRequests(numPosters, postsPerUser int) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := 0; i < numPosters; i++ {
		wg.Add(1)
		go sendRequestsBySingleUser(ctx, 1000, postsPerUser, uuid.New(), &wg)
	}

	wg.Wait()
}

func sendRequestsBySingleUser(ctx context.Context, interval, totalRequests int, userID uuid.UUID, wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(interval))
	defer ticker.Stop()
	defer wg.Done()

	for i := 0; i < totalRequests; i++ {
		select {
		case <-ticker.C:
			// log.Println("Posted")
			err := createPost(userID, fmt.Sprintf("text no. %d", i))
			if err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func createPost(userID uuid.UUID, text string) error {
	// create request body
	requestBody := CreatePostRequestBody{
		UserID: userID.String(),
		Text:   text,
	}

	// encode body into json
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}

	// send post request
	resp, err := http.Post("http://localhost:80/api/posts", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	return nil
}
