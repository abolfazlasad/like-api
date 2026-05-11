package memory

import (
	"errors"
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
	"sync"
)

type LikeRepositoryImpl struct {
	likes map[string]entities.Like
	mu    sync.RWMutex
}

func NewLikeRepository() repositories.LikeRepository {
	return &LikeRepositoryImpl{
		likes: make(map[string]entities.Like),
	}
}

func (r *LikeRepositoryImpl) Create(like entities.Like) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.likes[like.ID] = like
	return nil
}

func (r *LikeRepositoryImpl) Delete(userID, postID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, like := range r.likes {
		if like.UserID == userID && like.PostID == postID {
			delete(r.likes, id)
			return nil
		}
	}
	return errors.New("like not found")
}

func (r *LikeRepositoryImpl) FindByPostID(postID string) ([]entities.Like, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	likes := make([]entities.Like, 0)
	for _, like := range r.likes {
		if like.PostID == postID {
			likes = append(likes, like)
		}
	}
	return likes, nil
}

func (r *LikeRepositoryImpl) FindByUserID(userID string) ([]entities.Like, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	likes := make([]entities.Like, 0)
	for _, like := range r.likes {
		if like.UserID == userID {
			likes = append(likes, like)
		}
	}
	return likes, nil
}

func (r *LikeRepositoryImpl) Exists(userID, postID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, like := range r.likes {
		if like.UserID == userID && like.PostID == postID {
			return true
		}
	}
	return false
}
