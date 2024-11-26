package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"timelineDemo/api/handlers"
	"timelineDemo/internal/app/usecases"
	"timelineDemo/internal/domain/entities"
	"timelineDemo/internal/infrastructure/persistence"
)

const (
	port = 80
)

func main() {
	postsMap := make(map[string][]*entities.Post)
	postsRepository := persistence.NewPostsRepository(&postsMap)
	createPostUsecase := usecases.NewCreatePostsUsecase(postsRepository)
	getUserAndFolloweePostsUsecase := usecases.NewGetUserAndFolloweePostsUsecase(postsRepository)
	var userChannels = make(map[string]chan entities.TimelineEvent)
	var mu sync.Mutex

	http.HandleFunc("POST /api/posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreatePost(w, r, &mu, &userChannels, createPostUsecase)
	})

	http.HandleFunc("GET /api/{id}/sse", func(w http.ResponseWriter, r *http.Request) {
		handlers.SseTimeline(w, r, getUserAndFolloweePostsUsecase, &mu, &userChannels)
	})

	log.Println("Starting server...")
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalln(err)
	}
}
