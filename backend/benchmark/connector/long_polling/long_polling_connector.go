package longpolling

import (
	"benchmark/connector"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type LongPollingConnector struct {
}

func NewLongPollingConnector() connector.ConnectorInterface {
	return &LongPollingConnector{}
}

func (conn *LongPollingConnector) Connect(userID uuid.UUID, wg *sync.WaitGroup, connected chan<- struct{}) {
	defer wg.Done()

	// send request
	url := fmt.Sprintf("http://localhost:80/api/%s/polling?event_type=TimelineAccessed", userID.String())
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("User %s: error connecting: %v", userID, err)
		return
	}
	defer resp.Body.Close()

	// notification for connection established
	connected <- struct{}{}
	url = fmt.Sprintf("http://localhost:80/api/%s/polling?event_type=PollingRequest", userID.String())

	// read response
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			cancel()
			return
		}
		client := &http.Client{}
		resp, err = client.Do(req)
		cancel()
		if err != nil {
			fmt.Printf("Error connectiong to polling request for user %s: %v\n", userID.String(), err)
			return
		}

		byteArray, _ := io.ReadAll(resp.Body)
		text := string(byteArray)
		log.Println(text)
		if strings.Contains(text, "end") {
			return
		}
		resp.Body.Close()
	}
}
