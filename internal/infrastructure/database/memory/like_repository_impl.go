package memory

import (
	"errors"
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
	"sync"
)

type likeRepositoryImpl struct {
	likes map[string]entities.Like
	mu    sync.RWMutex
}

func newLikeRepository() repositories.LikeRepository {
	return &likeRepositoryImpl{
		likes: make(map[string]entities.Like),
	}
}

func (r *likeRepositoryImpl) Create(like entities.Like) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.likes[like.ID] = like
	return nil
}

func (r *likeRepositoryImpl) Delete(userID, postID string) error {
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

func (r *likeRepositoryImpl) FindByPostID(postID string) ([]entities.Like, error) {
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

func (r *likeRepositoryImpl) FindByUserID(userID string) ([]entities.Like, error) {
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

func (r *likeRepositoryImpl) Exists(userID, postID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, like := range r.likes {
		if like.UserID == userID && like.PostID == postID {
			return true
		}
	}
	return false
}
