package repositories

import "timelineDemo/internal/domain/entities"

type PostsRepositoryInterface interface {
	GetUserAndFolloweePosts(userID string) ([]*entities.Post, error)
	CreatePost(userID string, text string) (entities.Post, error)
}
