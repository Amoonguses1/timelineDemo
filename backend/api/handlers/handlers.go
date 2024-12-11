package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	usecases "timelineDemo/internal/app/usecases"
	"timelineDemo/internal/domain/entities"
)

const (
	expiredTime = 5 * time.Second
)

func CreatePost(w http.ResponseWriter, r *http.Request, mu *sync.Mutex, usersChan *map[string]chan entities.TimelineEvent, createPostUsecase usecases.CreatePostUsecaseInterface) {
	var body createPostRequestBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintln("Request body was invalid."), http.StatusBadRequest)
		return
	}

	post, err := createPostUsecase.CreatePost(body.UserID, body.Text)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not create a post."), http.StatusInternalServerError)
		return
	}

	go func(userChan *map[string]chan entities.TimelineEvent) {
		var posts []*entities.Post
		posts = append(posts, &post)

		mu.Lock()
		for _, userChan := range *usersChan {
			userChan <- entities.TimelineEvent{EventType: entities.PostCreated, Posts: posts}
		}
		mu.Unlock()
	}(usersChan)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(&post)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not encode response."), http.StatusInternalServerError)
		return
	}
}

func SseTimeline(w http.ResponseWriter, r *http.Request, u usecases.GetUserAndFolloweePostsUsecaseInterface, mu *sync.Mutex, usersChan *map[string]chan entities.TimelineEvent) {
	userID := r.PathValue("id")
	posts, err := u.GetUserAndFolloweePosts(userID)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not get posts"), http.StatusInternalServerError)
		return
	}

	mu.Lock()
	if _, exists := (*usersChan)[userID]; !exists {
		(*usersChan)[userID] = make(chan entities.TimelineEvent, 1)
	}
	userChan := (*usersChan)[userID]
	mu.Unlock()

	userChan <- entities.TimelineEvent{EventType: entities.TimelineAccessed, Posts: posts}

	flusher, _ := w.(http.Flusher)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
		case event := <-userChan:
			jsonData, err := json.Marshal(event)
			if err != nil {
				log.Println(err)
				return
			}

			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			flusher.Flush()
		case <-r.Context().Done():
			mu.Lock()
			delete(*usersChan, userID)
			mu.Unlock()
			return
		}
	}
}

func LongPollingTimeline(w http.ResponseWriter, r *http.Request, u usecases.GetUserAndFolloweePostsUsecaseInterface, mu *sync.Mutex, usersChan *map[string]chan entities.TimelineEvent) {
	userID := r.PathValue("id")
	var body longPollingTimelineRequestBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
        log.Println("decode error")
		http.Error(w, fmt.Sprintln("Request body was invalid."), http.StatusBadRequest)
		return
	}

	switch body.PollingEventType {
	case TimelineAccessed:
		handleTimelineAccess(w, userID, u)
	case PollingRequest:
		handlePollingRequest(w, r, userID, usersChan, mu)
	default:
		http.Error(w, fmt.Sprintln("Unknown event type"), http.StatusBadRequest)
	}
}

func handleTimelineAccess(w http.ResponseWriter, userID string, u usecases.GetUserAndFolloweePostsUsecaseInterface) {
	posts, err := u.GetUserAndFolloweePosts(userID)
	if err != nil {
		http.Error(w, "Could not get posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&posts); err != nil {
		http.Error(w, "Could not encode response.", http.StatusInternalServerError)
	}
}

func handlePollingRequest(w http.ResponseWriter, r *http.Request, userID string, usersChan *map[string]chan entities.TimelineEvent, mu *sync.Mutex) {
	ctx, cancel := context.WithTimeout(r.Context(), expiredTime)
	defer cancel()

	mu.Lock()
	if _, exists := (*usersChan)[userID]; !exists {
		(*usersChan)[userID] = make(chan entities.TimelineEvent, 1)
	}
	userChan := (*usersChan)[userID]
	mu.Unlock()

	select {
	case event := <-userChan:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(event.Posts); err != nil {
			http.Error(w, "Could not encode response.", http.StatusInternalServerError)
		}
	case <-ctx.Done():
		mu.Lock()
		delete(*usersChan, userID)
		mu.Unlock()
		w.WriteHeader(http.StatusNoContent)
	}
}
