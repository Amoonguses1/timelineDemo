package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	usecases "timelineDemo/internal/app/usecases"
	"timelineDemo/internal/domain/entities"
	fileio "timelineDemo/internal/infrastructure/fileIO"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WebSocketTimeline is a websocket handler for timeline.
// Keep the connection connected when a request is received and return posts.
// While a connection is alive, return a response every time an event comes in.
func WebSocketTimeline(w http.ResponseWriter, r *http.Request, u usecases.GetUserAndFolloweePostsUsecaseInterface, mu *sync.Mutex, usersChan *map[uuid.UUID]chan entities.TimelineEvent, isBench bool) {
	// Extract and parse the user ID from the request path.
	// If the ID is not a valid UUID, log the error and terminate the connection.
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		log.Println("Falied to parse userID:", err)
		return
	}

	// checkOrigin sets the CORS configucation.
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	// Upgrade the HTTP connection to a WebSocket connection.
	// If the upgrade fails, log the error and terminate connection.
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
	}
	defer ws.Close()
	log.Println("Client Connected")
	if isBench {
		fileio.WriteNewText("WSBenchLogs.txt", fmt.Sprintf("request comes\n%s: %v\n", userID.String()[:7], time.Now()))
	}

	// Set up WebSocket Close flow.
	ws.SetCloseHandler(func(code int, text string) error {
		log.Printf("WebSocket closed: User=%s, Code=%d, Reason=%s", userID, code, text)

		mu.Lock()
		delete(*usersChan, userID)
		mu.Unlock()

		return nil
	})

	// Fetch the user's posts and posts from followed users using the use case.
	// If retrieval fails, log the error and terminate connection.
	posts, err := u.GetUserAndFolloweePosts(userID)
	if err != nil {
		log.Println("Failed to retrieve posts:", err)
		return
	}

	// Create channels to receive the event notifications.
	mu.Lock()
	if _, exists := (*usersChan)[userID]; !exists {
		(*usersChan)[userID] = make(chan entities.TimelineEvent, 1)
	}
	userChan := (*usersChan)[userID]
	mu.Unlock()

	// Send the initial access response.
	if isBench {
		fileio.WriteNewText("WSBenchLogs.txt", fmt.Sprintf("response send\n%s: %v\n", userID.String()[:7], time.Now()))
	}
	err = ws.WriteJSON(entities.TimelineEvent{EventType: entities.TimelineAccessed, Posts: posts})
	if err != nil {
		log.Println("Failed to send the initial access response:", err)
		return
	}

	closeChan := make(chan struct{})
	go func() {
		defer close(closeChan)

		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					return
				}
				log.Println("WebSocket read error:", err)
				return
			}
		}
	}()

	// Continuously listen for notifications while the connection is active.
	for {
		select {
		// Handle incoming timeline events and send them to the client.
		case event := <-userChan:
			err = ws.WriteJSON(event)
			if err != nil {
				log.Println("Failed to send the event notification:", err)
				return
			}

		// Handle WebSocket disconnection and clean up the user channel.
		case <-r.Context().Done():
			// Remove the user from the active WebSocket channel map.
			mu.Lock()
			delete(*usersChan, userID)
			mu.Unlock()
			return

		// Handle WebSocket disconnection message and clean up the user channel.
		case <-closeChan:
			// Remove the user from the active WebSocket channel map.
			mu.Lock()
			delete(*usersChan, userID)
			mu.Unlock()
			return
		}
	}
}
