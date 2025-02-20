package persistence

import (
	"sync"
	"time"
	"timelineDemo/internal/domain/entities"
	"timelineDemo/internal/domain/repositories"

	"github.com/google/uuid"
)

type PostsRepositoryInMemory struct {
	postsMap *map[uuid.UUID][]*entities.Post
	lock     sync.RWMutex
}

func NewPostsRepository(postsMap *map[uuid.UUID][]*entities.Post) repositories.PostsRepositoryInterface {
	return &PostsRepositoryInMemory{postsMap: postsMap}
}

func (p *PostsRepositoryInMemory) GetUserAndFolloweePosts(userID uuid.UUID) ([]*entities.Post, error) {
	// p.lock.RLock()
	posts := (*p.postsMap)[userID]
	// p.lock.RUnlock()
	return posts, nil
}

func (p *PostsRepositoryInMemory) CreatePost(userID uuid.UUID, text string) (*entities.Post, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	post := entities.Post{
		ID:        uuid,
		UserID:    userID,
		Text:      text,
		CreatedAt: time.Now(),
	}

	p.lock.Lock()
	posts := (*p.postsMap)[userID]
	posts = append(posts, &post)
	(*p.postsMap)[userID] = posts

	p.lock.Unlock()
	return &post, nil
}
