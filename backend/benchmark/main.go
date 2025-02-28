package main

import (
	websocketconnector "benchmark/connector/websocket"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

func main() {
	// Determine connection type based on the command-line argument
	// conn := sse.NewSSEConnector()
	// conn := longpolling.NewLongPollingConnector()
	// conn := grpcconnector.NewGRPCConnector()
	conn := websocketconnector.NewWebSocketConnector()

	// Set up cases for benchmarking.
	benchCases := []struct {
		name           string
		numPosters     int
		PostsPerUser   int
		NumConnections int
	}{
		{
			name:           "influencer case",
			numPosters:     1,
			PostsPerUser:   1,
			NumConnections: 25,
		},
		{
			name:           "many posters case",
			numPosters:     5,
			PostsPerUser:   5,
			NumConnections: 1,
		},
		{
			name:           "multiple posters, multiple connected clients case",
			numPosters:     5,
			PostsPerUser:   5,
			NumConnections: 25,
		},
	}
	benchCase := benchCases[2]

	log.Println(benchCase.name)

	// connections
	var wg sync.WaitGroup
	connected := make(chan struct{}, benchCase.NumConnections)

	for i := 0; i < benchCase.NumConnections; i++ {
		wg.Add(1)
		go conn.Connect(uuid.New(), &wg, connected)
		time.Sleep(100 * time.Millisecond)
	}

	for i := 0; i < benchCase.NumConnections; i++ {
		<-connected
	}

	log.Println("Send requests")
	BenchRequests(benchCase.numPosters, benchCase.PostsPerUser)
	EndRequest()
	wg.Wait()
}
