package sse

import (
	"benchmark/connector"
	fileio "benchmark/fileIO"
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type SSEConnector struct {
}

func NewSSEConnector() connector.ConnectorInterface {
	return &SSEConnector{}
}

func (conn *SSEConnector) Connect(userID uuid.UUID, wg *sync.WaitGroup, connected chan<- struct{}) {
	defer wg.Done()

	// send request
	timestamp := time.Now().UTC().Format("15:04:05.000")
	fileio.WriteNewText("SSEBenchLogs.txt", fmt.Sprintf("send, %s, %s", userID.String()[:7], timestamp))
	url := fmt.Sprintf("http://localhost:80/api/%s/sse", userID.String())
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error connectiong to SSE for user %s: %v\n", userID.String(), err)
		return
	}

	// notification for connection established
	connected <- struct{}{}

	// read response
	scanner := bufio.NewScanner(resp.Body)
	defer resp.Body.Close()

	var responseSize int64
	var imageWg sync.WaitGroup

	for scanner.Scan() {
		// output response size
		responseSize += int64(len(scanner.Bytes()))
		text := scanner.Text()
		if len(text) > 0 && !strings.Contains(text, "TimelineAccessed") {
			timestamp := time.Now().UTC().Format("15:04:05.000")
			imageWg.Add(1)
			fileio.WriteNewText("SSEBenchLogs.txt", fmt.Sprintf("comes, %s, %s", userID.String()[:7], timestamp))
			go getImage(&imageWg, userID, "test.png")
		}
		if strings.Contains(text, "end") {
			imageWg.Wait()
			//fmt.Println(responseSize)
			return
		}
		// fmt.Printf("User: %s received: %s\n", userID.String(), text)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading for user %s: %v\n", userID.String(), err)
	}
}

func getImage(wg *sync.WaitGroup, userID uuid.UUID, imagePath string) {
	defer wg.Done()

	timestamp := time.Now().UTC().Format("15:04:05.000")
	fileio.WriteNewText("SSEBenchLogsImage.txt", fmt.Sprintf("send, %s, %s", userID.String()[:7], timestamp))
	url := fmt.Sprintf("http://localhost:80/api/%s/getimg?file=%s", userID.String(), imagePath)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error connectiong to SSE for user %s: %v\n", userID.String(), err)
		return
	}
	defer resp.Body.Close()

	timestamp = time.Now().UTC().Format("15:04:05.000")
	fileio.WriteNewText("SSEBenchLogsImage.txt", fmt.Sprintf("comes, %s, %s", userID.String()[:7], timestamp))
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return
	}
	log.Println(len(body))
}
