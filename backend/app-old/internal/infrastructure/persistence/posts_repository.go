package persistence

import (
	"timelineDemo-old/internal/domain/entities"
	"timelineDemo-old/internal/domain/repositories"
)

type PostsRepositoryInMemory struct {
	postsMap *map[string][]*entities.Post
}

func NewPostsRepository(postsMap *map[string][]*entities.Post) repositories.PostsRepositoryInterface {
	return &PostsRepositoryInMemory{postsMap: postsMap}
}

func (p *PostsRepositoryInMemory) GetUserAndFolloweePosts(userID string) ([]*entities.Post, error) {
	posts := (*p.postsMap)[userID]
	return posts, nil
}

func (p *PostsRepositoryInMemory) CreatePost(userID string, text string) (entities.Post, error) {
	post := entities.Post{
		UserID: userID,
		Text:   text,
	}
	posts := (*p.postsMap)[userID]
	posts = append(posts, &post)
	(*p.postsMap)[userID] = posts

	return post, nil
}
