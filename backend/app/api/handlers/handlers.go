package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	usecases "timelineDemo/internal/app/usecases"
	"timelineDemo/internal/domain/entities"
	fileio "timelineDemo/internal/infrastructure/fileIO"

	"github.com/google/uuid"
)

func CreatePost(w http.ResponseWriter, r *http.Request, mu *sync.Mutex, usersChan *map[uuid.UUID]chan entities.TimelineEvent, createPostUsecase usecases.CreatePostUsecaseInterface) {
	var body createPostRequestBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintln("Request body was invalid."), http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(body.UserID)
	if err != nil {
		http.Error(w, fmt.Sprintln("invalid user id"), http.StatusBadRequest)
		return
	}
	log.Println("createPost request arrived:", userID.String()[:7])

	post, err := createPostUsecase.CreatePost(userID, body.Text)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not create a post."), http.StatusInternalServerError)
		return
	}

	go func(userChan *map[uuid.UUID]chan entities.TimelineEvent) {
		var posts []*entities.Post
		posts = append(posts, post)

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

func SseTimeline(w http.ResponseWriter, r *http.Request, u usecases.GetUserAndFolloweePostsUsecaseInterface, mu *sync.Mutex, usersChan *map[uuid.UUID]chan entities.TimelineEvent, isBench bool) {
	log.Println("app-new sse connection called")
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, fmt.Sprintln("invalid user id"), http.StatusBadRequest)
		return
	}
	if isBench {
		timestamp := time.Now().Format("15:04:05.000")
		fileio.WriteNewText("SSEBenchLogs.txt", fmt.Sprintf("comes, %s, %s", userID.String()[:7], timestamp))
	}

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

	flusher, _ := w.(http.Flusher)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	log.Println("header set")

	// send initial response
	jsonData, err := json.Marshal(entities.TimelineEvent{EventType: entities.TimelineAccessed, Posts: posts})
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	flusher.Flush()

	for {
		select {
		case event := <-userChan:
			log.Println("event comes in app-new")
			jsonData, err := json.Marshal(event)
			if err != nil {
				log.Println(err)
				return
			}

			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			if isBench {
				timestamp := time.Now().Format("15:04:05.000")
				fileio.WriteNewText("SSEBenchLogs.txt", fmt.Sprintf("send, %s, %s", userID.String()[:7], timestamp))
			}
			flusher.Flush()
		case <-r.Context().Done():
			log.Println("Connection closed")
			mu.Lock()
			delete(*usersChan, userID)
			mu.Unlock()
			return
		}
	}
}
