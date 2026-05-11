package memory

import (
	"errors"
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
	"sync"
)

type postRepositoryImpl struct {
	posts map[string]entities.Post
	mu    sync.RWMutex
}

func NewPostRepository() repositories.PostRepository {
	return &postRepositoryImpl{
		posts: make(map[string]entities.Post),
	}
}

func (r *postRepositoryImpl) Create(post entities.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.posts[post.ID] = post
	return nil
}

func (r *postRepositoryImpl) FindAll() ([]entities.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	posts := make([]entities.Post, 0, len(r.posts))
	for _, post := range r.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *postRepositoryImpl) FindByID(id string) (entities.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, exists := r.posts[id]
	if !exists {
		return entities.Post{}, errors.New("post not found")
	}
	return post, nil
}

func (r *postRepositoryImpl) FindByUserID(userID string) ([]entities.Post, error) {
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

func (r *postRepositoryImpl) Update(post entities.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.posts[post.ID]; !exists {
		return errors.New("post not found")
	}
	r.posts[post.ID] = post
	return nil
}

func (r *postRepositoryImpl) Exists(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.posts[id]
	return exists
}
