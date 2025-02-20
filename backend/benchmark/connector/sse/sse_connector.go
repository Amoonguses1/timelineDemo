package sse

import (
	"benchmark/connector"
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"sync"

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
