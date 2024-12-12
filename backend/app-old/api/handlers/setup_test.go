package handlers

import (
	"os"
	"sync"
	"testing"
	"timelineDemo-old/internal/app/usecases"
	"timelineDemo-old/internal/domain/entities"
	"timelineDemo-old/internal/domain/repositories"
	"timelineDemo-old/internal/infrastructure/persistence"

	"github.com/stretchr/testify/suite"
)

func TestMain(m *testing.M) {
	code := m.Run()

	os.Exit(code)
}

type HandlersTestSuite struct {
	suite.Suite
	createPostUsecase              usecases.CreatePostUsecaseInterface
	getUserAndFolloweePostsUsecase usecases.GetUserAndFolloweePostsUsecaseInterface
	postsRepository                repositories.PostsRepositoryInterface
	userChannels                   map[string]chan entities.TimelineEvent
	mu                             sync.Mutex
}

func (s *HandlersTestSuite) SetupTest() {
	postMap := make(map[string][]*entities.Post)
	s.postsRepository = persistence.NewPostsRepository(&postMap)
	s.createPostUsecase = usecases.NewCreatePostsUsecase(s.postsRepository)
	s.getUserAndFolloweePostsUsecase = usecases.NewGetUserAndFolloweePostsUsecase(s.postsRepository)

	s.mu = sync.Mutex{}
	s.userChannels = make(map[string]chan entities.TimelineEvent)
}
