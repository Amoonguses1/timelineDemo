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

var IsBench bool

func main() {
	// Determine if the server should run in benchmark mode based on the command-line argument.
	if len(os.Args) > 2 && os.Args[2] == "bench" {
		IsBench = true
		log.Println("Benchmark mode")
	}

	// Initialize usecases and repositories for managing posts.
	postsMap := make(map[uuid.UUID][]*entities.Post)
	postsRepository := persistence.NewPostsRepository(&postsMap)
	createPostUsecase := usecases.NewCreatePostsUsecase(postsRepository)
	getUserAndFolloweePostsUsecase := usecases.NewGetUserAndFolloweePostsUsecase(postsRepository)

	// Initialize channels and a mutex for managing real-time notifications.
	var userChannels = make(map[uuid.UUID]chan entities.TimelineEvent)
	var mu sync.Mutex

	// Create an HTTP request multiplexer.
	mux := http.NewServeMux()

	// Register an endpoint for creating posts.
	// This handler creates a new post and sends real-time notifications to users following the poster.
	mux.HandleFunc("POST /api/posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreatePost(w, r, &mu, &userChannels, createPostUsecase)
	})

	// Register an endpoint to provide real-time timeline updates using Server-Sent Events (SSE).
	// This handler establishes a text/event-stream connection for a user identified by {id}
	// and continuously streams timeline updates, including new posts or deletions from followed users.
	mux.HandleFunc("GET /api/{id}/sse", func(w http.ResponseWriter, r *http.Request) {
		handlers.SseTimeline(w, r, getUserAndFolloweePostsUsecase, &mu, &userChannels, IsBench)
	})

	// Register an endpoint for polling to fetch timeline updates.
	// This handler supports two types of requests:
	// 1. Initial access request to fetch existing timeline data.
	// 2. Long polling request that waits for new updates before responding.
	mux.HandleFunc("GET /api/{id}/polling", func(w http.ResponseWriter, r *http.Request) {
		handlers.LongPollingTimeline(w, r, getUserAndFolloweePostsUsecase, &mu, &userChannels)
	})

	// Register a WebSocket endpoint to provide real-time timeline updates.
	// This handler establishes a WebSocket connection for a user identified by {id}
	// and streams timeline updates including new posts or deletions from followed users.
	mux.HandleFunc("/api/{id}/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.WebSocketTimeline(w, r, getUserAndFolloweePostsUsecase, &mu, &userChannels, IsBench)
	})

	handlersWithCORS := middlewares.CORS(mux)

	// Set up a gRPC server.
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	timelinegrpc.RegisterTimelineServiceServer(s, server.NewGrpcServer(getUserAndFolloweePostsUsecase, &mu, &userChannels, IsBench))

	// Start the gRPC server in a separate goroutine.
	go func() {
		log.Printf("Starting gRPC server on port: %d", grpcPort)
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// activate an HTTP server.
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handlersWithCORS,
	}

	// Start the HTTP server in a separate goroutine.
	go func() {
		log.Printf("Starting HTTP server on port: %d", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve HTTP: %v", err)
		}
	}()

	// Handle shutdown signals to gracefully stop servers.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Stopping servers...")

	// Gracefully stop the gRPC server.
	s.GracefulStop()

	// Gracefully shut down the HTTP server with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server Shutdown: %v", err)
	}

	log.Println("Servers stopped successfully.")
}
