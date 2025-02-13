package usecases

import (
	"timelineDemo/internal/domain/entities"
	"timelineDemo/internal/domain/repositories"

	"github.com/google/uuid"
)

type CreatePostUsecaseInterface interface {
	CreatePost(userID uuid.UUID, text string) (*entities.Post, error)
}

type createPostUsecase struct {
	postsRepository repositories.PostsRepositoryInterface
}

func NewCreatePostsUsecase(postsRepository repositories.PostsRepositoryInterface) CreatePostUsecaseInterface {
	return &createPostUsecase{postsRepository: postsRepository}
}

func (u *createPostUsecase) CreatePost(userID uuid.UUID, text string) (*entities.Post, error) {
	post, err := u.postsRepository.CreatePost(userID, text)
	if err != nil {
		return nil, err
	}

	return post, nil
}
