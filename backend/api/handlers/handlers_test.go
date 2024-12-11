package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
	"timelineDemo/internal/domain/entities"

	"github.com/stretchr/testify/suite"
)

func (s *HandlersTestSuite) TestCreatePost() {
	tests := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{
			name:         "create post",
			body:         `{ "user_id": "user1", "text": "test2" }`,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid JSON body",
			body:         `{ "text": "test2" `,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("POST", "/api/posts", strings.NewReader(test.body))
		rr := httptest.NewRecorder()

		CreatePost(rr, req, &s.mu, &s.userChannels, s.createPostUsecase)

		if rr.Code != test.expectedCode {
			s.T().Errorf("%s: wrong code returned; expected %d, but got %d", test.name, test.expectedCode, rr.Code)
		}
	}
}

func (s *HandlersTestSuite) TestSseTimeline() {
	user1 := "user1"
	s.newTestPost(fmt.Sprintf(`{ "user_id": "%s", "text": "test5" }`, user1))

	tests := []struct {
		name          string
		userID        string
		expectedCount int
	}{
		{
			name:          "get posts",
			userID:        user1,
			expectedCount: 1,
		},
		{
			name:          "get no posts",
			userID:        "",
			expectedCount: 0,
		},
		{
			name:          "get posts already posted and posts posted during timeline access",
			userID:        user1,
			expectedCount: 2,
		},
	}

	for _, test := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(
			"GET",
			"/api/{id}/sse",
			strings.NewReader(""),
		).WithContext(ctx)
		req.SetPathValue("id", test.userID)

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			SseTimeline(rr, req, s.getUserAndFolloweePostsUsecase, &s.mu, &s.userChannels)
		}()

		if test.name == "get posts already posted and posts posted during timeline access" {
			time.Sleep(100 * time.Millisecond)
			s.newTestPost(fmt.Sprintf(`{ "user_id": "%s", "text": "test5" }`, test.userID))
		}

		wg.Wait()
		scanner := bufio.NewScanner(rr.Body)
		var posts []entities.Post

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data:") {
				jsonData := strings.TrimPrefix(line, "data: ")
				var timelineEvent entities.TimelineEvent

				err := json.Unmarshal([]byte(jsonData), &timelineEvent)
				if err != nil {
					s.T().Errorf("Failed to decode JSON: %v", err)
				}
				for _, post := range timelineEvent.Posts {
					posts = append(posts, *post)
				}
			}
		}

		if len(posts) != test.expectedCount {
			s.T().Errorf("%s: wrong number of posts returned; expected %d, but got %d", test.name, test.expectedCount, len(posts))
		}
	}
}

func (s *HandlersTestSuite) TestLongPollingTimeline() {
	user1 := "user1"
	s.newTestPost(fmt.Sprintf(`{ "user_id": "%s", "text": "test5" }`, user1))

	tests := []struct {
		name         string
		userID       string
		body         string
		expectedCode int
	}{
		{
			name:         "get posts",
			userID:       user1,
			body:         fmt.Sprintf(`{"polling_event_type": "%s"}`, TimelineAccessed),
			expectedCode: http.StatusOK,
		},
		{
			name:         "get posts already posted and posts posted during timeline access",
			userID:       user1,
			body:         fmt.Sprintf(`{"polling_event_type": "%s"}`, PollingRequest),
			expectedCode: http.StatusOK,
		},
		{
			name:         "no new posts",
			userID:       user1,
			body:         fmt.Sprintf(`{"polling_event_type": "%s"}`, PollingRequest),
			expectedCode: http.StatusNoContent,
		},
	}

	for _, test := range tests {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(
			"GET",
			"/api/{id}/polling",
			strings.NewReader(test.body),
		)
		req.SetPathValue("id", test.userID)

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			LongPollingTimeline(rr, req, s.getUserAndFolloweePostsUsecase, &s.mu, &s.userChannels)
		}()

		if test.name == "get posts already posted and posts posted during timeline access" {
			time.Sleep(100 * time.Millisecond)
			s.newTestPost(fmt.Sprintf(`{ "user_id": "%s", "text": "test5" }`, test.userID))
		}

		wg.Wait()

		if rr.Code != test.expectedCode {
			s.T().Errorf("%s: wrong code returned; expected %d, but got %d", test.name, test.expectedCode, rr.Code)
		}
	}
}

func (s *HandlersTestSuite) newTestPost(body string) {
	req := httptest.NewRequest("POST", "/api/posts", strings.NewReader(body))
	rr := httptest.NewRecorder()

	CreatePost(rr, req, &s.mu, &s.userChannels, s.createPostUsecase)
}

// TestHandlersTestSuite runs all of the tests attached to HandlersTestSuite.
func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}
