package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	timelinegrpc "grpc-test/protogen/timeline"
	"grpc-test/timeline"

	"google.golang.org/grpc"
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
	// timeline.BasicTimelineResponse()

	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	timelinegrpc.RegisterTimelineServiceServer(s, timeline.NewGrpcServer())

	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}
