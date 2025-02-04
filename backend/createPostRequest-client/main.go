package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func main() {
	// setup default values
	intervalMs := 100
	totalRequests := 10

	// Check command line args
	if len(os.Args) > 2 {
		if v, err := strconv.Atoi(os.Args[1]); err == nil {
			intervalMs = v
		} else {
			fmt.Println("Invalid intervalMs, using default:", intervalMs)
		}

		if v, err := strconv.Atoi(os.Args[2]); err == nil {
			totalRequests = v
		} else {
			fmt.Println("Invalid totalRequests, using default:", totalRequests)
		}
	} else {
		fmt.Printf("Usage: go run main.go <intervalMs> <totalRequests>\n"+
			"Using default params: intervalMs=%d, totalRequests=%d\n", intervalMs, totalRequests)

	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Send the specified number of requests at the specified intervals
	SendRequests(ctx, intervalMs, totalRequests)

}

func SendRequests(ctx context.Context, interval, totalRequests int) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(interval))
	defer ticker.Stop()

	for i := 0; i < totalRequests; i++ {
		select {
		case <-ticker.C:
			userID, _ := uuid.NewRandom()
			err := sendRequest(userID, fmt.Sprintf("text no. %d", i))
			if err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

type CreatePostRequestBody struct {
	UserID string `json:"user_id,omitempty"`
	Text   string `json:"text"`
}

func sendRequest(userID uuid.UUID, text string) error {
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
