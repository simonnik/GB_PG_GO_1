package service

import "github.com/simonnik/GB_PG_GO_1/hw5/pkg/post/storage"

type PostService struct {
	store storage.DB
}

func NewPostService(db storage.DB) *PostService {
	return &PostService{db}
}

// GetPosts returns list of posts
func (p *PostService) GetPosts() ([]*storage.Post, error) {
	posts, err := p.store.GetPosts()
	return posts, err
}
