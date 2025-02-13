package handlers

import (
	"os"
	"sync"
	"testing"
	"timelineDemo/internal/app/usecases"
	"timelineDemo/internal/domain/entities"
	"timelineDemo/internal/domain/repositories"
	"timelineDemo/internal/infrastructure/persistence"

	"github.com/google/uuid"
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
	userChannels                   map[uuid.UUID]chan entities.TimelineEvent
	mu                             sync.Mutex
}

func (s *HandlersTestSuite) SetupTest() {
	postMap := make(map[uuid.UUID][]*entities.Post)
	s.postsRepository = persistence.NewPostsRepository(&postMap)
	s.createPostUsecase = usecases.NewCreatePostsUsecase(s.postsRepository)
	s.getUserAndFolloweePostsUsecase = usecases.NewGetUserAndFolloweePostsUsecase(s.postsRepository)

	s.mu = sync.Mutex{}
	s.userChannels = make(map[uuid.UUID]chan entities.TimelineEvent)
}
