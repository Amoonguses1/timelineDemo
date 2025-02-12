package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
	server "timelineDemo/api"
	"timelineDemo/api/handlers"
	"timelineDemo/api/middlewares"
	timelinegrpc "timelineDemo/grpc/protogen/timeline"
	"timelineDemo/internal/app/usecases"
	"timelineDemo/internal/domain/entities"
	"timelineDemo/internal/infrastructure/persistence"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

const (
	port     = 80
	grpcPort = 8081
)

func main() {
	postsMap := make(map[uuid.UUID][]*entities.Post)
	postsRepository := persistence.NewPostsRepository(&postsMap)
	createPostUsecase := usecases.NewCreatePostsUsecase(postsRepository)
	getUserAndFolloweePostsUsecase := usecases.NewGetUserAndFolloweePostsUsecase(postsRepository)
	var userChannels = make(map[uuid.UUID]chan entities.TimelineEvent)
	var mu sync.Mutex

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreatePost(w, r, &mu, &userChannels, createPostUsecase)
	})

	mux.HandleFunc("GET /api/{id}/sse", func(w http.ResponseWriter, r *http.Request) {
		handlers.SseTimeline(w, r, getUserAndFolloweePostsUsecase, &mu, &userChannels)
	})

	mux.HandleFunc("GET /api/{id}/polling", func(w http.ResponseWriter, r *http.Request) {
		handlers.LongPollingTimeline(w, r, getUserAndFolloweePostsUsecase, &mu, &userChannels)
	})

	// Register a WebSocket endpoint to provide real-time timeline updates.
	// This handler establishes a WebSocket connection for a user identified by {id}
	// and streams timeline updates including new posts or deletions from followed users.
	mux.HandleFunc("/api/{id}/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.WebSocketTimeline(w, r, getUserAndFolloweePostsUsecase, &mu, &userChannels)
	})

	handlersWithCORS := middlewares.CORS(mux)

	// setup grpc server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	timelinegrpc.RegisterTimelineServiceServer(s, server.NewGrpcServer(getUserAndFolloweePostsUsecase, &mu, &userChannels))

	// activate grpc server
	go func() {
		log.Printf("Starting gRPC server on port: %d", grpcPort)
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// activate http server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handlersWithCORS,
	}
	go func() {
		log.Printf("Starting HTTP server on port: %d", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve HTTP: %v", err)
		}
	}()

	// signal handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Stopping servers...")

	// stop grpc server
	s.GracefulStop()

	// stop http server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server Shutdown: %v", err)
	}

	log.Println("Servers stopped successfully.")
}
