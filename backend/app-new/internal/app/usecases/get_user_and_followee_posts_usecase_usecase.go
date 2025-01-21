package usecases

import (
	"timelineDemo/internal/domain/entities"
	"timelineDemo/internal/domain/repositories"

	"github.com/google/uuid"
)

type GetUserAndFolloweePostsUsecaseInterface interface {
	GetUserAndFolloweePosts(userID uuid.UUID) ([]*entities.Post, error)
}

type getUserAndFolloweePostsUsecase struct {
	postsRepository repositories.PostsRepositoryInterface
}

func NewGetUserAndFolloweePostsUsecase(postsRepository repositories.PostsRepositoryInterface) GetUserAndFolloweePostsUsecaseInterface {
	return &getUserAndFolloweePostsUsecase{postsRepository: postsRepository}
}

func (p *getUserAndFolloweePostsUsecase) GetUserAndFolloweePosts(userID uuid.UUID) ([]*entities.Post, error) {
	posts, err := p.postsRepository.GetUserAndFolloweePosts(userID)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
