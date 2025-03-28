package websocketconnector

import (
	"benchmark/connector"
	fileio "benchmark/fileIO"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketConnector struct {
}

func NewWebSocketConnector() connector.ConnectorInterface {
	return &WebSocketConnector{}
}

func (c *WebSocketConnector) Connect(userID uuid.UUID, wg *sync.WaitGroup, connected chan<- struct{}) {
	defer wg.Done()

	// connect to the WebSocket endpoint
	timestamp := time.Now().UTC().Format("15:04:05.000")
	fileio.WriteNewText("WSBenchLogs.txt", fmt.Sprintf("send, %s, %s", userID.String()[:7], timestamp))
	url := fmt.Sprintf("ws://localhost:80/api/%s/ws", userID.String())
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	// notification for connection established
	connected <- struct{}{}
	var receivedBytes int64
	var imageWg sync.WaitGroup

	for {
		_, message, err := conn.ReadMessage()
		receivedBytes += int64(len(message))
		if !strings.Contains(string(message), "TimelineAccessed") {
			timestamp := time.Now().UTC().Format("15:04:05.000")
			fileio.WriteNewText("WSBenchLogs.txt", fmt.Sprintf("comes, %s, %s", userID.String()[:7], timestamp))

			imageWg.Add(1)
			go getImage(&imageWg, userID)
		}

		if err != nil {
			log.Println("read:", err)
			break
		}

		if strings.Contains(string(message), "end") {
			imageWg.Wait()
			// log.Println("Close connection")
			CloseWebSocket(conn)
			break
		}
		// log.Printf("message: %s", message)
	}

	// fmt.Println(receivedBytes)
}

func CloseWebSocket(conn *websocket.Conn) {
	// send close frame
	message := "Client is closing the connection"
	err := conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, message))

	if err != nil {
		log.Println("Failed to send close frame:", err)
	}
}

func getImage(wg *sync.WaitGroup, userID uuid.UUID) {
	defer wg.Done()

	// connect to the WebSocket endpoint
	timestamp := time.Now().UTC().Format("15:04:05.000")
	fileio.WriteNewText("WSBenchLogsImage.txt", fmt.Sprintf("send, %s, %s", userID.String()[:7], timestamp))
	url := fmt.Sprintf("ws://localhost:80/api/%s/getimg/ws", userID.String())
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte("test.png"))
	if err != nil {
		log.Fatal("write:", err)
	}

	var receivedBytes int64

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		receivedBytes += int64(len(message))
	}

	fileio.WriteNewText("WSBenchLogsImage.txt", fmt.Sprintf("comes, %s, %s", userID.String()[:7], timestamp))
	log.Printf("Image resp: %d\n", receivedBytes)
}
