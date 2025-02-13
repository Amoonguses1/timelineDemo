package repositories

import (
	"timelineDemo/internal/domain/entities"

	"github.com/google/uuid"
)

type PostsRepositoryInterface interface {
	GetUserAndFolloweePosts(userID uuid.UUID) ([]*entities.Post, error)
	CreatePost(userID uuid.UUID, text string) (*entities.Post, error)
}
