package main

import (
	"fmt"
	"grpc-test/timeline"
	"log"
	"time"
)

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().Format("15:04:05" + " " + string(bytes)))
}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	// timeline.BasicPosts()
	timeline.BasicTimelineResponse()
}
