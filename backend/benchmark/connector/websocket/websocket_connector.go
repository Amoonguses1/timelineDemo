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
	fileio.WriteNewText("WSBenchLogs.txt", fmt.Sprintf("request send\n%s: %v\n", userID.String()[:7], time.Now()))
	url := fmt.Sprintf("ws://localhost:80/api/%s/ws", userID.String())
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	// notification for connection established
	connected <- struct{}{}

	for {
		_, message, err := conn.ReadMessage()
		fileio.WriteNewText("WSBenchLogs.txt", fmt.Sprintf("response received\n%s: %v\n", userID.String()[:7], time.Now()))

		if err != nil {
			log.Println("read:", err)
			break
		}

		if strings.Contains(string(message), "end") {
			log.Println("Close connection")
			CloseWebSocket(conn)
			break
		}
		log.Printf("message: %s", message)
	}
}

func CloseWebSocket(conn *websocket.Conn) {
	// send close frame
	message := "Client is closing the connection"
	err := conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, message))

	if err != nil {
		log.Println("Failed to send close frame:", err)
	} else {
		log.Println("Close frame sent")
	}
}
