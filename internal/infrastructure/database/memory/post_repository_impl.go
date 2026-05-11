package memory

import (
	"errors"
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
	"sync"
)

type PostRepositoryImpl struct {
	posts map[string]entities.Post
	mu    sync.RWMutex
}

func NewPostRepository() repositories.PostRepository {
	return &PostRepositoryImpl{
		posts: make(map[string]entities.Post),
	}
}

func (r *PostRepositoryImpl) Create(post entities.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.posts[post.ID] = post
	return nil
}

func (r *PostRepositoryImpl) FindAll() ([]entities.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	posts := make([]entities.Post, 0, len(r.posts))
	for _, post := range r.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepositoryImpl) FindByID(id string) (entities.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, exists := r.posts[id]
	if !exists {
		return entities.Post{}, errors.New("post not found")
	}
	return post, nil
}

func (r *PostRepositoryImpl) FindByUserID(userID string) ([]entities.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	posts := make([]entities.Post, 0)
	for _, post := range r.posts {
		if post.UserID == userID {
			posts = append(posts, post)
		}
	}
	return posts, nil
}

func (r *PostRepositoryImpl) Update(post entities.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.posts[post.ID]; !exists {
		return errors.New("post not found")
	}
	r.posts[post.ID] = post
	return nil
}

func (r *PostRepositoryImpl) Exists(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.posts[id]
	return exists
}
