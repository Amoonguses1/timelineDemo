package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

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
