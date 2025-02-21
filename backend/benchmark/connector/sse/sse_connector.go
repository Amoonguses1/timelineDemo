package sse

import (
	"benchmark/connector"
	fileio "benchmark/fileIO"
	"bufio"
	"fmt"
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
	fileio.WriteNewText("SSEBenchLogs.txt", fmt.Sprintf("request send\n%s: %v\n", userID.String()[:7], time.Now()))
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
	for scanner.Scan() {
		fileio.WriteNewText("SSEBenchLogs.txt", fmt.Sprintf("response received\n%s: %v\n", userID.String()[:7], time.Now()))
		text := scanner.Text()
		if strings.Contains(text, "end") {
			return
		}
		fmt.Printf("User: %s received: %s\n", userID.String(), text)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading for user %s: %v\n", userID.String(), err)
	}

	resp.Body.Close()
}
