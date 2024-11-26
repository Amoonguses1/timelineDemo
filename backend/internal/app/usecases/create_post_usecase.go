package usecases

import (
	"timelineDemo/internal/domain/entities"
	"timelineDemo/internal/domain/repositories"
)

type CreatePostUsecaseInterface interface {
	CreatePost(userID string, text string) (entities.Post, error)
}

type createPostUsecase struct {
	postsRepository repositories.PostsRepositoryInterface
}

func NewCreatePostsUsecase(postsRepository repositories.PostsRepositoryInterface) CreatePostUsecaseInterface {
	return &createPostUsecase{postsRepository: postsRepository}
}

func (u *createPostUsecase) CreatePost(userID string, text string) (entities.Post, error) {
	post, err := u.postsRepository.CreatePost(userID, text)
	if err != nil {
		return entities.Post{}, err
	}

	return post, nil
}
