package usecases

import (
	"timelineDemo/internal/domain/entities"
	"timelineDemo/internal/domain/repositories"
)

type GetUserAndFolloweePostsUsecaseInterface interface {
	GetUserAndFolloweePosts(userID string) ([]*entities.Post, error)
}

type getUserAndFolloweePostsUsecase struct {
	postsRepository repositories.PostsRepositoryInterface
}

func NewGetUserAndFolloweePostsUsecase(postsRepository repositories.PostsRepositoryInterface) GetUserAndFolloweePostsUsecaseInterface {
	return &getUserAndFolloweePostsUsecase{postsRepository: postsRepository}
}

func (p *getUserAndFolloweePostsUsecase) GetUserAndFolloweePosts(userID string) ([]*entities.Post, error) {
	posts, err := p.postsRepository.GetUserAndFolloweePosts(userID)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
